package main

import (
	"fmt"
	"math/rand"
	"os"
	sat "sat/pkg"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Errorf("Error: no CNF files supplied")
	}
	filePath := os.Args[1]

	rand.Seed(time.Now().UTC().UnixNano())

	sat.StartTimer()
	formula, formulaState, err := sat.ParseCNFFile(filePath)
	if err != nil {
		fmt.Errorf("Error", err)
		return
	}

	//formula.PrintBooleanFormula()

	sat.DebugLine("~~~Solving...")

	runCounter := 0
	runsBeforeIncreasingBacktrackingRate := 6
	isSat, runOutOfBacktracks, state := false, true, &sat.BooleanFormulaState{}
	formulaState.StateSetUp()

	for runOutOfBacktracks {
		runCounter++
		isSat, runOutOfBacktracks, state = formula.SolveFormula(formulaState.Copy())
		if runOutOfBacktracks && runCounter%runsBeforeIncreasingBacktrackingRate == 0 {
			//Restart!
			newBacktrackLimit := formula.BacktrackingLimit + formula.BacktrackingLimitIncreaseRate
			sat.DebugFormat("~~~Back track limit of %d hit %d times! Starting run #%d with backtracking limit %d \n", formula.BacktrackingLimit, runsBeforeIncreasingBacktrackingRate, runCounter+1, newBacktrackLimit)
			formula.BacktrackingLimit = newBacktrackLimit
		}
	}

	sat.StopTimer()

	if isSat { // For debugging purposes to see whether our assignment even work
		test := state.Debug_CheckAssignmentIsSat()
		sat.DebugLine("is this assignment really sat?", test)
		//sat.PrintBooleanFormulaState(state)
	}
	sat.DebugFormat("Solution: Is sat: %v \n", isSat)
	sat.DebugLine("Number of runs needed: ", runCounter)

	//Standard output for autograder
	fmt.Println("")
	pathTokens := strings.Split(filePath, "/")
	time := sat.GetElapsedNano() / 1000000000
	assignmentString := ""
	//TODO: come back to uncomment this later when it's less annoying to dev (clogs terminal with text)
	// for varIndx, assignment := range state.Assignments {
	// 	assignmentString += fmt.Sprintf("%d %t ", varIndx, assignment == sat.POS)
	// }
	if isSat {
		fmt.Printf("{\"Instance\": \"%s\", \"Time\": \"%.2f\", \"Result\": \"SAT\", \"Solution\": \"%s\"}", pathTokens[len(pathTokens)-1], time, strings.TrimSuffix(assignmentString, " "))
	} else {
		fmt.Printf("{\"Instance\": \"%s\", \"Time\": \"%.2f\", \"Result\": \"UNSAT\"}", pathTokens[len(pathTokens)-1], time)
	}
	fmt.Println("")
}
