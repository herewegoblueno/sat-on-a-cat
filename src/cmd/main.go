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
	formula, formulaState, err := sat.ParseCNFFile(filePath)
	if err != nil {
		fmt.Errorf("Error", err)
	}

	formula.PrintBooleanFormula()
	fmt.Println("Solving...")
	isSat, state := formula.SolveFormula(formulaState)
	fmt.Printf("Solution: Is sat: %v \n", isSat)
	if isSat {
		sat.PrintBooleanFormulaState(state)
	}
}
