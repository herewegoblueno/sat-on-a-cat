package pkg

import "fmt"

func DebugLine(s ...interface{}) {
	fmt.Println(s...)
}

func DebugFormat(s string, p ...interface{}) {
	fmt.Printf(s, p...)
}
