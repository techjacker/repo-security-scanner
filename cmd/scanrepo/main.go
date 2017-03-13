package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/techjacker/diffence"
)

func main() {

	rPath := flag.String("rules", "", "path to custom rules in JSON format")
	flag.Parse()

	info, _ := os.Stdin.Stat()
	if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		log.Fatalln("The command is intended to work with pipes.")
		return
	}

	var (
		err   error
		rules *[]diffence.Rule
	)

	if len(*rPath) > 0 {
		rules, err = diffence.LoadRulesJSON(*rPath)
	} else {
		rules, err = diffence.LoadDefaultRules()
	}
	if err != nil {
		log.Fatalf("Cannot load rules\n%s", err)
		return
	}

	diff := diffence.DiffChecker{
		Rules:   rules,
		Ignorer: diffence.NewIgnorerFromFile(".secignore"),
	}
	res, err := diff.Check(bufio.NewReader(os.Stdin))
	if err != nil {
		log.Fatalf("Error reading diff\n%s\n", err)
		return
	}

	matches := res.Matches()
	if matches > 0 {
		i := 1
		fmt.Printf("Diff contains %d offenses\n\n", matches)
		for filename, rule := range res.MatchedRules {
			fmt.Printf("------------------\n")
			fmt.Printf("Violation %d\n", i)
			fmt.Printf("File: %s\n", filename)
			fmt.Printf("Reason: %#v\n\n", rule[0].Caption)
			i++
		}
		// finding violations constitutes an error
		os.Exit(1)
		return
	}
	fmt.Printf("Diff contains NO offenses\n\n")
	os.Exit(0)
}
