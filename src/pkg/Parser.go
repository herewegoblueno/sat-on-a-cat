package pkg

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// TODO: non deterministic parsing dropping vars somehow?

func ParseCNFFile(filename string) (*BooleanFormula, *BooleanFormulaState, error) {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file) // TODO: error on > 65536 chars
	scanner.Split(bufio.ScanLines)    // split on lines

	currFormula := BooleanFormula{
		make(map[VarIndex]*SATVar),
		make(map[ClauseIndex]*SATClause),
		make([]VarIndex, 0),
		0.7,  //VarBranchingOrderShuffleDistance[0, 1]
		50,   //VarBranchingOrderShuffleChance[0, 100]
		0,    //BacktrackCounter
		500,  //BacktrackingLimit
		1000, //BacktrackingLimitIncreaseRate
	}

	initialState := BooleanFormulaState{
		&currFormula,
		make(map[VarIndex]VarState),
		make(map[ClauseIndex]WatchedLiterals),
		make(map[VarIndex][]ClauseIndex),
		make(map[ClauseIndex]bool),
		make(map[ClauseIndex]VarIndex),
		make(map[VarIndex]VarState),
		true,
	}

	clauseNum := 0
	for scanner.Scan() { //Line by line...
		tokens := strings.Split(scanner.Text(), " ")

		if len(tokens) > 0 && tokens[0] != "p" && tokens[0] != "c" {

			currClause := SATClause{
				ClauseIndex(clauseNum),
				make(map[VarIndex]VarState),
			}
			for _, tok := range tokens {
				// skip at the end of the line, don't bother converting to int
				if tok == "0" || tok == "" {
					continue
				}

				// convert var string to num
				var_as_num, err := strconv.Atoi(tok)
				if err != nil {
					return nil, nil, fmt.Errorf("error while converting tok to num")
				}

				var currVar *SATVar
				currVarIndex := VarIndex(math.Abs(float64(var_as_num)))

				existingVariable, ok := currFormula.Vars[currVarIndex]
				if ok {
					currVar = existingVariable
				} else {
					currVar = &SATVar{
						VarIndex(currVarIndex),
						make(map[ClauseIndex]VarState),
						DEFAULT, //Will get overwritten if it's wrong anyways
					}
					initialState.PureVariables[currVarIndex] = DEFAULT
					currFormula.Vars[currVarIndex] = currVar
					currFormula.VarBranchingOrder = append(currFormula.VarBranchingOrder, VarIndex(currVarIndex))
				}

				var newState VarState

				if var_as_num > 0 {
					newState = POS
				} else {
					newState = NEG
				}

				//If the variable is already in this clause and has an opposite value
				//then just remove it from the clause entirely
				previousAppearanceInCurrentClause, ok := currVar.ClauseAppearances[currClause.Index]
				if ok && previousAppearanceInCurrentClause != newState {
					//Stop parsing the whole clause if there's both positive and negative in the same clause
					//Because this clause is essentially unconstrained (there's not way it can be unsatisfied)
					//TODO: Haven't found a way for it to re-check the purity of the variables it used to contain
					for varIndx := range currClause.Instances {
						delete(currFormula.Vars[varIndx].ClauseAppearances, currClause.Index)
					}
					//Make the map empty to act like it's empty (code down below will handle that gracefully), then end
					currClause.Instances = make(map[VarIndex]VarState)
					break
				} else {
					currClause.Instances[currVarIndex] = newState
					currVar.ClauseAppearances[currClause.Index] = newState

					_, wasPureLastTimeRead := initialState.PureVariables[currVarIndex]
					if currVar.LastSeenState == DEFAULT {
						currVar.LastSeenState = newState
						initialState.PureVariables[currVarIndex] = newState
					} else if wasPureLastTimeRead && currVar.LastSeenState != newState {
						delete(initialState.PureVariables, currVarIndex)
					} else if wasPureLastTimeRead { //No point in keeping track if we know it's not pure anymore
						currVar.LastSeenState = newState
					}
				}
			}

			//Only add the clause if it's not empty
			clauseLength := len(currClause.Instances)
			if clauseLength > 0 {
				currFormula.Clauses[ClauseIndex(clauseNum)] = &currClause
				if clauseLength == 1 {
					//Pick one (the only, actually) variable
					for varIndx := range currClause.Instances {
						initialState.UnitClauses[ClauseIndex(clauseNum)] = varIndx
						break
					}

				}
			}
			clauseNum += 1
		}
	}

	//Set the brancing order based on the number of clauses that variables come in
	//(so we branch on variables that matter more)
	sort.Slice(currFormula.VarBranchingOrder, func(i, j int) bool {
		iApprearances := len(currFormula.Vars[currFormula.VarBranchingOrder[i]].ClauseAppearances)
		jApprearances := len(currFormula.Vars[currFormula.VarBranchingOrder[j]].ClauseAppearances)
		return iApprearances > jApprearances
	})

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return &currFormula, &initialState, nil
}
