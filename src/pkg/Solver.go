package pkg

//Do some operations that shouldn't need to be repeated with every solver reset
//If these already make the formula unsat, it will be caught down the line...
func (state *BooleanFormulaState) StateSetUp() {
	state.SetWatcherVariables()

	//Clear out the unit clauses first
	err := state.UnitClauseElimination()
	if !state.Sat || err != nil {
		return
	}

	state.PureLiteralElimination()
	if !state.Sat {
		return
	}

	// if there are no more clauses, terminate before branching
	if len(state.DeletedClauses) == len(state.Formula.Clauses) {
		return
	}
}

func (b *BooleanFormula) SolveFormula(initialState *BooleanFormulaState) (bool, bool, *BooleanFormulaState) {

	//Shuffle to add a bit of randomness
	b.ShuffleFormulaVariableBranchingOrder()

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

	state.PureLiteralElimination()
	if !state.Sat {
		return false, false, nil
	}

	// if there are no more clauses, terminate before branching
	if len(state.DeletedClauses) == len(state.Formula.Clauses) {
		return true, false, state
	}

	for _, varIdx := range state.Formula.VarBranchingOrderShuffled {
		_, ok := state.Assignments[varIdx]

		if !ok {
			//DebugFormat("DEPTH: %d, Branching On V%d~ \n", len(state.Assignments), varIdx)
			stateCopy := state.Copy()
			isPure, shouldSkip, assignment := stateCopy.AssignmentFromDynamicLargestCombinedSum(varIdx)
			if shouldSkip {
				continue
			}
			stateCopy.AssignmentPropagation(varIdx, assignment)
			solved, runOutOfBacktracks, solvedState := stateCopy.SolveFromState()

			if runOutOfBacktracks {
				return false, true, nil
			}

			if solved {
				return solved, false, solvedState
			}

			if isPure { //No point in continuing if it was pure and the initial assignment didn't work...
				continue
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

func (state *BooleanFormulaState) AssignmentFromDynamicLargestCombinedSum(variable VarIndex) (bool, bool, VarState) {
	posCount := 0
	negCount := 0

	for clauseIdx, varState := range state.Formula.Vars[variable].ClauseAppearances {
		if _, ok := state.DeletedClauses[clauseIdx]; !ok {
			if varState == POS {
				posCount += 1
			} else {
				negCount += 1
			}
		}
	}

	if posCount == 0 && negCount == 0 {
		return false, true, POS
	}

	//If we realize this is pure, great news! Even if we don't end up finding sat, we can
	//net our future sibling states get this information by passing it into our parent...
	if negCount == 0 {
		DebugFormat("Discovered that V%d is %v pure throughout %d instances!\n", variable, int(POS), posCount)
		state.Parent.PureVariables[variable] = POS
		return true, false, POS

	} else if posCount == 0 {
		DebugFormat("Discovered that V%d is %v pure throughout %d instances!\n", variable, int(NEG), negCount)
		state.Parent.PureVariables[variable] = NEG
		return true, false, NEG
	}

	if posCount > negCount {
		return false, false, POS
	} else {
		return false, false, NEG
	}
}

func (state *BooleanFormulaState) SetWatcherVariables() {
	//First assign everyone who has a sign larger than 1 some watched variables...
	for clauseIndex, clause := range state.Formula.Clauses {

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
				state.VariablesKeepingTrackOfWhereTheyreBeingWatched[varIndex] = append(state.VariablesKeepingTrackOfWhereTheyreBeingWatched[varIndex], clauseIndex)
			} else if counter == 2 {
				watchedLiteralHolder.left = varIndex
				state.VariablesKeepingTrackOfWhereTheyreBeingWatched[varIndex] = append(state.VariablesKeepingTrackOfWhereTheyreBeingWatched[varIndex], clauseIndex)
			} else {
				break
			}
		}
		//DebugLine("given!", clauseIndex, watchedLiteralHolder)
		state.ClauseWatchedLiterals[clauseIndex] = watchedLiteralHolder
	}
}
