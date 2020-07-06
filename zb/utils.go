package main

import (
	"fmt"
)

func Logf(fmtStr string, argv ...interface{}) {
	rig := fmt.Sprintf(fmtStr, argv...)
	fmt.Printf("ZB>  %s\n", rig)
}
