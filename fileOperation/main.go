package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	fl := New("sample.txt")
	value, erro := fl.Read()
	m := make(map[string]int)
	if erro != nil {
		print(erro.Error())
	} else {
		//	fmt.Println(value)
		slice := convertToslice(value)

		for _, v := range slice {
			m[v] = countLength(slice, v)
		}
		for i, v := range m {
			fmt.Println("  ", i, " :", v)
		}
		max, str := getMaxCount(m)
		fmt.Println("max count is ", max, "str is ", str)
	}

}
func countLength(s []string, str string) int {
	i := 0
	for _, v := range s {
		if str == v {
			i++
		}

	}
	return i
}

func convertToslice(s string) []string {
	var slice []string
	str := ""
	for _, r := range s {
		if r == ' ' || r == '.' || r == '\n' {
			if str != "" {
				slice = append(slice, str)
				str = ""

			}
		} else {
			str = str + string(r)

		}

	}
	return slice

}

type FileWriter struct {
	name string
}

func New(name string) *FileWriter {
	return &FileWriter{name: name}
}
func (f *FileWriter) Read() (string, error) {
	file, erro := os.OpenFile(f.name, os.O_RDONLY, 0444)
	if erro != nil {
		return "", erro
	} else {
		defer file.Close()
		res, err := io.ReadAll(file)
		if err != nil {
			return "", err

		} else {
			str := string(res)
			return str, nil
		}
	}

}       

func getMaxCount(m map[string]int) (int, string) {
	maxV := 0
	str := ""
	for K, v := range m {
		if maxV < v {
			maxV = v
			str = K
		}
	}
	return maxV, str
}
