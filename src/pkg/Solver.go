package pkg

import (
	"math"
	"math/rand"
	"sort"
)

// Do some operations that shouldn't need to be repeated with every solver reset
// If these already make the formula unsat, it will be caught down the line...
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
	b.CopyShuffledFormulaVariableBranchingOrder(initialState)
	initialState.VarBranchingOrderPointer = &initialState.VarBranchingOrderLocal

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

	//First, since it's time to branch, maybe we shoudl check if we should make a new iteration order...
	if (state.Depth % DEPTH_LIFETIME_FOR_SORTING_ORDERS) == 0 {
		state.VarBranchingOrderLocal = append([]VarIndex(nil), *state.Parent.VarBranchingOrderPointer...)
		state.VarBranchingOrderPointer = &state.VarBranchingOrderLocal
		newVariableScoring := state.ScoreVariablesForNewBranchingOrder()

		sort.Slice(state.VarBranchingOrderLocal, func(i, j int) bool {
			iApprearances := (*newVariableScoring)[state.VarBranchingOrderLocal[i]]
			jApprearances := (*newVariableScoring)[state.VarBranchingOrderLocal[j]]
			return iApprearances > jApprearances
		})
	}

	//This loop will stop at fist unassigned var that actually has remaining instances
	for _, varIdx := range *state.VarBranchingOrderPointer {
		_, ok := state.Assignments[varIdx]

		if !ok {
			stateCopy := state.Copy() //Making a clean copy for the negation later on

			isPure, shouldSkip, assignment := state.AssignmentFromDynamicLargestCombinedSum(varIdx)
			if shouldSkip {
				continue
			}
			state.AssignmentPropagation(varIdx, assignment)
			solved, runOutOfBacktracks, solvedState := state.SolveFromState()

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

			stateCopy.AssignmentPropagation(varIdx, Negate(assignment))
			solved, runOutOfBacktracks, solvedState = stateCopy.SolveFromState()

			if runOutOfBacktracks {
				return false, true, nil
			}

			return solved, false, solvedState
		}
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

	//Skip it if it has no instances
	if posCount == 0 && negCount == 0 {
		return false, true, POS
	}

	if negCount == 0 {
		return true, false, POS

	} else if posCount == 0 {
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

		state.ClauseWatchedLiterals[clauseIndex] = watchedLiteralHolder
	}
}

// TODO: we can also just update our map of pure variables from here
func (state *BooleanFormulaState) ScoreVariablesForNewBranchingOrder() *map[VarIndex]float64 {
	scores := make(map[VarIndex]float64)
	for varIndx, variable := range state.Formula.Vars {
		if _, ok := state.Assignments[varIndx]; !ok {
			continue
		}

		posCount := 0
		negCount := 0

		for clauseIdx, varState := range variable.ClauseAppearances {
			if _, ok := state.DeletedClauses[clauseIdx]; !ok {
				if varState == POS {
					posCount++
				} else {
					negCount++
				}
			}
		}
		scores[varIndx] = math.Max(float64(negCount)/float64(posCount), float64(posCount)/float64(negCount))
		scores[varIndx] += rand.Float64()
	}

	return &scores
}
