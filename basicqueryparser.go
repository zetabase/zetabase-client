package zetabase

import (
	"github.com/alecthomas/participle"
)

// Base comparison: an identifier then a comparison operator then a literal
type BQCompQ struct {
	Field    string   `@Ident`
	Operator string   `@("=" | ">" | "<" | ">=" | "<=" | "~" | "!=" | "not")`
	Value    *BQValue `@@`
}

// Root (expression): a comparison followed by one or more logical conjunctions...
type BQRootQ struct {
	Clause       *BQClause     `@@`
	Conjunctions []*BQLogicalQ `@@*`
}

// A "clause": a subexpression in parentheses OR a single field comparison
type BQClause struct {
	Subexpression *BQRootQ `"(" @@ ")"`
	Comparison    *BQCompQ ` | @@`
}

// A "logical": an operator with a second clause
type BQLogicalQ struct {
	Operation string    `@("and" | "or")`
	Query     *BQClause `@@`
}

// A value: string or number literal
type BQValue struct {
	String  *string  `   @String`
	Number  *float64 ` | @Float`
	Integer *int64   ` | @Int`
}

type BQParser struct {
	input string
}

func (b *BQValue) Value() interface{} {
	if b.String != nil {
		return *b.String
	} else if b.Integer != nil {
		return *b.Integer
	} else if b.Number != nil {
		return *b.Number
	}
	return nil
}

func (b *BQCompQ) ToQuery() SubQueryConvertible {
	fld := b.Field
	switch b.Operator {
	case ">":
		if b.Value.Number == nil && b.Value.Integer == nil {
			return nil
		}
		return QGt(fld, b.Value.Value())
	case "<":
		if b.Value.Number == nil && b.Value.Integer == nil {
			return nil
		}
		return QLt(fld, b.Value.Value())
	case ">=":
		if b.Value.Number == nil && b.Value.Integer == nil {
			return nil
		}
		return QGte(fld, b.Value.Value())
	case "<=":
		if b.Value.Number == nil && b.Value.Integer == nil {
			return nil
		}
		return QLte(fld, b.Value.Value())
	case "~":
		if b.Value.String == nil {
			return nil
		}
		qt := QText(fld, *b.Value.String)
		return qt
	case "!=", "not":
		neq := QNEq(fld, b.Value.Value())
		return neq
	default: // TODO jv - fill in the rest of these operators
		qe := QEq(fld, b.Value.Value())
		return qe
	}
}

func mergeLogicals(base SubQueryConvertible, logicals []*BQLogicalQ) SubQueryConvertible {
	if len(logicals) == 0 {
		return base
	} else {
		logical0 := logicals[0]
		switch logical0.Operation {
		case "and":
			return mergeLogicals(QAnd(base, logical0.Query.ToQuery()), logicals[1:])
		default:
			return mergeLogicals(QOr(base, logical0.Query.ToQuery()), logicals[1:])
		}
	}
}

func (b *BQRootQ) ToQuery() SubQueryConvertible {
	if b.Clause == nil {
		return nil
	}
	root := b.Clause.ToQuery()
	return mergeLogicals(root, b.Conjunctions)
}

func (b *BQClause) ToQuery() SubQueryConvertible {
	if b.Subexpression != nil {
		return b.Subexpression.ToQuery()
	} else {
		// do standard comparison
		return b.Comparison.ToQuery()
	}
}

func (b *BQParser) Parse() (*BQRootQ, error) {
	parser, err := participle.Build(&BQRootQ{})
	if err != nil {
		return nil, err
	}
	rig := &BQRootQ{}
	err = parser.ParseString(b.input, rig)
	if err != nil {
		return nil, err
	}

	//log.Printf("Rig = %v\n", rig)
	return rig, nil
}

func NewBQParser(s string) *BQParser {
	return &BQParser{input: s}
}
