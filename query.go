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

func Select(args ...string) *Query {
	q := new(Query)
	q.queryType = "SELECT"

	return q
}

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

func (q *Query) Limit(l int) *Query {
	q.limit = l

	return q
}

func (q *Query) Offset(o int) *Query {
	q.offset = o

	return q
}

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

func (q *Query) AddJoin(j ...Join) *Query {
	q.joins = append(q.joins, j...)
	return q
}

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

func (q *Query) Where(c Criteria) *Query {
	q.where = c
	return q
}

func (q *Query) AndWhere(c Criteria) *Query {
	if len(q.where.expressions) == 0 {
		q.where = c
	} else {
		q.where = And(q.where, c)
	}
	return q
}

func (q *Query) OrWhere(c Criteria) *Query {
	if len(q.where.expressions) == 0 {
		q.where = c
	} else {
		q.where = Or(q.where, c)
	}
	return q
}

func (q *Query) Having(c Criteria) *Query {
	q.having = c
	return q
}

func (q *Query) AndHaving(c Criteria) *Query {
	if len(q.having.expressions) == 0 {
		q.having = c
	} else {
		q.having = And(q.where, c)
	}
	return q
}

func (q *Query) OrHaving(c Criteria) *Query {
	if len(q.having.expressions) == 0 {
		q.having = c
	} else {
		q.having = Or(q.where, c)
	}
	return q
}

func (q *Query) identfierQuote(identifier string) string {
	return q.identifierLeftQuote + identifier + q.identifierRightQuote
}
