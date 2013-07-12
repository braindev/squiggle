package squiggle

import (
	"fmt"
	"strings"
)

type Criteria struct {
	and         bool
	expressions []interface{}
}

// returns a criteria as an SQL string
func (c Criteria) String() string {
	var parts []string

	for _, expression := range c.expressions {
		switch expression.(type) {
		default:
			panic(fmt.Sprintf("unexpected type %T in criteria", expression))
		case string:
			parts = append(parts, expression.(string))
		case Criteria:
			parts = append(parts, `(`+(expression.(Criteria)).String()+`)`)
		}
	}

	if c.and {
		return strings.Join(parts, " AND ")
	}
	return strings.Join(parts, " OR ")
}

// Creates a criteria with the logic of AND.  Accepts any number of arguments
// of type string or squiggle.Criteria.
//
// 	squiggle.And("a=1", squiggle.Or("b=2", "c=3", squiggle.And("d=4", "e=5")))
// 	// => a=1 AND (b=2 OR c=3 OR (d=4 AND e=5))
func And(args ...interface{}) Criteria {
	c := Criteria{and: true}
	for _, arg := range args {
		switch arg.(type) {
		default:
			panic(fmt.Sprintf("unexpected type %T used in And()", arg))
		case string:
			c.expressions = append(c.expressions, arg)
		case Criteria:
			c.expressions = append(c.expressions, arg)
		}
	}

	return c
}

// This is the same as the And() function except it creates a Criteria with OR
// logic
func Or(args ...interface{}) Criteria {
	c := Criteria{and: false}
	for _, arg := range args {
		switch arg.(type) {
		default:
			panic(fmt.Sprintf("unexpected type %T used in Or()", arg))
		case string:
			c.expressions = append(c.expressions, arg)
		case Criteria:
			c.expressions = append(c.expressions, arg)
		}
	}

	return c
}
