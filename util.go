package zetabase

import (
	"github.com/hashicorp/go-version"
	"github.com/zetabase/zetabase-client/zbprotocol"
	"regexp"
	"time"
)

func IsSemVerVersionAtLeast(userVersion, minVersion string) bool {
	usr, err := version.NewVersion(userVersion)
	if err != nil {
		return false
	}
	min, err := version.NewVersion(minVersion)
	if err != nil {
		return false
	}
	if usr.LessThan(min) {
		return false
	}
	return true
}

func EmptySignature() *zbprotocol.EcdsaSignature {
	return &zbprotocol.EcdsaSignature{
		R: "",
		S: "",
	}
}

type NonceMaker struct {
	initValu int64
	inChan   chan chan int64
}

func (n *NonceMaker) runLoop() {
	for ch := range n.inChan {
		n.initValu += 1
		ch <- n.initValu
	}
}

func (n *NonceMaker) Get() int64 {
	ch := make(chan int64, 1)
	n.inChan <- ch
	v := <-ch
	return v
}

func NewNonceMaker() *NonceMaker {
	n := &NonceMaker{
		initValu: time.Now().UnixNano(),
		inChan:   make(chan chan int64, 16),
	}
	go n.runLoop()
	return n
}
func ValidatePhoneNumber(s string) bool {
	reg, err := regexp.Compile("^\\+(?:[0-9] ?){6,14}[0-9]$")
	if err != nil {
		return false
	}
	return reg.MatchString(s)
}
