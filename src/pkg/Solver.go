package pkg

import "fmt"

func (b *BooleanFormula) SolveFormula(initialState *BooleanFormulaState) (bool, *BooleanFormulaState) {
	//First assign everyone who has a sign larger than 1 some watched variables...
	for clauseIndex, clause := range b.Clauses {

		if len(clause.Instances) < 2 {
			continue
		}

		watchedLiteralHolder := WatchedLiterals{}
		//Ugly way to get two random keys...
		counter := 0
		for varIndex := range clause.Instances {
			counter++
			if counter == 1 {
				watchedLiteralHolder.right = varIndex
				initialState.VariablesKeepingTrackOfWhereTheyreBeingWatched[varIndex] = append(initialState.VariablesKeepingTrackOfWhereTheyreBeingWatched[varIndex], clauseIndex)
			} else if counter == 2 {
				watchedLiteralHolder.left = varIndex
				initialState.VariablesKeepingTrackOfWhereTheyreBeingWatched[varIndex] = append(initialState.VariablesKeepingTrackOfWhereTheyreBeingWatched[varIndex], clauseIndex)
			} else {
				break
			}
		}
		initialState.ClauseWatchedLiterals[clauseIndex] = watchedLiteralHolder
	}

	fmt.Println("before solving state")
	PrintBooleanFormulaState(initialState)
	fmt.Println("map of ClauseWatchedLiterals", initialState.ClauseWatchedLiterals)
	fmt.Println("map of VariablesKeepingTrackOfWhereTheyreBeingWatched", initialState.VariablesKeepingTrackOfWhereTheyreBeingWatched)
	fmt.Println("map of UnitClauses", initialState.UnitClauses)
	fmt.Println("now solve", initialState.Sat)

	//Now solve the state...
	solveable, solveableState := initialState.SolveFromState()
	if solveable {
		test := solveableState.CheckAssignmentIsSat() // For debugging purposes to see whether our assignment even work
		fmt.Println("is this sat?", test)
	}
	return solveable, solveableState
}

func (state *BooleanFormulaState) SolveFromState() (bool, *BooleanFormulaState) {
	//Clear out the unit clauses first
	if !state.Sat {
		return false, state
	}

	err := state.UnitClauseElimination()
	if !state.Sat || err != nil {
		return false, nil
	}

	//Clear out the unit clauses first
	state.PureLiteralElimination()
	if !state.Sat {
		return false, nil
	}

	// if there are no more clauses, terminate before branching
	if len(state.DeletedClauses) == len(state.Formula.Clauses) {
		return true, state
	}

	positiveState := state.Copy()
	// TODO: could make a better heuristic as opposed to arbitrarily picking a variable to branch on
	// loop over the variables in formula
	for idx, _ := range state.Formula.Vars {
		// fmt.Println("examining candidate", idx)
		// check if the variable is assigned
		fmt.Println("looking for sat on this", state.Assignments, idx)
		_, ok := positiveState.Assignments[idx]
		// if not assigned
		if !ok {
			// fmt.Println("examining candidate that is not assigned", idx)
			assignment := positiveState.AssignmentFromDynamicLargestCombinedSum(idx)

			fmt.Println("guessing var", idx, assignment)

			positiveState.AssignmentPropagation(idx, assignment)
			solved, solvedState := positiveState.SolveFromState()

			fmt.Println("neg guessing var", idx, assignment, solved)

			if solved {
				return solved, solvedState
			} else {
				negativeState := state.Copy()
				negativeState.AssignmentPropagation(idx, Negate(assignment))
				fmt.Println("unsat on both, look for other vars", state.Assignments, idx)
				return negativeState.SolveFromState()
			}
		}
	}

	return false, state
}

func (state *BooleanFormulaState) CheckAssignmentIsSat() bool {
	for clauseIdx, clause := range state.Formula.Clauses {
		satisfiedClause := false
		for literal, asgn := range clause.Instances {
			if state.Assignments[literal] == asgn {
				// found one satisfying
				satisfiedClause = true
				break
			}
		}
		if !satisfiedClause {
			// curr assignments don't satisfy clause
			fmt.Println("error: assignment doesn't satisfy", clauseIdx, state.Formula.Clauses[clauseIdx])
			return false
		}
	}
	// if it reaches here, every clause is satisfied, so it is SAT
	return true
}

func (state *BooleanFormulaState) AssignmentFromDynamicLargestCombinedSum(variable VarIndex) VarState {
	posCount := 0
	negCount := 0
	for clauseIdx, clause := range state.Formula.Clauses {
		if _, ok := state.DeletedClauses[clauseIdx]; !ok {
			if clause.Instances[variable] == POS {
				posCount += 1
			} else {
				negCount += 1
			}
		}
	}
	if posCount > negCount {
		return POS
	} else {
		return NEG
	}
}
