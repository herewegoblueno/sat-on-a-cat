package main

import (
	sat "sat/pkg"
)

func main() {
	sat.ParseCNFFile("toy_simple.cnf")
}
