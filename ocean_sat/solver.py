#!/bin/python3
import sys
from copy import copy, deepcopy
from collections import Counter
import random

# Feel free to change the provided types and parsing code to match
# your preferred representation of formulas, clauses, and literals.

class Literal:
    def __init__(self, name, sign):
        self.name = name  # integer
        self.sign = sign  # boolean

    def __repr__(self):
        return ("-" if not self.sign else "") + self.name

    def __eq__(self, other):
        if type(other) != Literal:
            return False
        return self.name == other.name and self.sign == other.sign

    def __hash__(self):
      return hash((self.name, self.sign))


class Clause:
    def __init__(self, id, literalSet):
        self.id = id
        self.literalSet = literalSet

    def __repr__(self):
        return f"{self.id}: {str(self.literalSet)}"

    def __eq__(self, other):
        if type(other) != Clause:
            return False
        return self.id == other.id


# Read and parse a cnf file, returning the variable set and clause set
def readInput(cnfFile):
    variableSet = []
    clauseSet = []
    nextCID = 0
    with open(cnfFile, "r") as f:
        for line in f.readlines():
            tokens = line.strip().split()
            if tokens and tokens[0] != "p" and tokens[0] != "c":
                literalSet = []
                for lit in tokens[:-1]:
                    sign = lit[0] != "-"
                    variable = lit.strip("-")

                    literalSet.append(Literal(variable, sign))
                    if variable not in variableSet:
                        variableSet.append(variable)

                clauseSet.append(Clause(nextCID, literalSet))
                nextCID += 1

    return variableSet, clauseSet

def unitClauseElim(varbset, varAssignment, formula):
    cleaned_formula = deepcopy(formula)

    # remove all clauses that contain any literals that's assigned true
    for clause in formula:
        for cliteral in clause.literalSet:
            if cliteral in varAssignment:
                cleaned_formula.remove(clause)
                break

    # remove all literals that are assigned false
    for clause in cleaned_formula:
        clause.literalSet = [cliteral for cliteral in clause.literalSet if cliteral.name in varbset]

    new_formula = []

    for clause in cleaned_formula:

        # found a unit clause
        if len(clause.literalSet) == 1:
            unit_literal = clause.literalSet[0]
            reverse_unit_literal = Literal(unit_literal.name, not(unit_literal.sign))

            # assign the value to unit_literal
            if unit_literal.sign:
                varAssignment.append(Literal(unit_literal.name, True))
            else:
                varAssignment.append(Literal(unit_literal.name, False))

            for index, check_clause in enumerate(cleaned_formula):
                check_literal_list = check_clause.literalSet
                if unit_literal in check_literal_list:
                    # unit literal in the clause
                    continue # remove clause by adding nothing to new_formula
                elif reverse_unit_literal in check_literal_list:
                  	# remove a specific literal from the clause
                    new_check_clause = deepcopy(check_clause)
                    new_check_clause.literalSet.remove(reverse_unit_literal)
                    new_formula.append(new_check_clause)
                else:
                  	# unit_literal is not in this clause, keep as it is
                    new_formula.append(check_clause)

            # remove unit literal from varbset now that it's assigned
            varbset.remove(unit_literal.name)

            # finished removing one unit_literal
            return unitClauseElim(varbset, varAssignment, new_formula)

    # didn't find any unit_literal
    return varbset, varAssignment, cleaned_formula


def pureLiteralElim(varbset, varAssignment, formula):
    for variable in varbset:
        to_be_removed = []
        positive = False
        negative = False

        for idx, clause in enumerate(formula):
            if Literal(variable, True) in clause.literalSet:
                positive = True
                to_be_removed.append(idx)

            if Literal(variable, False) in clause.literalSet:
                negative = True
                to_be_removed.append(idx)

        if not (positive and negative):
            # variable is pure

            formula = [i for i in formula if formula.index(i) not in to_be_removed]
            # remove all clauses containing +/-x

            if positive:
                varAssignment.append(Literal(variable, True))
            else:
                varAssignment.append(Literal(variable, False))

            varbset.remove(variable)# remove the variable from varbset
            return pureLiteralElim(varbset, varAssignment, formula)

    # didn't find any pure literal
    return varbset, varAssignment, formula

# helper function checks whether a formula is sat given the varAssignment
def checkSat(varAssignment, formula):
    for clause in formula:
        oneTrue = False
        for cliteral in clause.literalSet:
            if cliteral in varAssignment:
                oneTrue = True
        if oneTrue:
            continue
        else:
            return False # unsat
    return True # sat

def solve(varbset, varAssignment, formula):
    varbset, new_varAssignment, new_formula = unitClauseElim(varbset, varAssignment, formula)

    varbset, new_varAssignment, new_formula = pureLiteralElim(varbset, new_varAssignment, new_formula)


    # unsat if the formula contains empty clause, so none of the assignment counts, returning [] as varAssignment
    for clause in new_formula:
        if clause.literalSet == []:
            return varbset, [], new_formula

    # return current varAssignment if the formula has no clause
    if new_formula == []:
        return varbset, varAssignment, new_formula

    if len(varbset) != 0:
        x = varbset[0] # pick the first var in varbSet
        varAssignment.append(Literal(x, True))
        del varbset[0]
        varAssignment_before = deepcopy(varAssignment)
        varbset_before = deepcopy(varbset)

        newVarbSet, newVarAssign, newFormula = solve(varbset, varAssignment, new_formula)
        if checkSat(newVarAssign, newFormula):
            return newVarbSet, newVarAssign, newFormula
        else:
            # assigning True didn't work, so assign to False
            varAssignment_before.remove(Literal(x, True))
            varAssignment_before.append(Literal(x, False))
            return solve(varbset_before, varAssignment_before, new_formula)

    return varbset, varAssignment, formula


# Print the result in DIMACS format
def printOutput(assignment):
    result = ""
    isSat = (assignment != [])
    if isSat:
        assignment = sorted(assignment, key=lambda x: int(x.name))
        for var in assignment:
            result += " " + str(var)

    print(f"s {'SATISFIABLE' if isSat else 'UNSATISFIABLE'}")
    if isSat:
        print(f"v{result} 0")

# checks whether the varAssignment is consistent with varbSet, each var should have exactly one assignment
def checkConsistency(varbSet, varAssignment):
    if len(varbSet) != len(varAssignment):
        return False
    for var in varbSet:
        if Literal(var, True) in varAssignment and Literal(var, False) in varAssignment: # contains both positive and negative
            return False
        if not(Literal(var, True) in varAssignment) and not(Literal(var, False) in varAssignment): # contains no assignment
            return False
    return True

if __name__ == "__main__":
    inputFile = sys.argv[1]
    varbset, clauseSet = readInput(inputFile)

    prev_varbset = deepcopy(varbset)

    from time import time
    start = time()
    new_varbset, varAssignment, formula = solve(varbset, [], clauseSet)
    end = time()

    # print("*********finished solve*********")
    print("c solving", inputFile)
    if len(prev_varbset) != len(varAssignment): #some variable didn't get assigned, so none of the assignments count and it's unsat
        printOutput([])
    else:
        # result is sat
        printOutput(varAssignment)
        print("c Assignments result in sat:", checkSat(varAssignment, clauseSet))
        print("c Assignment consistent with varbset:", checkConsistency(prev_varbset, varAssignment))

    print("c Finished in %.3fs" % (end-start))


    # TODO: Another SAT solver to verify unsatisfiable: https://pypi.org/project/satispy/

# https://www.cs.ubc.ca/~hoos/SATLIB/benchm.html
# consistent and no contradictions, actually a solution
# https://jgalenson.github.io/research.js/demos/minisat.html
