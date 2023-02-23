package test

import (
	"fmt"
	sat "sat/pkg"
	"testing"
)

func BenchmarkSolving(b *testing.B) {
	formula, formulaState, err := sat.ParseCNFFile("../../input/C181_3151.cnf")
	if err != nil {
		fmt.Errorf("Error %v", err)
		return
	}

	runCounter := 0
	runsBeforeIncreasingBacktrackingLimit := 50
	isSat, runOutOfBacktracks, state := false, true, &sat.BooleanFormulaState{}
	formulaState.StateSetUp()

	for runOutOfBacktracks {
		runCounter++
		if runCounter > 150 { //TODO: maybe remove this later
			break
		}
		isSat, runOutOfBacktracks, state = formula.SolveFormula(formulaState.Copy())
		if runOutOfBacktracks && runCounter%runsBeforeIncreasingBacktrackingLimit == 0 {
			//Restart!
			newBacktrackLimit := formula.BacktrackingLimit + formula.BacktrackingLimitIncreaseRate
			sat.DebugFormat("~~~Back track limit of %d hit %d times! Starting run #%d with backtracking limit %d \n", formula.BacktrackingLimit, runsBeforeIncreasingBacktrackingLimit, runCounter+1, newBacktrackLimit)
			formula.BacktrackingLimit = newBacktrackLimit
		}
	}

	if isSat { // For debugging purposes to see whether our assignment even work
		test := state.Debug_CheckAssignmentIsSat()
		sat.DebugLine("is this assignment really sat?", test)
		//sat.PrintBooleanFormulaState(state)
	}
	sat.DebugFormat("Solution: Is sat: %v \n", isSat)
	sat.DebugLine("Number of runs needed: ", runCounter)
}
