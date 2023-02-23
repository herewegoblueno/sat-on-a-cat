package main

import (
	"fmt"
	"math"
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

	runCounterTotal := 0
	runsCounterSinceLastRestart := 0

	runsBeforeIncreasingBacktrackingLimitCurrent := sat.RUNS_BEFORE_INCREASING_BACKTRACK_LIMIT_MAX
	currentBacktrackingLimitIncrement := sat.BACKTRACK_LIMIT_INCREMENT_MIN

	isSat, runOutOfBacktracks, state := false, true, &sat.BooleanFormulaState{}
	formulaState.StateSetUp()

	for runOutOfBacktracks {
		runCounterTotal++
		runsCounterSinceLastRestart++
		isSat, runOutOfBacktracks, state = formula.SolveFormula(formulaState.Copy())
		if runOutOfBacktracks && runsCounterSinceLastRestart == runsBeforeIncreasingBacktrackingLimitCurrent {
			//Restart with a new limit!
			newBacktrackLimit := formula.BacktrackingLimit + currentBacktrackingLimitIncrement
			sat.DebugFormat("~~~Back track limit of %d hit %d times! Starting run #%d with backtracking limit %d (increase of %d)\n",
				formula.BacktrackingLimit, runsBeforeIncreasingBacktrackingLimitCurrent, runCounterTotal+1, newBacktrackLimit, currentBacktrackingLimitIncrement)
			formula.BacktrackingLimit = newBacktrackLimit
			runsCounterSinceLastRestart = 0

			runsBeforeIncreasingBacktrackingLimitCurrent = sat.Clamp(
				int(math.Floor(float64(runsBeforeIncreasingBacktrackingLimitCurrent)*sat.DECREASE_RATE_FOR_RUNS_THRESHOLD)),
				sat.RUNS_BEFORE_INCREASING_BACKTRACK_LIMIT_MAX,
				sat.RUNS_BEFORE_INCREASING_BACKTRACK_LIMIT_MIN,
			)

			currentBacktrackingLimitIncrement = int(math.Floor(float64(currentBacktrackingLimitIncrement) * sat.BACKTRACK_LIMIT_INCREMENT_INCREASE_RATE))

		}
	}

	sat.StopTimer()

	if isSat { // For debugging purposes to see whether our assignment even work
		test := state.Debug_CheckAssignmentIsSat()
		sat.DebugLine("is this assignment really sat?", test)
		//sat.PrintBooleanFormulaState(state)
	}
	sat.DebugFormat("Solution: Is sat: %v \n", isSat)
	sat.DebugLine("Number of runs needed: ", runCounterTotal)

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
