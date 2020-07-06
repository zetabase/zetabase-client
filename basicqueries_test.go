package zetabase

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"testing"
)

func Test_SimpleDSL(t *testing.T) {
	qry := QAnd(QOr(QEq("uid", "jason"), QEq("uid", "charlotte")), QEq("type", 1))
	rig := qry.ToSubQuery("usr", "table")
	bs, err := proto.Marshal(rig)
	if err != nil {
		t.Fatalf("Error: %s", err.Error())
	} else {
		fmt.Printf("Success: got %d bytes after serialization\n", len(bs))
	}
}
