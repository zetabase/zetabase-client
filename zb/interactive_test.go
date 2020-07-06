package main

import "testing"

func Test_shellParse(t *testing.T) {
	cmd := "ls -l \"jason's folder\""
	arr := shellParse(cmd)
	if len(arr) != 3 {
		t.Fatalf("Wrong len %d", len(arr))
	}
	if arr[0] != "ls" || arr[1] != "-l" || arr[2] != "jason's folder" {
		t.Fatalf("Wrong values %v", arr)
	}
	cmd = "ls -l \"jason's \\\"folder\\\"\""
	arr = shellParse(cmd)
	if len(arr) != 3 {
		t.Fatalf("Wrong len %d", len(arr))
	}
	if arr[0] != "ls" || arr[1] != "-l" || arr[2] != "jason's \"folder\"" {
		t.Fatalf("Wrong values %v", arr)
	}
}