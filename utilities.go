package main

import (
	"encoding/json"
	"fmt"
)

func StringInSlice(a string, list [13]string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func VerifyIndex(a int) int {
	if a >= 0 {
		return a
	}
	return 0
}

func PrettyPrint(v map[int]*Body) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
