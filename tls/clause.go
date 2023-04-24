package tls

// Clause
type Clause struct {
	Condition []string
	Params    []interface{}
	End       bool
}

func NewClause() Clause {
	return Clause{
		Condition: make([]string, 0),
		Params:    make([]interface{}, 0),
	}
}
