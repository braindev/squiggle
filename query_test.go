package squiggle

import (
	"testing"
)

func Test_Select(t *testing.T) {
	q := Select()
	if q.queryType != "SELECT" {
		t.Errorf("query type for Select() should be \"SELECT\"")
	}
}

func Test_SetIdentifierQuotes(t *testing.T) {
	q1 := Select().SetIdentifierQuotes(`"`)
	if q1.identifierLeftQuote != `"` && q1.identifierRightQuote != `"` {
		t.Error("identifier quotes are incorrect")
	}

	q2 := Select().SetIdentifierQuotes(`[`, `]`)
	if q2.identifierLeftQuote != `[` && q2.identifierRightQuote != `]` {
		t.Error("identifier quotes are incorrect")
	}
}

func Test_Limit(t *testing.T) {
	q1 := Select().Limit(10)

	if q1.limit != 10 {
		t.Error("limit is incorrect")
	}
}

func Test_Offset(t *testing.T) {
	q1 := Select().Offset(20)

	if q1.offset != 20 {
		t.Error("offset is incorrect")
	}
}

func Test_AddFrom(t *testing.T) {
	q := Select().AddFrom(From{Table: "table1"}, "table2")
	if len(q.from) != 2 {
		t.Errorf("wrong number of froms")
	} else {
		if q.from[0].Table != "table1" {
			t.Errorf("table1 is missing or in wrong order in froms slice")
		}
		if q.from[1].Table != "table2" {
			t.Errorf("table2 is missing or in wrong order in froms slice")
		}
	}
}

func Test_AddField(t *testing.T) {
	q := Select().AddField("field1", Field{Name: "field2"})
	if len(q.fields) != 2 {
		t.Errorf("wrong number of fields")
	} else {
		if q.fields[0].Name != "field1" {
			t.Error("field1 is missing on in wrong order in fields slice")
		}
		if q.fields[1].Name != "field2" {
			t.Error("field2 is missing on in wrong order in fields slice")
		}
	}
}

func Test_AddOrdering(t *testing.T) {
	q1 := Select().AddOrdering("foo")

	q2 := Select().AddOrdering(Ordering{Schema: "a", Table: "b", Field: "c", Desc: true})

	if len(q1.orderings) != 1 {
		t.Error("Invalid number of orderings")
	} else {
		if q1.orderings[0].Field != "foo" || q1.orderings[0].Desc != false {
			t.Error("Invalid ordering from string parameter")
		}
	}

	if len(q2.orderings) != 1 {
		t.Error("Invalid number of orderings")
	} else {
		if q2.orderings[0].Schema != "a" || q2.orderings[0].Table != "b" || q2.orderings[0].Desc != true || q2.orderings[0].Field != "c" {
			t.Error("Invalid ordering from string parameter")
		}
	}
}

func Test_AddGrouping(t *testing.T) {
	q1 := Select().AddGrouping("foo", Grouping{Field: "bar"})

	if q1.groupings[0].Field != "foo" || q1.groupings[1].Field != "bar" {
		t.Error("Invalid groupings")
	}
}

func Test_AddJoin(t *testing.T) {
	q1 := Select().AddJoin(Join{Type: "left", Table: "foo", On: And("a = b")})
	join := q1.joins[0]
	if join.Type != "left" || join.Table != "foo" || !join.On.and || join.On.expressions[0].(string) != "a = b" {
		t.Error("Invalid join")
	}
}

func Test_Add(t *testing.T) {
	q1 := Select().Add(Field{Name: "baz"}, From{Table: "bar"}, Join{Type: "inner", Table: "foo", On: And("bar.foo_id = foo.id")}, Grouping{Table: "bar", Field: "baz"}, Ordering{Field: "baz"})
	if len(q1.joins) != 1 {
		t.Error("foin count is wrong")
	}
	if len(q1.from) != 1 {
		t.Error("from count is wrong")
	}
	if len(q1.fields) != 1 {
		t.Error("field count is wrong")
	}
	if len(q1.groupings) != 1 {
		t.Error("grouping count is wrong")
	}
	if len(q1.orderings) != 1 {
		t.Error("ordering count is wrong")
	}
}

func Test_FieldsString(t *testing.T) {
	var expected string

	q1 := Select()

	if q1.FieldsString() != " *" {
		t.Error("FieldsString() invalid for no selected fields")
	}

	q2 := Select().
		SetIdentifierQuotes("[", "]").
		Add(Field{Schema: "db1", Table: "table1", Name: "field1", Alias: "alias1"}).
		AddField("field2", Field{Expression: "SUM(x)", Alias: "exp_alias"})
	expected = ` [db1].[table1].[field1] AS [alias1], [field2], SUM(x) AS [exp_alias]`
	if str := q2.FieldsString(); str != expected {
		t.Errorf("FieldString invalid. expected: `%s` found `%s`", expected, str)
	}
}

func Test_FromString(t *testing.T) {
	q1 := Select()

	if q1.FromString() != "" {
		t.Error("FromString() should return an empty string for no fields")
	}

	q2 := Select().
		AddFrom("table1").
		SetIdentifierQuotes("[", "]").
		Add(From{Schema: "db2", Table: "table2", Alias: "alias2"})

	expected := " FROM [table1], [db2].[table2] [alias2]"
	str := q2.FromString()
	if expected != str {
		t.Errorf("FromString() invalid.  Expected `%s` found `%s`", expected, str)
	}
}

func Test_JoinsString(t *testing.T) {
	q1 := Select()

	if q1.JoinsString() != "" {
		t.Error("JoinsString() should return an empty string for a query with no joins")
	}

	q2 := Select().
		SetIdentifierQuotes("[", "]").
		Add(Join{Type: "inner", Schema: "db1", Table: "foo", Alias: "a", On: And("bar.foo_id = foo.id")})

	expected := " INNER JOIN [db1].[foo] [a] ON bar.foo_id = foo.id"
	str := q2.JoinsString()
	if str != expected {
		t.Errorf("JoinsString() invalid expected `%s` found `%s`", expected, str)
	}
}

func Test_GroupingsString(t *testing.T) {
	q1 := Select()
	if q1.GroupingsString() != "" {
		t.Error("GroupingsString() should return an empty string for a query with no groupings")
	}

	q1.AddGrouping(Grouping{Schema: "db", Table: "table", Field: "f1"})
	q1.AddGrouping(Grouping{Schema: "db", Table: "table", Field: "f2"})
	q1.SetIdentifierQuotes(`[`, `]`)
	str := q1.GroupingsString()
	expected := " GROUP BY [db].[table].[f1], [db].[table].[f2]"
	if expected != str {
		t.Errorf("GroupingsString() returned `%s` expected `%s`", str, expected)
	}
}

func Test_OrderingsString(t *testing.T) {
	q1 := Select()
	if q1.OrderingsString() != "" {
		t.Error("OrderingsString() should return an empty string for a query with no orderings")
	}
	q1.AddOrdering("foo", Ordering{Schema: "db", Table: "t", Field: "f", Desc: true}).
		SetIdentifierQuotes(`[`, `]`)

	expected := " ORDER BY [foo] ASC, [db].[t].[f] DESC"
	str := q1.OrderingsString()

	if str != expected {
		t.Errorf("OrderingsString() returned `%s` expected `%s`", str, expected)
	}
}

func Test_String(t *testing.T) {
	q1 := Select().
		AddFrom("users").
		AddField("id", "username").
		Where(Or("username = $1", "updated_at >= $2")).
		AddOrdering("username").
		SetIdentifierQuotes("`")

	str := q1.String()
	expected := "SELECT `id`, `username` FROM `users` WHERE username = $1 OR updated_at >= $2 ORDER BY `username` ASC"
	if str != expected {
		t.Errorf("String() returned `%s` expected `%s`", str, expected)
	}

	q2 := Select().
		AddFrom("users").
		AddField(Field{Expression: "COUNT(*)", Alias: "user_type_count"}).
		AddGrouping("user_type").
		Having(And("COUNT(*) < 10", "COUNT(*) < 20"))
	str = q2.String()
	expected = "SELECT COUNT(*) AS user_type_count FROM users GROUP BY user_type HAVING COUNT(*) < 10 AND COUNT(*) < 20"
	if str != expected {
		t.Errorf("String() returned `%s` expected `%s`", str, expected)
	}
}

func Test_Where(t *testing.T) {
	q1 := Select().Where(And("a = 1"))

	if q1.where.expressions[0].(string) != "a = 1" {
		t.Error("Where() did not create expected criteria")
	}

	q1.Where(And("b = 1"))

	if q1.where.expressions[0].(string) != "b = 1" {
		t.Error("Where() did not override previous criteria")
	}
}

func Test_AndWhere(t *testing.T) {
	q1 := Select().Where(And("a=1")).AndWhere(And("b=1"))
	if !q1.where.and || q1.where.expressions[1].(Criteria).expressions[0].(string) != "b=1" {
		t.Error("AndWhere() did not create the expected criteria")
	}
}

func Test_OrWhere(t *testing.T) {
	q1 := Select().Where(And("a=1")).OrWhere(And("b=1"))
	if q1.where.and || q1.where.expressions[1].(Criteria).expressions[0].(string) != "b=1" {
		t.Error("OrWhere() did not create the expected criteria")
	}
}

func Test_Having(t *testing.T) {
	q1 := Select().Having(And("a = 1"))

	if q1.having.expressions[0].(string) != "a = 1" {
		t.Error("Having() did not create expected criteria")
	}

	q1.Having(And("b = 1"))

	if q1.having.expressions[0].(string) != "b = 1" {
		t.Error("Having() did not override previous criteria")
	}
}

func Test_AndHaving(t *testing.T) {
	q1 := Select().Having("a=1").AndHaving("b=1")
	if !q1.having.and || q1.having.expressions[1].(Criteria).expressions[0].(string) != "b=1" {
		t.Error("AndHaving() did not create the expected criteria")
	}
}

func Test_OrHaving(t *testing.T) {
	q1 := Select().Having("a=1").OrHaving("b=1")
	if q1.having.and || q1.having.expressions[1].(Criteria).expressions[0].(string) != "b=1" {
		t.Error("OrHaving() did not create the expected criteria")
	}
}
