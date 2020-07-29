package zetabase

import (
	"errors"
	"log"
	"math/rand"
	"testing"
)

func prepData(n int) ([]string, [][]byte) {
	var ks []string
	var bs [][]byte
	for i := 0; i < n; i++ {
		k := randJsonStringRunes(12)
		dlen := rand.Intn(100) + 28
		v := []byte(randJsonStringRunes(dlen))
		ks = append(ks, k)
		bs = append(bs, v)
	}
	return ks, bs
}

func validatePage(pg [][]byte, maxBytes int) error {
	s, n := 0, 0
	defer func(){
		log.Printf("Total size: %d bytes (%d items)\n", s, n)
	}()
	for _, v := range pg {
		s += len(v)
		n += 1
	}
	if s > maxBytes {
		log.Printf("Total size %d > %d max\n", s, maxBytes)
		return errors.New("TooBig")
	}
	return nil
}

func Test_Pagify(t *testing.T) {
	keys, values := prepData(100)
	maxBytes := 800

	ppgs := makePutPages(nil, keys, values, uint64(maxBytes))
	_, vpgs, _ := ppgs.pagify()

	for _, pg := range vpgs {
		if e := validatePage(pg, maxBytes); e != nil {
			t.Fatalf("Error validating page: %s", e.Error())
		}
	}

	errpgs := makePutPages(nil, []string{"k"}, [][]byte{[]byte(randJsonStringRunes(GrpcMaxBytes))}, GrpcMaxBytes / 2)
	_, _, err := errpgs.pagify()
	if err == nil {
		t.Fatalf("There should be an error")
	} else {
		log.Printf("Got error: %s!\n", err.Error())
	}
}