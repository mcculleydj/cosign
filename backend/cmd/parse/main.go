package main

import (
	"backend/internal/database"
	"backend/pkg/parse"
	"flag"
	"fmt"
)

func main() {
	dev := flag.Bool("d", false, "Use local mongo")
	user := flag.String("u", "", "Mongo username")
	password := flag.String("p", "", "Mongo password")
	populateBills := flag.Bool("b", false, "Populate bills")
	populateMembers := flag.Bool("m", false, "Populate members")
	populateCells := flag.Bool("c", false, "Populate cells and member counts")
	populateSubjects := flag.Bool("s", false, "Populate policy areas and subjects")
	flag.Parse()

	var uri string
	if *dev {
		uri = "mongodb://localhost:27017"
	} else {
		uri = fmt.Sprintf("mongodb+srv://%s:%s@cluster0.wht7g.mongodb.net/admin?retryWrites=true&w=majority", *user, *password)
	}

	if err := database.Connect(uri); err != nil {
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
