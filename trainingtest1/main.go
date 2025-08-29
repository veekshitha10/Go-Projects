package main

import "fmt"

func main() {
	str2 := "veek3Shit78ha"

	for _, v := range str2 {
		if v >= 97 && v <= 122 {
			fmt.Print(string(v - 32))

		} else {
			fmt.Print(string(v))
		}
	}
}
