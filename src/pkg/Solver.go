package pkg

func (b *BooleanFormula) SolveFormula(initialState *BooleanFormulaState) (bool, bool, *BooleanFormulaState) {

	//Shuffle to add a bit of randomness
	b.ShuffleFormulaVariableBranchingOrder()

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
		//DebugLine("given!", clauseIndex, watchedLiteralHolder)
		initialState.ClauseWatchedLiterals[clauseIndex] = watchedLiteralHolder
	}

	// DebugLine("before solving state")
	// PrintBooleanFormulaState(initialState)
	// DebugLine("map of ClauseWatchedLiterals", initialState.ClauseWatchedLiterals)
	// DebugLine("map of VariablesKeepingTrackOfWhereTheyreBeingWatched", initialState.VariablesKeepingTrackOfWhereTheyreBeingWatched)
	// DebugLine("map of UnitClauses", initialState.UnitClauses)
	// DebugLine("now solve", initialState.Sat)

	//Now solve the state...
	return initialState.SolveFromState()
}

func (state *BooleanFormulaState) SolveFromState() (bool, bool, *BooleanFormulaState) {

	if !state.Sat {
		return false, false, state
	}

	if state.Formula.BacktrackCounter > state.Formula.BacktrackingLimit {
		return false, true, nil
	}

	//Clear out the unit clauses first
	err := state.UnitClauseElimination()
	if !state.Sat || err != nil {
		return false, false, nil
	}

	//Clear out the unit clauses first
	state.PureLiteralElimination()
	if !state.Sat {
		return false, false, nil
	}

	// if there are no more clauses, terminate before branching
	if len(state.DeletedClauses) == len(state.Formula.Clauses) {
		return true, false, state
	}

	// TODO: could make a better heuristic as opposed to arbitrarily picking a variable to branch on
	// loop over the variables in formula for branching
	for _, varIdx := range state.Formula.VarBranchingOrder {
		_, ok := state.Assignments[varIdx]

		if !ok {
			//DebugFormat("DEPTH: %d, Branching On V%d~ \n", len(state.Assignments), varIdx)
			stateCopy := state.Copy()
			assignment := stateCopy.AssignmentFromDynamicLargestCombinedSum(varIdx)
			stateCopy.AssignmentPropagation(varIdx, assignment)
			solved, runOutOfBacktracks, solvedState := stateCopy.SolveFromState()

			if runOutOfBacktracks {
				return false, true, nil
			}

			if solved {
				return solved, false, solvedState
			}

			state.Formula.BacktrackCounter++
			//DebugLine("Backtrack #", state.Formula.BacktrackCounter)

			stateCopy = state.Copy()
			stateCopy.AssignmentPropagation(varIdx, Negate(assignment))
			solved, runOutOfBacktracks, solvedState = stateCopy.SolveFromState()

			if runOutOfBacktracks {
				return false, true, nil
			}

			if solved {
				return solved, false, solvedState
			}

			state.Formula.BacktrackCounter++
			//DebugLine("Backtrack #", state.Formula.BacktrackCounter)
		}

		//This loop will bottom out if we've gone branched on all the
		//unassigned variables and none of them lead to sat
	}

	return false, false, state
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
