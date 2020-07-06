package main

import "testing"

func Test_ValidatePhoneNumber(t *testing.T) {
	num := "+12035613094"
	if !validatePhoneNumber(num) {
		t.Fatalf("Should have validated %s\n", num)
	}
	num = "2035613094"
	if validatePhoneNumber(num) {
		t.Fatalf("Should not have validated %s\n", num)
	}
	num = "+120A5613094"
	if validatePhoneNumber(num) {
		t.Fatalf("Should not have validated %s\n", num)
	}
}
