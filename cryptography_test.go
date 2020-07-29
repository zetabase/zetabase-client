package zetabase

import (
	"github.com/zetabase/zetabase-client/zbprotocol"
	"math/rand"
)

var jsonRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ{}:,\"[]")

func randJsonStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = jsonRunes[rand.Intn(len(jsonRunes))]
	}
	return string(b)
}

func makeFakePairs(n int) []*zbprotocol.DataPair {
	var arr []*zbprotocol.DataPair
	for i := 0; i < n; i++ {
		k := randJsonStringRunes(32)
		v := randJsonStringRunes(256)
		arr = append(arr, &zbprotocol.DataPair{
			Key:   k,
			Value: []byte(v),
		})
	}
	return arr
}

/*
func Test_MultiPutExtraSigningBytesBenchmark(t *testing.T) {
	nPairs := 1000000
	nIters := 10

	testPairs := makeFakePairs(nPairs)

	t0 := time.Now()
	for i := 0; i < nIters; i++ {
		MultiPutExtraSigningBytes(testPairs)
	}
	t1 := time.Now()
	for i := 0; i < nIters; i++ {
		MultiPutExtraSigningBytesMurmur3(testPairs)
	}
	t2 := time.Now()
	for i := 0; i < nIters; i++ {
		MultiPutExtraSigningBytesAbbrev(testPairs)
	}
	t3 := time.Now()
	for i := 0; i < nIters; i++ {
		MultiPutExtraSigningBytesMurmur3Sliding(testPairs)
	}
	t4 := time.Now()

	ms := int(t1.Sub(t0).Seconds()*1000)
	ms2 := int(t2.Sub(t1).Seconds()*1000)
	ms3 := int(t3.Sub(t2).Seconds()*1000)
	ms4 := int(t4.Sub(t3).Seconds()*1000)

	log.Printf("MD5: %d iterations of %d pairs each took:\t\t %d milliseconds.", nIters, nPairs, ms)
	log.Printf("Murmur3: %d iterations of %d pairs each took:\t\t %d milliseconds.", nIters, nPairs, ms2)
	log.Printf("MD5-Abbrev: %d iterations of %d pairs each took:\t %d milliseconds.", nIters, nPairs, ms3)
	log.Printf("Murmur3-sliding: %d iterations of %d pairs each took:\t %d milliseconds.", nIters, nPairs, ms4)
}
*/