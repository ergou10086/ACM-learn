// Code generated by copypasta/template/leetcode/generator_test.go
package main

import (
	"github.com/EndlessCheng/codeforces-go/leetcode/testutil"
	"testing"
)

func Test(t *testing.T) {
	t.Log("Current test is [b]")
	examples := [][]string{
		{
			`"aba"`, 
			`6`,
		},
		{
			`"abc"`, 
			`3`,
		},
		{
			`"ltcd"`, 
			`0`,
		},
		{
			`"noosabasboosa"`, 
			`237`,
		},
		
	}
	targetCaseNum := 0 // -1
	if err := testutil.RunLeetCodeFuncWithExamples(t, countVowels, examples, targetCaseNum); err != nil {
		t.Fatal(err)
	}
}
// https://leetcode-cn.com/contest/weekly-contest-266/problems/vowels-of-all-substrings/
