package pkg

import (
	"fmt"
	"math"
	"math/rand"
)

func (b *BooleanFormula) PrintBooleanFormula() {
	DebugLine("~~Printing Vars~~")
	for _, v := range b.Vars {
		PrintSATVar(v)
	}
	DebugLine("")
	DebugLine("")

	DebugLine("~~Printing Clauses~~")
	for _, c := range b.Clauses {
		PrintSATClause(c)
	}
	DebugLine("")
	DebugLine("")
}

func PrintVarState(v VarState) string {
	if v == POS {
		return "POS"
	} else {
		return "NEG"
	}
}

func PrintSATClause(c *SATClause) {
	DebugFormat("C%v \n", c.Index)
	for varIndx, varState := range c.Instances {
		DebugFormat("V%v: %v   ", varIndx, PrintVarState(varState))
	}
	DebugLine("\n===")
}

func PrintSATVar(v *SATVar) {
	DebugFormat("V%v \n", v.Index)
	for cIndx, varState := range v.ClauseAppearances {
		DebugFormat("C%v: %v   ", cIndx, PrintVarState(varState))
	}
	DebugLine("\n===")
}

func PrintBooleanFormulaState(s *BooleanFormulaState) {
	//For now, just print the assignments...
	irrelevantVariables := "Irrelevant Variables: "
	for varIndx := range s.Formula.Vars {
		if assignmnet, ok := s.Assignments[varIndx]; ok {
			DebugFormat("V%v: %v    ", varIndx, PrintVarState(assignmnet))
		} else {
			irrelevantVariables += fmt.Sprintf("V%v, ", varIndx)
		}
	}
	DebugLine("")
	DebugLine("")
	DebugLine(irrelevantVariables)
	DebugLine("")
}

func (b *BooleanFormulaState) Copy() *BooleanFormulaState {
	new_b := BooleanFormulaState{
		Parent:  b,
		Formula: b.Formula,
		Depth:   b.Depth + 1,
		//VarBranchingOrderLocal is intentionally not copied over
		VarBranchingOrderPointer: b.VarBranchingOrderPointer,
		Assignments:              make(map[VarIndex]VarState),
		ClauseWatchedLiterals:    make(map[ClauseIndex]WatchedLiterals),
		VariablesKeepingTrackOfWhereTheyreBeingWatched: make(map[VarIndex][]ClauseIndex),
		DeletedClauses: make(map[ClauseIndex]bool),
		UnitClauses:    make(map[ClauseIndex]VarIndex),
		PureVariables:  make(map[VarIndex]VarState),
		Sat:            b.Sat,
	}
	for k, v := range b.Assignments {
		new_b.Assignments[k] = v
	}
	for k, v := range b.ClauseWatchedLiterals {
		new_b.ClauseWatchedLiterals[k] = v
	}
	for k, v := range b.VariablesKeepingTrackOfWhereTheyreBeingWatched {
		new_b.VariablesKeepingTrackOfWhereTheyreBeingWatched[k] = v
	}
	for k, v := range b.DeletedClauses {
		new_b.DeletedClauses[k] = v
	}
	for k, v := range b.UnitClauses {
		new_b.UnitClauses[k] = v
	}
	for k, v := range b.PureVariables {
		new_b.PureVariables[k] = v
	}
	return &new_b
}

func Negate(v VarState) VarState {
	if v == POS {
		return NEG
	} else {
		return POS
	}
}

func (formula *BooleanFormula) CopyShuffledFormulaVariableBranchingOrder(state *BooleanFormulaState) {

	//Copying over
	state.VarBranchingOrderLocal = append([]VarIndex(nil), formula.VarBranchingOrderOriginal...)

	if VAR_BRANCHING_ORDER_SHUFFLE_CHANCE == 0 || VAR_BRANCHING_ORDER_SHUFFLE_DISTANCE == 0 {
		return
	}

	shuffleRadius := int(math.Floor(float64(len(formula.Vars)) * VAR_BRANCHING_ORDER_SHUFFLE_DISTANCE))
	if shuffleRadius == 0 {
		return
	}

	shuffleCounter := 0
	for index := range state.VarBranchingOrderLocal {
		if rand.Intn(100) > VAR_BRANCHING_ORDER_SHUFFLE_CHANCE {
			continue
		}

		offset := (rand.Intn(shuffleRadius*2) - shuffleRadius)
		newIndex := index + offset

		//Making sure its in bounds...
		if (newIndex) >= len(state.VarBranchingOrderLocal) {
			newIndex = index - offset
		}

		if (newIndex) < 0 {
			continue
		}

		shuffleCounter++
		state.VarBranchingOrderLocal[index], state.VarBranchingOrderLocal[newIndex] = state.VarBranchingOrderLocal[newIndex], state.VarBranchingOrderLocal[index]
	}
}

func Clamp(num int, max int, min int) int {
	if num > max {
		return max
	}
	if num < min {
		return min
	}
	return num
}
