// Package squiggle provides functionality to build SQL queries.
//
// example
//
//   import(
//    	sgl "github.com/braindev/squiggle"
//   )
//
//   func main() {
//   	query := sgl.Select().
//   		AddFrom("users").
//   		AddField("username", "id", "created_at").
//   		Where(sgl.Or("username = $1", sql.And("is_admin = $2", "is_deleted = $3"))).
//   		Limit(10).Offset(5)
//   	fmt.Println(query.String()) // ==> SELECT username, id, created_at FROM users WHERE 
//   		// username = $1 OR (is_admin = $2 AND is_deleted = $3) LIMIT 10 OFFSET 5
//   }
package squiggle

import (
	"fmt"
	"strings"
)

type Join struct {
	Type   string
	On     Criteria
	Schema string
	Table  string
	Alias  string
}

type From struct {
	Schema string
	Table  string
	Alias  string
}

type Field struct {
	Schema     string
	Table      string
	Name       string
	Expression string
	Alias      string
}

type Grouping struct {
	Schema string
	Table  string
	Field  string
}

type Ordering struct {
	Schema string
	Table  string
	Field  string
	Desc   bool
}

type Query struct {
	queryType            string
	from                 []From
	fields               []Field
	groupings            []Grouping
	orderings            []Ordering
	joins                []Join
	where                Criteria
	having               Criteria
	limit                int
	offset               int
	identifierLeftQuote  string
	identifierRightQuote string
}

// Create a new SELECT query
// 	q := squiggle.Select()
func Select() *Query {
	q := new(Query)
	q.queryType = "SELECT"

	return q
}

// Sets identifier quotes.  The default is to have no identifier quotes.  This
// method accepts one or two string arguments.  If one string argument is 
// passed the identifier quotes be the same on both sides of the identifier.
// If two arguments are passed then the first will be on the left side of the 
// identifier and the second will be on the right.` 
// 
// Examples:
//
// 	squiggle.Select().SetIdentifierQuotes("`").Add(squiggle.From{Schema: "db1", Table: "users", Alias: "u"}).String()
// 	// => "SELECT * FROM `db1`.`users` `u`"
//
// 	squiggle.Select().SetIdentifierQuotes("[", "]").Add(sqiggle.Field{Schema: "db1", Table: "users", Name: "field", Alias: "f1"}).AddFrom("users").String()
// 	// => "SELECT [db1].[users].[field] AS [f1] FROM [users]"
func (q *Query) SetIdentifierQuotes(quotes ...string) *Query {
	if len(quotes) == 1 {
		q.identifierLeftQuote = quotes[0]
		q.identifierRightQuote = quotes[0]
	}
	if len(quotes) == 2 {
		q.identifierLeftQuote = quotes[0]
		q.identifierRightQuote = quotes[1]
	}
	return q
}

// Set a limit on a query
func (q *Query) Limit(l int) *Query {
	q.limit = l

	return q
}

// Set an offset on a query
func (q *Query) Offset(o int) *Query {
	q.offset = o

	return q
}

// Add a table to the from clause of a query.  This method will accept any
// number or arguments of type sting or squiggle.From.  If the argument is of
// type string it is the same as passing an argument
// squiggle.From{Table: "<string>"}
//
// 	// these two are the same
// 	sqiggle.Select().AddFrom("foo")
// 	sqiggle.Select().AddFrom(sqiggle.From{Table: "foo"})
func (q *Query) AddFrom(froms ...interface{}) *Query {
	for _, from := range froms {
		switch from.(type) {
		default:
			panic(fmt.Sprintf("unexpected type %T used in AddFrom()", from))
		case string:
			q.from = append(q.from, From{Table: from.(string)})
		case From:
			q.from = append(q.from, from.(From))
		}
	}

	return q
}
// Adds a field to a select query.  The method accepts a variable number of 
// arguments of type string or sqiggle.Field.  When an argument is simply of
// type string it's treated as if it were sqiggle.Field{Name: <string>}
//
// 	sqiggle.Select().SetIdentifierQuotes(`"`).AddField("user_type", Field{Expression: "AVG(age)"})
// 	// => SELECT "user_type", AVG(age)
func (q *Query) AddField(fields ...interface{}) *Query {
	for _, field := range fields {
		switch field.(type) {
		default:
			panic(fmt.Sprintf("unexpected type %T used in AddField()", field))
		case string:
			q.fields = append(q.fields, Field{Name: field.(string)})
		case Field:
			q.fields = append(q.fields, field.(Field))
		}
	}

	return q
}

// Add ordering to a query.  This method will accept any number of aguments of
// type string or squiggle.Ordering.  Passing a string is a shortcut for
// passing squiggle.Ordering{Field: "<string>", Desc: false}
//
// 	squiggle.Select().AddOrdering("foo", squiggle.Ordering{Field: "Bar", Desc: true})
// 	// => SELECT ... ORDER BY foo ASC, Bar DESC
func (q *Query) AddOrdering(orderings ...interface{}) *Query {
	for _, ordering := range orderings {
		switch ordering.(type) {
		default:
			panic(fmt.Sprintf("unexpected type %T used in AddOrdering()", ordering))
		case Ordering:
			q.orderings = append(q.orderings, ordering.(Ordering))
		case string:
			q.orderings = append(q.orderings, Ordering{Field: ordering.(string)})
		}
	}

	return q
}

// Add groupings to a query.  Accepts any number of aguments of type string or
// Grouping.  When a string is passed as an argument it is the same as passing
// squiggle.Grouping{Field: "<string>"}
//
// 	squiggle.Select().AddGrouping("foo", Grouping{Field: "bar", Table: "baz"})
// 	// => SELECT ... GROUP BY foo, baz.bar
func (q *Query) AddGrouping(groupings ...interface{}) *Query {
	for _, grouping := range groupings {
		switch grouping.(type) {
		default:
			panic(fmt.Sprintf("unexpected type %T used in AddGrouping()", grouping))
		case Grouping:
			q.groupings = append(q.groupings, grouping.(Grouping))
		case string:
			q.groupings = append(q.groupings, Grouping{Field: grouping.(string)})
		}
	}
	return q
}

// Add joins to a query.  Accepts any number of arguments of type Join
func (q *Query) AddJoin(j ...Join) *Query {
	q.joins = append(q.joins, j...)
	return q
}

// A generic way to add joins, orderings, fields, froms, and criteria to a
// query.  This method will accept any number of arguments of types
// Grouping, Ordering, Field, From, or Join
func (q *Query) Add(things ...interface{}) *Query {
	for _, thing := range things {
		switch thing.(type) {
		default:
			panic(fmt.Sprintf("unexpected type %T used in Add()", thing))
		case Grouping:
			q.AddGrouping(thing.(Grouping))
		case Ordering:
			q.AddOrdering(thing.(Ordering))
		case Field:
			q.AddField(thing.(Field))
		case From:
			q.AddFrom(thing.(From))
		case Join:
			q.AddJoin(thing.(Join))
		}
	}

	return q
}

// returns the fields and expressions portions of the query as an SQL string
func (q *Query) FieldsString() string {
	sql := ""
	var fields []string
	if len(q.fields) == 0 {
		sql = sql + " *"
	} else {
		for _, field := range q.fields {
			fieldStr := ""
			if field.Expression == "" {
				if field.Schema != `` {
					fieldStr = fieldStr + q.identfierQuote(field.Schema) + "."
				}
				if field.Table != `` {
					fieldStr = fieldStr + q.identfierQuote(field.Table) + "."
				}
				fieldStr = fieldStr + q.identfierQuote(field.Name)
			} else {
				fieldStr = fieldStr + field.Expression
			}
			if field.Alias != `` {
				fieldStr = fieldStr + " AS " + q.identfierQuote(field.Alias)
			}
			fields = append(fields, fieldStr)
		}
		sql = sql + " " + strings.Join(fields, ", ")
	}

	return sql
}

// returns the from portion of the query as an SQL string
func (q *Query) FromString() string {
	sql := ""
	if len(q.from) > 0 {
		var fromStrings []string
		sql = sql + " FROM "
		for _, from := range q.from {
			fromStr := ""
			if from.Schema != "" {
				fromStr = fromStr + q.identfierQuote(from.Schema) + "."
			}
			fromStr = fromStr + q.identfierQuote(from.Table)
			if from.Alias != "" {
				fromStr = fromStr + " " + q.identfierQuote(from.Alias)
			}
			fromStrings = append(fromStrings, fromStr)
		}
		sql = sql + strings.Join(fromStrings, ", ")
	}

	return sql
}

// returns the joins portion of the query as a string
func (q *Query) JoinsString() string {
	sql := ""
	joinStrings := []string{}
	for _, join := range q.joins {
		joinStr := " " + strings.ToUpper(join.Type) + " JOIN "
		if join.Schema != "" {
			joinStr = joinStr + q.identfierQuote(join.Schema) + "."
		}
		joinStr = joinStr + q.identfierQuote(join.Table)
		if join.Alias != "" {
			joinStr = joinStr + " " + q.identfierQuote(join.Alias)
		}
		if len(join.On.expressions) > 0 {
			joinStr = joinStr + " ON " + join.On.String()
		}
		joinStrings = append(joinStrings, joinStr)
	}

	sql = sql + strings.Join(joinStrings, " ")

	return sql
}

// returns the grouping portion of the query as a string
func (q *Query) GroupingsString() string {
	sql := ""
	if len(q.groupings) > 0 {
		sql = sql + " GROUP BY "
		var groupingsStrings []string
		for _, grouping := range q.groupings {
			groupingStr := q.identfierQuote(grouping.Field)
			if grouping.Table != "" {
				groupingStr = q.identfierQuote(grouping.Table) + "." + groupingStr
			}
			if grouping.Schema != "" {
				groupingStr = q.identfierQuote(grouping.Schema) + "." + groupingStr
			}
			groupingsStrings = append(groupingsStrings, groupingStr)
		}
		sql = sql + strings.Join(groupingsStrings, ", ")
	}

	return sql
}

// 	returns the orderings portion of the query as a string
func (q *Query) OrderingsString() string {
	sql := ""
	if len(q.orderings) > 0 {
		var orderingsStrings []string
		for _, ordering := range q.orderings {
			orderingStr := q.identfierQuote(ordering.Field)
			if ordering.Table != "" {
				orderingStr = q.identfierQuote(ordering.Table) + "." + orderingStr
			}
			if ordering.Schema != "" {
				orderingStr = q.identfierQuote(ordering.Schema) + "." + orderingStr
			}
			if ordering.Desc {
				orderingStr = orderingStr + " DESC"
			} else {
				orderingStr = orderingStr + " ASC"
			}
			orderingsStrings = append(orderingsStrings, orderingStr)
		}
		sql = sql + " ORDER BY " + strings.Join(orderingsStrings, ", ")
	}

	return sql
}

// Turns the query into a string of SQL
func (q *Query) String() string {
	// <QUERY TYPE>
	sql := q.queryType

	// <FIELDS>
	sql = sql + q.FieldsString()

	// <FROM>
	sql = sql + q.FromString()

	// <JOINS>
	sql = sql + q.JoinsString()

	// <WHERE>
	if len(q.where.expressions) > 0 {
		sql = sql + " WHERE " + q.where.String()
	}

	// <GROUPS>
	sql = sql + q.GroupingsString()

	// <HAVING>
	if len(q.having.expressions) > 0 {
		sql = sql + " HAVING " + q.having.String()
	}

	// <ORDER>
	sql = sql + q.OrderingsString()

	// <LIMIT OFFSET>
	if q.limit > 0 {
		sql = sql + fmt.Sprintf(" LIMIT %d", q.limit)
	}
	if q.offset > 0 {
		sql = sql + fmt.Sprintf(" OFFSET %d", q.offset)
	}

	return sql
}

// Add criteria to the "where" portion of a query.  This method accepts a
// parameter of type string or squiggle.Criteria.  The criteria can be
// created by using the squiggle.And and squiggle.Or functions.  When an
// argument of type string is passed it's the same as passing
// squiggle.And("<string>")  Note that Where will replace and previously
// created criteria.
//
// 	squiggle.Select().Where("a=?")
// 	// => WHERE a=?
// 	squiggle.Select().Where(squiggle.Or("a=?", squiggle.And("b=?", "c=?)))
// 	// => WHERE a=? OR (b=? AND c=?)
func (q *Query) Where(c interface{}) *Query {
	var criteria Criteria
	switch c.(type) {
	default:
		panic(fmt.Sprintf("unexpected type %T used in Where()", c))
	case string:
		criteria = And(c.(string))
	case Criteria:
		criteria = c.(Criteria)
	}

	q.where = criteria
	return q
}

// This is the same as the Where() method except it appends with AND logic
//
// 	squiggle.Select().Where("a=1").AndWhere("b=2")
// 	// => SELECT ... WHERE a=1 AND (b=2)
func (q *Query) AndWhere(c interface{}) *Query {
	var criteria Criteria
	switch c.(type) {
	default:
		panic(fmt.Sprintf("unexpected type %T used in AndWhere()", c))
	case string:
		criteria = And(c.(string))
	case Criteria:
		criteria = c.(Criteria)
	}
	
	if len(q.where.expressions) == 0 {
		q.where = criteria
	} else {
		q.where = And(q.where, criteria)
	}
	return q
}

// This is the same as the Where() method except it appends with OR logic
//
// 	squiggle.Select().Where("a=1").OrWhere("b=2")
// 	// => SELECT ... WHERE a=1 OR (b=2)
func (q *Query) OrWhere(c interface{}) *Query {
	var criteria Criteria
	switch c.(type) {
	default:
		panic(fmt.Sprintf("unexpected type %T used in OrWhere()", c))
	case string:
		criteria = And(c.(string))
	case Criteria:
		criteria = c.(Criteria)
	}

	if len(q.where.expressions) == 0 {
		q.where = criteria
	} else {
		q.where = Or(q.where, criteria)
	}
	return q
}

// This is the same as the Where() method except it add criteria to the HAVING
// portion of the query rather than the WHERE portion
func (q *Query) Having(c interface{}) *Query {
	var criteria Criteria
	switch c.(type) {
	default:
		panic(fmt.Sprintf("unexpected type %T used in Having()", c))
	case string:
		criteria = And(c.(string))
	case Criteria:
		criteria = c.(Criteria)
	}

	q.having = criteria
	return q
}

// This is the same as the Having() method except it appends with AND logic
//
// 	squiggle.Select().Having("a=1").AndHaving("b=2")
// 	// => SELECT ... HAVING a=1 AND (b=2)
func (q *Query) AndHaving(c interface{}) *Query {
	var criteria Criteria
	switch c.(type) {
	default:
		panic(fmt.Sprintf("unexpected type %T used in AndHaving()", c))
	case string:
		criteria = And(c.(string))
	case Criteria:
		criteria = c.(Criteria)
	}

	if len(q.having.expressions) == 0 {
		q.having = criteria
	} else {
		q.having = And(q.where, criteria)
	}
	return q
}

// This is the same as the Having() method except it appends with OR logic
//
// 	squiggle.Select().Having("a=1").OrHaving("b=2")
// 	// => SELECT ... HAVING a=1 OR (b=2)
func (q *Query) OrHaving(c interface{}) *Query {
	var criteria Criteria
	switch c.(type) {
	default:
		panic(fmt.Sprintf("unexpected type %T used in OrHaving()", c))
	case string:
		criteria = And(c.(string))
	case Criteria:
		criteria = c.(Criteria)
	}

	if len(q.having.expressions) == 0 {
		q.having = criteria
	} else {
		q.having = Or(q.where, criteria)
	}
	return q
}

func (q *Query) identfierQuote(identifier string) string {
	return q.identifierLeftQuote + identifier + q.identifierRightQuote
}
