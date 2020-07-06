package zetabase

import (
	"log"
	"testing"
)

func Test_IsSemVerVersionAtLeast(t *testing.T) {
	userVersions := []string{"1.2", "1.1", "1.1-alpha", "1.4.1", "1.2", "1.2-alpha", "0.0.0.9", "1.1.0.1"}
	minVersions := []string{"1.2", "1.1-alpha", "1.1-alpha", "1.3", "1.3.1", "1.2", "0.0.9", "1.1.0.1.1"}
	for i := 0; i < len(userVersions); i++ {
		minV := minVersions[i]
		usrV := userVersions[i]
		sbt := IsSemVerVersionAtLeast(usrV, minV)
		correctValu := true
		if i >= 4 {
			correctValu = false
		}
		if sbt != correctValu {
			t.Fatalf("Should be %v: %s vs. %s", correctValu, usrV, minV)
		} else {
			log.Printf("Correct (%s vs %s)\n", usrV, minV)
		}
	}
}