package pkg

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func ParseCNFFile(filename string) (BooleanFormula, error) {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file) // TODO: error on > 65536 chars
	scanner.Split(bufio.ScanLines)    // split on lines

	clauseNum := 0
	currFormula := BooleanFormula{
		make(map[VarIndex]SATVar),
		make(map[ClauseIndex]SATClause),
		true,
		[]ClauseIndex{},
	}

	for scanner.Scan() { //Line by line...
		tokens := strings.Split(scanner.Text(), " ")

		if len(tokens) > 0 && tokens[0] != "p" && tokens[0] != "c" {

			currClause := SATClause{
				ClauseIndex(clauseNum),
				make(map[VarIndex]VarState),
			}
			for _, tok := range tokens {
				// skip at the end of the line, don't bother converting to int
				if tok == "0" {
					continue
				}

				// convert var string to num
				var_as_num, err := strconv.Atoi(tok)
				if err != nil {
					fmt.Println("error while converting tok to num")
					break // TODO: standardize error
				}

				var currVar SATVar
				currVarIndex := VarIndex(math.Abs(float64(var_as_num)))

				existingVariable, ok := currFormula.vars[currVarIndex]
				if ok {
					currVar = existingVariable
				} else {
					currVar = SATVar{
						VarIndex(currVarIndex),
						true,
						make(map[ClauseIndex]VarState),
						DEFAULT, //Will get overwritten if it's wrong anyways
					}
				}

				currFormula.vars[currVarIndex] = currVar
				var newState VarState

				if var_as_num > 0 {
					newState = POS
				} else {
					newState = NEG
				}

				//If the variable is already in this clause and has an opposite value
				//then just remove it from the clause entirely
				previousAppearanceInCurrentClause, ok := existingVariable.clauseAppearances[currClause.index]
				if ok && previousAppearanceInCurrentClause != newState {
					//Remove the variable from the clause
					delete(currClause.instances, currVarIndex)
					delete(currVar.clauseAppearances, currClause.index)
				} else {
					currClause.instances[currVarIndex] = newState
					currVar.clauseAppearances[currClause.index] = newState

					//Checking if still pure!
					if currVar.lastSeenState == DEFAULT {
						currVar.lastSeenState = newState
					} else if currVar.isPure && currVar.lastSeenState != newState {
						currVar.isPure = false
					} else if currVar.isPure { //No point in keeping track if we know it's not pure anymore
						currVar.lastSeenState = newState
					}
				}

			}

			//Only add the clause if it's not empty
			if len(currClause.instances) > 0 {
				currFormula.clauses[ClauseIndex(clauseNum)] = currClause
			}
			clauseNum += 1
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return currFormula, nil
}
