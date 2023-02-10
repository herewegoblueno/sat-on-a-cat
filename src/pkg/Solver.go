package pkg

import "fmt"

//This only picks it one at a time, so you want to call this repeatedly
//Note that this is a destructive operation; it removes elements from multiple data structures
func (f *BooleanFormula) UnitClauseElimination() {
	if (len(f.unitClauses)) > 0 {
		unitClauseIndx, otherUnitClauseIndxs := f.unitClauses[0], f.unitClauses[1:]
		unitClause := f.clauses[unitClauseIndx]
		var unitVarIndx VarIndex
		var unitInstanceState VarState

		if len(unitClause.instances) != 1 {
			//TODO: come back to this later
			fmt.Printf("error: non-unit clause being treated as a unit clase!")
		}

		//There's only one instance so this should only run once
		for varIndx, instanceState := range unitClause.instances {
			unitInstanceState = instanceState
			unitVarIndx = varIndx
		}

		delete(f.clauses, unitClauseIndx)
		delete(f.vars[unitVarIndx].clauseAppearances, unitClauseIndx)
		f.UnitPropagation(unitVarIndx, unitInstanceState)
		f.unitClauses = otherUnitClauseIndxs
	}
}

//Note that this can be a destructive operation
//This can also set a formula to unsat
func (f *BooleanFormula) UnitPropagation(indx VarIndex, propagatedState VarState) {
	for clauseIndx, instanceState := range f.vars[indx].clauseAppearances {

		if instanceState == propagatedState {
			//Completely delete the clause
			delete(f.clauses, clauseIndx)
			for _, variable := range f.vars {
				delete(variable.clauseAppearances, clauseIndx) //TODO: is this memory safe?
			}

		} else {
			//Remove the instance from the clause
			delete(f.clauses[clauseIndx].instances, indx)
			if len(f.clauses[clauseIndx].instances) == 0 {
				f.sat = false
			}
		}
	}
}
