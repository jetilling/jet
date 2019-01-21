package main

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
