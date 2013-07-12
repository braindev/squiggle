# Squiggle - an SQL query builder for golang

Have a program where it's awkward or impossible to compose an SQL query in sinlge place?  Tired of concatenating pieces of SQL manually?  Then Squiggle is for you!  Squiggle is not an ORM but a powerful light weight tool to build SQL queries piece by piece.  Squiggle draws inspiration from amazing projects like [Squel.js](http://hiddentao.github.io/squel/) and [Sequel](http://sequel.rubyforge.org).

## Project Goals

- Intuitive easy to use/read/remember DSL for query creation
- Cover 100% of common use cases for SQL query building

## Example Use

```go
import sgl "github.com/braindev/squiggle"

func main() {
  query := sgl.Select().
    AddFrom("users").
    AddField("username", "id", "created_at").
    Where(sgl.Or("username = $1", sql.And("is_admin = $2", "is_deleted = $3"))).
    Limit(10).Offset(5)
  fmt.Println(query.String()) 
  // ==> SELECT username, id, created_at FROM users WHERE username = $1 OR (is_admin = $2 AND is_deleted = $3) LIMIT 10 OFFSET 5

}
```

## API

#### `Select()` - creates a new query of type SELECT

```go
squiggle.Select()
```

#### `AddFrom(string/squiggle.From...)` - appends to the FROM clause of the query

```go
squiggle.Select().
  AddFrom("users", From{Schema: "db1", Table: "table1", Alias: "t1"}).
  AddFrom("foo").
  String()
// => "SELECT * FROM users, db1.table1 t1, foo"
```

#### `AddField(string/squiggle.Field...)` - appends to the field/expression portion of the query

```go
squiggle.Select().
  AddField("username", "id", squiggle.Field{Schema: "db", Table: "table", Name: "field", Alias: "f1"}).
  String()
// => "SELECT username, id, db.table.field AS f1"

squiggle.Select().
  AddField(squiggle.Field{Expression:"NOW()", Alias: "the_time"}).
  SetIdentifierQuotes("`")
  String()
// => "SELECT NOW() AS `the_time`" NOTE that the expression naturally doesn't get identifier quotes 
```

## TODO

- Support query types other than SELECT
