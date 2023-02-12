package main

import (
	"fmt"
	"os"
	sat "sat/pkg"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Errorf("Error: no CNF files supplied")
	}
	filePath := os.Args[1]
	sat.ParseCNFFile(filePath)
}
