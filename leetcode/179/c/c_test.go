// Code generated by generator_test.
package main

import (
	"github.com/EndlessCheng/codeforces-go/leetcode/testutil"
	"testing"
)

func Test(t *testing.T) {
	t.Log("Current test is [c]")
	exampleIns := [][]string{{`1`, `0`, `[-1]`, `[0]`}, {`6`, `2`, `[2,2,-1,2,2,2]`, `[0,0,1,0,0,0]`}, {`7`, `6`, `[1,2,3,4,5,6,-1]`, `[0,6,5,4,3,2,1]`}, {`15`, `0`, `[-1,0,0,1,1,2,2,3,3,4,4,5,5,6,6]`, `[1,1,1,1,1,1,1,0,0,0,0,0,0,0,0]`}, {`4`, `2`, `[3,3,-1,2]`, `[0,0,162,914]`}}
	exampleOuts := [][]string{{`0`}, {`1`}, {`21`}, {`3`}, {`1076`}}
	// TODO: 测试参数的下界和上界！
	// custom test cases or WA cases.
	//exampleIns = append(exampleIns, []string{``})
	//exampleOuts = append(exampleOuts, []string{``})
	targetCaseNum := 0
	if err := testutil.RunLeetCodeFuncWithCase(t, numOfMinutes, exampleIns, exampleOuts, targetCaseNum); err != nil {
		t.Fatal(err)
	}
}
// https://leetcode-cn.com/contest/weekly-contest-179/problems/time-needed-to-inform-all-employees/
