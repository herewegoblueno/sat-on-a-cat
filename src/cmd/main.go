package main

import (
	"os"
	sat "sat/pkg"
)

func main() {
	filePath := os.Args[1]
	sat.ParseCNFFile(filePath)
}
