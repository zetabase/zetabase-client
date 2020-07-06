package zetabase

import (
	"log"
	"testing"
)

func Test_BasicParsing(t *testing.T) {
	parser := NewBQParser("rig = \"hello\"")
	res, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	} else {
		log.Printf("Res: %v\n", res)
	}

	parser = NewBQParser("rig BLAH \"hello\"")
	_, err = parser.Parse()
	if err == nil {
		t.Fatalf("Should have had an error.")
	} else {
		log.Printf("Correct error: %s\n", err.Error())
	}


	parser = NewBQParser(" rig   = \"hello\"")
	_, err = parser.Parse()
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	} else {
		log.Printf("Correctly parsed despite weird whitespace.\n")
	}
}

func Test_BasicParsing_and(t *testing.T) {
	parser := NewBQParser("rig = \"hello\" and tim = 55")
	res, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	} else {
		if res.Conjunctions == nil {
			t.Fatalf("Wrong dimensions: %d x %v", len(res.Conjunctions), res.Clause)
		} else {
			andCl := res.Conjunctions[0]
			if andCl.Query.Comparison.Field != "tim" {
				t.Fatalf("Wrong fields: %s x %s", andCl.Query.Comparison.Field, andCl.Query.Comparison.Field)
			} else if res.Clause.Comparison.Value.String == nil || andCl.Query.Comparison.Value.Integer == nil {
				t.Fatalf("Wrong values: %v x %v", andCl.Query.Comparison.Value, andCl.Query.Comparison.Value)
			} else {
				log.Printf("Parsed correctly: %s x %d.", *res.Clause.Comparison.Value.String, *andCl.Query.Comparison.Value.Integer)
			}
		}
	}
}

func Test_BasicParsing_subexpr(t *testing.T) {
	parser := NewBQParser("rig = \"hello\" and (tim = 55 and fff = 222)")
	res, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	} else {
		if res.Conjunctions == nil {
			t.Fatalf("Wrong dimensions: %d x %v", len(res.Conjunctions), res.Clause)
		} else {
			if res.Conjunctions[0].Query.Subexpression == nil {
				t.Fatalf("Should have subexpression: %v", res)
			}
			if res.Conjunctions[0].Query.Subexpression == nil || res.Conjunctions[0].Query.Subexpression.Conjunctions == nil {
				t.Fatalf("Should have subexpression conjunction: %v", res)
			} else {
				fldN := res.Conjunctions[0].Query.Subexpression.Conjunctions[0].Query.Comparison.Field
				if fldN != "fff" {
					t.Fatalf("Wrong field name in final conjunction: %s\n", fldN)
				} else {
					log.Printf("Got correct field names: %s, %s, %s\n", res.Clause.Comparison.Field, res.Conjunctions[0].Query.Subexpression.Clause.Comparison.Field, res.Conjunctions[0].Query.Subexpression.Conjunctions[0].Query.Comparison.Field)
				}
			}
		}
	}
}


func Test_BasicParsing_matchop(t *testing.T) {
	parser := NewBQParser("rig ~ \"hello goodbye\"")
	res, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parsing error: %s\n", err.Error())
	} else {
		if res.Clause != nil && res.Clause.Comparison != nil {
			if res.Clause.Comparison.Operator != "~" {
				t.Fatalf("Wrong operator: %s\n", res.Clause.Comparison.Operator)
			} else if res.Clause.Comparison.Field != "rig" {
				t.Fatalf("Wrong field: %s\n", res.Clause.Comparison.Field)
			} else if res.Clause.Comparison.Value.String == nil {
				t.Fatalf("No query data.\n")
			} else if *res.Clause.Comparison.Value.String != "hello goodbye" {
				t.Fatalf("Wrong query: %s\n", *res.Clause.Comparison.Value.String)
			} else {
				log.Printf("Successfully parsed match operator %s: %v\n", res.Clause.Comparison.Operator, res)
			}
		} else {
			t.Fatalf("Wrong structure: %v\n", res)
		}
	}
}

func Test_Cmd1(t *testing.T) {
	qryStr := "(age > 30.0)"
	parser := NewBQParser(qryStr)
	res, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parsing error: %s\n", err.Error())
	} else {
		qry := res.ToQuery().ToSubQuery("tblowner", "tbl")
		log.Printf("Query: %v\n", qry)
	}
}
/*
*/
