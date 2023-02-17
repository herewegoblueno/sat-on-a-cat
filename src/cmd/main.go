package main

import (
	"fmt"
	"os"
	sat "sat/pkg"
	"strings"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Errorf("Error: no CNF files supplied")
	}
	filePath := os.Args[1]

	sat.StartTimer()
	formula, formulaState, err := sat.ParseCNFFile(filePath)
	if err != nil {
		fmt.Errorf("Error", err)
	}

	formula.PrintBooleanFormula()
	sat.DebugLine("Solving...")

	isSat, state := formula.SolveFormula(formulaState)
	sat.StopTimer()

	if isSat { // For debugging purposes to see whether our assignment even work
		test := state.Debug_CheckAssignmentIsSat()
		sat.DebugLine("is this assignment really sat?", test)
		sat.PrintBooleanFormulaState(state)
	}
	sat.DebugFormat("Solution: Is sat: %v \n", isSat)

	//Standard output for autograder
	fmt.Println("")
	pathTokens := strings.Split(filePath, "/")
	time := sat.GetElapsedNano() / 1000000000
	assignmentString := ""
	for varIndx, assignment := range state.Assignments {
		assignmentString += fmt.Sprintf("%d %t ", varIndx, assignment == sat.POS)
	}
	if isSat {
		fmt.Printf("{\"Instance\": \"%s\", \"Time\": \"%.2f\", \"Result\": \"SAT\", \"Solution\": \"%s\"}", pathTokens[len(pathTokens)-1], time, strings.TrimSuffix(assignmentString, " "))
	} else {
		fmt.Printf("{\"Instance\": \"%s\", \"Time\": \"%.2f\", \"Result\": \"UNSAT\"}", pathTokens[len(pathTokens)-1], time)
	}
	fmt.Println("")
}
