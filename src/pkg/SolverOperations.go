package pkg

import "fmt"

//Does all of it at once
func (s *BooleanFormulaState) UnitClauseElimination() error {
	for len(s.UnitClauses) != 0 && s.Sat {
		//Pick one at a time
		var unitClauseIndx ClauseIndex
		for k, _ := range s.UnitClauses {
			unitClauseIndx = k
			break
		}
		delete(s.UnitClauses, unitClauseIndx)
		unitClause := s.Formula.Clauses[unitClauseIndx]

		if len(unitClause.Instances) != 1 {
			return fmt.Errorf("error: non-unit clause being treated as a unit clause!")
		}

		var unitVarIndx VarIndex
		var unitInstanceState VarState

		//There's only one instance so this should only run once
		for varIndx, instanceState := range unitClause.Instances {
			unitInstanceState = instanceState
			unitVarIndx = varIndx
		}

		s.DeletedClauses[unitClauseIndx] = true
		delete(s.PureVariables, unitVarIndx)
		s.AssignmentPropagation(unitVarIndx, unitInstanceState)
	}
	return nil
}

//Note that this can be a destructive operation
//This can also set a formula to unsat
func (s *BooleanFormulaState) AssignmentPropagation(newlyAsgnVar VarIndex, propagatedState VarState) {

	s.Assignments[newlyAsgnVar] = propagatedState

	for clauseIndx, instanceState := range s.Formula.Vars[newlyAsgnVar].ClauseAppearances {
		if _, ok := s.DeletedClauses[clauseIndx]; !ok {
			continue
		}
		if instanceState == propagatedState {
			s.DeletedClauses[clauseIndx] = true
		} else {
			//If it's a unit clause and there's a mismatch, then we're unsat
			if _, ok := s.UnitClauses[clauseIndx]; ok {
				s.Sat = false
				return
			}
		}
	}

	//check watch literals of other clauses
	for _, clauseIndx := range s.VariablesKeepingTrackOfWhereTheyreBeingWatched[newlyAsgnVar] {
		var otherWatchedLiteral VarIndex
		var isRight bool

		if _, ok := s.DeletedClauses[clauseIndx]; ok {
			continue
		}

		if s.ClauseWatchedLiterals[clauseIndx].left == newlyAsgnVar {
			isRight = false
			otherWatchedLiteral = s.ClauseWatchedLiterals[clauseIndx].right
		} else {
			isRight = true
			otherWatchedLiteral = s.ClauseWatchedLiterals[clauseIndx].left
		}

		//Pick another watcher to replace!
		foundReplacementWatchedLiteral := false
		clause := s.Formula.Clauses[clauseIndx]
		for watchLiteralCandidateIndx, _ := range clause.Instances {
			// make sure this candidate is not already assigned and that it's not the other watched literal
			_, isAssigned := s.Assignments[watchLiteralCandidateIndx]
			if !isAssigned && watchLiteralCandidateIndx != otherWatchedLiteral {
				// can be used as watch literal
				if isRight {
					s.ClauseWatchedLiterals[clauseIndx] = WatchedLiterals{right: watchLiteralCandidateIndx, left: otherWatchedLiteral}
				} else {
					s.ClauseWatchedLiterals[clauseIndx] = WatchedLiterals{left: watchLiteralCandidateIndx, right: otherWatchedLiteral}
				}
				s.VariablesKeepingTrackOfWhereTheyreBeingWatched[watchLiteralCandidateIndx] = append(s.VariablesKeepingTrackOfWhereTheyreBeingWatched[watchLiteralCandidateIndx], clauseIndx)
				foundReplacementWatchedLiteral = true
				break
			}
		}

		if !foundReplacementWatchedLiteral { //That means this clause has become a unit clause! (┛◉Д◉)┛彡┻━┻
			delete(s.ClauseWatchedLiterals, clauseIndx)
			s.UnitClauses[clauseIndx] = true
		}
	}

	delete(s.VariablesKeepingTrackOfWhereTheyreBeingWatched, newlyAsgnVar)
}

//TODO: we don't currently handle state well such that
//we can easily detect new pure literals AFTER initial parsing
func (s *BooleanFormulaState) PureLiteralElimination() {
	for varIndx, varState := range s.PureVariables {
		if !s.Sat {
			return
		}
		delete(s.PureVariables, varIndx)
		s.AssignmentPropagation(varIndx, varState)
	}
}
