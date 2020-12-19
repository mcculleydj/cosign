package main

import (
	"backend/internal/database"
	"backend/pkg/parse"
	"flag"
	"fmt"
)

func main() {
	populateBills := flag.Bool("b", false, "Populate bills")
	populateMembers := flag.Bool("m", false, "Populate members")
	populateCells := flag.Bool("c", false, "Populate cells and member counts")
	populateSubjects := flag.Bool("s", false, "Populate policy areas and subjects")
	flag.Parse()

	if err := database.Connect(); err != nil {
		panic("Mongo connect error: " + err.Error())
	}
	defer database.Disconnect()

	if !*populateBills && !*populateMembers && !*populateCells && !*populateSubjects {
		*populateBills = true
		*populateMembers = true
		*populateCells = true
		*populateSubjects = true
	}

	if *populateBills {
		fmt.Println("Populating bills collection...")
		err := parse.PopulateBills()
		if err != nil {
			panic("Populate bills error: " + err.Error())
		}
	}

	if *populateMembers {
		fmt.Println("Populating members collection...")
		err := parse.PopulateMembers()
		if err != nil {
			panic("Populate members error: " + err.Error())
		}
	}

	if *populateCells {
		fmt.Println("Populating cells collection...")
		err := parse.PopulateCells()
		if err != nil {
			panic("Populate cells error: " + err.Error())
		}
		fmt.Println("Populating member counts...")
		err = parse.PopulateCounts()
		if err != nil {
			panic("Populate counts error: " + err.Error())
		}
	}

	if *populateSubjects {
		fmt.Println("Populating policy areas and subjects collection...")
		err := parse.PopulateSubjects()
		if err != nil {
			panic("Populate subjects error: " + err.Error())
		}
	}
}
