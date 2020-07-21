package main

import (
	"fmt"
	"math/rand"
)

func Logf(fmtStr string, argv ...interface{}) {
	rig := fmt.Sprintf(fmtStr, argv...)
	fmt.Printf("ZB>  %s\n", rig)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randLetterStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
