package squiggle

import (
	"fmt"
	"strings"
)

type Criteria struct {
	and bool
	expressions []interface{}
}

func (c Criteria) String() string {
	var parts []string

	for _, expression := range c.expressions {
		switch expression.(type) {
		default:
			panic(fmt.Sprintf("unexpected type %T in criteria", expression))
		case string:
			parts = append(parts, expression.(string))
		case Criteria:
			parts = append(parts, `(` + (expression.(Criteria)).String() + `)`)
		}
	}

	if c.and {
		return strings.Join(parts, " AND ")
	}
	return strings.Join(parts, " OR ")
}

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
