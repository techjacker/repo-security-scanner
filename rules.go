package main

import (
	"fmt"

	df "github.com/techjacker/diffence"
)

// RuleResults is the struct for revealing the result of diff rule analysis
type RuleResults struct {
	Pass bool
	Info string
}

func checkRuleResults(ruleResults df.Results, commitID string) RuleResults {

	dirty := false
	info := fmt.Sprintf("Commit ID: %s\n", commitID)

	for k, v := range ruleResults {
		if len(v) > 0 {
			dirty = true
			info += fmt.Sprintf("File %s violates %d rules:\n", k, len(v))
			for _, r := range v {
				info += fmt.Sprintf("\n%s\n", r.String())
			}
		}
	}

	return RuleResults{
		Pass: !dirty,
		Info: info,
	}
}

func runRules(rules *[]df.Rule, URL string) (df.Results, error) {
	resp, err := getGithubDiff(URL)
	defer resp.Body.Close()
	if err != nil {
		return df.Results{}, err
	}
	return df.CheckDiffs(resp.Body, rules)
}
