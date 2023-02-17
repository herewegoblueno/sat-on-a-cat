import unittest
from solver import *

class SolverTests(unittest.TestCase):

    ######## unitClauseElim ##########

    # base case
    def test_unit_clause_base_case(self):
        # input
        varbset = []
        varAssignment = []
        formula = []
        # output
        expected_varbset = []
        expected_varAssignment = []
        expected_formula = []
        actual_varbset, actual_varAssignment, actual_formula = unitClauseElim(varbset, varAssignment, formula)
        self.assertEqual(actual_varbset, expected_varbset)
        self.assertEqual(actual_varAssignment, expected_varAssignment)
        self.assertEqual(actual_formula, expected_formula)

    # base case
    def test_unit_clause_one_case(self):
        # input
        varbset = ['1']
        varAssignment = []
        formula = [
            Clause(0, [Literal("1", False)])
        ]
        # output
        expected_varbset = []
        expected_varAssignment = [Literal("1", False)]
        expected_formula = []
        actual_varbset, actual_varAssignment, actual_formula = unitClauseElim(varbset, varAssignment, formula)
        self.assertEqual(actual_varbset, expected_varbset)
        self.assertEqual(actual_varAssignment, expected_varAssignment)
        self.assertEqual(actual_formula, expected_formula)

    # remove all
    def test_unit_clause_remove_all(self):
        # input
        varbset = ['1', '2', '3']
        varAssignment = []
        formula = [
            Clause(0, [Literal("1", False)]), 
            Clause(1, [Literal("1", True), Literal("2", True)]), 
            Clause(2, [Literal("3", False)])
        ]
        # output
        expected_varbset = []
        expected_varAssignment = [Literal("1", False), Literal("2", True), Literal("3", False)]
        expected_formula = []
        actual_varbset, actual_varAssignment, actual_formula = unitClauseElim(varbset, varAssignment, formula)
        self.assertEqual(actual_varbset, expected_varbset)
        self.assertEqual(actual_varAssignment, expected_varAssignment)
        self.assertEqual(actual_formula, expected_formula)

    # remove one, not all literals
    def test_unit_clause_leftover(self):
        # input
        varbset = ['1', '2', '3']
        varAssignment = []
        formula = [
            Clause(0, [Literal("1", False)]), 
            Clause(1, [Literal("1", True), Literal("2", True), Literal("3", True)]), 
            Clause(2, [Literal("1", False), Literal("3", False)]), 
            Clause(3, [Literal("2", True), Literal("3", False)])
        ]
        # output
        expected_varbset = ['2', '3']
        expected_varAssignment = [Literal("1", False)]
        expected_formula = [
            Clause(1, [Literal("2", True), Literal("3", True)]), 
            Clause(3, [Literal("2", True), Literal("3", False)])
        ]
        actual_varbset, actual_varAssignment, actual_formula = unitClauseElim(varbset, varAssignment, formula)
        self.assertEqual(actual_varbset, expected_varbset)
        self.assertEqual(actual_varAssignment, expected_varAssignment)
        self.assertEqual(actual_formula, expected_formula)

    # empty clause
    def test_unit_clause_empty_clause(self):
        # input
        varbset = ['1', '2', '3', '4']
        varAssignment = []
        formula = [
            Clause(0, [Literal("1", False), Literal("2", True), Literal("3", True)]), 
            Clause(1, [Literal("1", True), Literal("4", True)]), 
            Clause(2, [Literal("1", False), Literal("2", False)]), 
            Clause(3, [Literal("1", False), Literal("2", True), Literal("4", True)]),
            Clause(4, [Literal("4", False)])
        ]
        # output
        expected_formula = [
            Clause(3, [])
        ]
        actual_varbset, actual_varAssignment, actual_formula = unitClauseElim(varbset, varAssignment, formula)
        self.assertEqual(actual_formula, expected_formula)

    ######## pureLiteralElim ##########

    # base
    def test_pureLiteralElim_baseCase(self):
        # input
        varbset = []
        varAssignment = []
        formula = []
        # output
        expected_varbset = []
        expected_varAssignment = []
        expected_formula = []
        actual_varbset, actual_varAssignment, actual_formula = pureLiteralElim(varbset, varAssignment, formula)
        self.assertEqual(actual_varbset, expected_varbset)
        self.assertEqual(actual_varAssignment, expected_varAssignment)
        self.assertEqual(actual_formula, expected_formula)

    # base
    def test_pureLiteralElim_oneCase(self):
        # input
        varbset = ['1']
        varAssignment = []
        formula = [
            Clause(0, [Literal("1", True)]), 
            Clause(1, [Literal("1", True)])
        ]
        # output
        expected_varbset = []
        expected_varAssignment = [Literal("1", True)]
        expected_formula = []
        actual_varbset, actual_varAssignment, actual_formula = pureLiteralElim(varbset, varAssignment, formula)
        self.assertEqual(actual_varbset, expected_varbset)
        self.assertEqual(actual_varAssignment, expected_varAssignment)
        self.assertEqual(actual_formula, expected_formula)

    # elimAll
    def test_pureLiteralElim_elimAll(self):
        # input
        varbset = ['1', '2', '3', '4']
        varAssignment = []
        formula = [
            Clause(0, [Literal("1", False), Literal("2", True), Literal("3", True)]), 
            Clause(1, [Literal("1", True), Literal("3", True)]), 
            Clause(2, [Literal("1", False), Literal("2", False)]), 
            Clause(3, [Literal("1", False), Literal("2", True), Literal("4", True)]),
            Clause(4, [Literal("2", False), Literal("4", False)])
        ]
        # output
        expected_varbset = []
        expected_varAssignment = [Literal("3", True), Literal("1", False), Literal("2", False), Literal("4", False)]
        expected_formula = []
        actual_varbset, actual_varAssignment, actual_formula = pureLiteralElim(varbset, varAssignment, formula)
        self.assertEqual(actual_varbset, expected_varbset)
        self.assertEqual(actual_varAssignment, expected_varAssignment)
        self.assertEqual(actual_formula, expected_formula)

    ######## checkSat ##########

    def test_checkSat_baseCase(self):
        # input
        varAssignment = []
        formula = []
        # output
        self.assertTrue(checkSat(varAssignment, formula))

    def test_checkSat_someCase(self):
        # input
        varAssignment = [Literal("1", True), Literal("2", False)]
        formula = [
            Clause(0, [Literal("1", True), Literal("2", True)]), 
            Clause(1, [Literal("2", False)])
        ]
        # output
        self.assertTrue(checkSat(varAssignment, formula))

    def test_checkSat_larger_sat(self):
        # input
        varAssignment = [Literal("1", False), Literal("2", False), Literal("3", True)]
        varAssignment2 = [Literal("1", True), Literal("2", False), Literal("3", True)]
        varAssignment3 = [Literal("1", True), Literal("2", True), Literal("3", True)]
        formula = [ 
            Clause(0, [Literal("1", False), Literal("2", True), Literal("3", True)]), 
            Clause(1, [Literal("1", True), Literal("2", False), Literal("3", False)]), 
            Clause(2, [Literal("1", False), Literal("2", False), Literal("3", True)])
        ]
        # output
        # multiple satisfies
        self.assertTrue(checkSat(varAssignment, formula))
        self.assertTrue(checkSat(varAssignment2, formula))
        self.assertTrue(checkSat(varAssignment3, formula))

    def test_checkSat_unsat(self):
        # input
        varAssignment = [Literal("1", True), Literal("2", True)]
        formula = [
            Clause(0, [Literal("1", True), Literal("2", True)]), 
            Clause(1, [Literal("2", False)])
        ]
        # output
        self.assertFalse(checkSat(varAssignment, formula))

    def test_checkSat_larger_unsat(self):
        # input
        varAssignment = [Literal("1", False), Literal("2", False), Literal("3", True)]
        varAssignment2 = [Literal("1", False), Literal("2", True), Literal("3", True)]
        varAssignment3 = [Literal("1", True), Literal("2", False), Literal("3", False)]
        formula = [
            Clause(0, [Literal("1", True), Literal("2", True), Literal("3", True)]), 
            Clause(1, [Literal("1", True), Literal("2", True), Literal("3", False)]), 
            Clause(2, [Literal("1", True), Literal("2", False), Literal("3", True)]), 
            Clause(3, [Literal("1", False), Literal("2", True), Literal("3", True)]), 
            Clause(4, [Literal("1", True), Literal("2", False), Literal("3", False)]), 
            Clause(5, [Literal("1", False), Literal("2", False), Literal("3", True)]), 
            Clause(6, [Literal("1", False), Literal("2", True), Literal("3", False)]), 
            Clause(7, [Literal("1", False), Literal("2", False), Literal("3", False)])
        ]
        # output
        # multiple unsatisfy
        self.assertFalse(checkSat(varAssignment, formula))
        self.assertFalse(checkSat(varAssignment2, formula))
        self.assertFalse(checkSat(varAssignment3, formula))

    ######## checkConsistency ##########

    def test_checkConsistency_baseCase(self):
        # input
        varAssignment = []
        varbSet = []
        # output
        self.assertTrue(checkConsistency(varbSet, varAssignment))

    def test_checkConsistency_someCase(self):
        # input
        varAssignment = [Literal("1", True), Literal("2", False)]
        varbSet = ['1','2']
        # output
        self.assertTrue(checkConsistency(varbSet, varAssignment))

    def test_checkConsistency_inconsistent(self):
        # input
        varAssignment = [Literal("1", True), Literal("1", False)]
        varAssignment2 = [Literal("1", True), Literal("1", True), Literal("2", False)]
        varAssignment3 = [Literal("1", True)]
        varAssignment4 = []
        varbSet = ['1','2']
        # output
        self.assertFalse(checkConsistency(varbSet, varAssignment))
        self.assertFalse(checkConsistency(varbSet, varAssignment2))
        self.assertFalse(checkConsistency(varbSet, varAssignment3))
        self.assertFalse(checkConsistency(varbSet, varAssignment4))

if __name__ == '__main__':
    unittest.main()