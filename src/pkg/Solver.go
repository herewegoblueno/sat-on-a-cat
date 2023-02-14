package pkg

func (b *BooleanFormula) SolveFormula(initialState *BooleanFormulaState) (bool, *BooleanFormulaState) {
	//First assign everyone who has a sign larger than 1 some watched variables...
	for clauseIndex, clause := range b.Clauses {
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
	}

	//Now solve the state...
	return initialState.SolveFromState()
}

func (state *BooleanFormulaState) SolveFromState() (bool, *BooleanFormulaState) {
	//Clear out the unit clauses first
	error := state.UnitClauseElimination()
	if !state.Sat || error != nil {
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

	// TODO: could make a better heuristic as opposed to arbitrarily picking a variable to branch on
	// loop over the variables in formula
	for idx, _ := range state.Formula.Vars {
		// check if the variable is assigned
		_, ok := state.Assignments[idx]
		// if not assigned
		if !ok {
			assignment := state.AssignmentFromDynamicLargestCombinedSum(idx)
			copyOfState := state.Copy()

			state.AssignmentPropagation(idx, assignment)
			solved, solvedState := state.SolveFromState()

			if solved {
				return solved, solvedState
			} else {
				copyOfState.AssignmentPropagation(idx, Negate(assignment))
				return copyOfState.SolveFromState()
			}
		}
	}

	//TODO: I think this should be unreachable...error here?
	return true, state
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
