package main

import (
	"backend/internal/database"
	"flag"
	"fmt"
)

func main() {
	dev := flag.Bool("d", false, "Use local mongo")
	user := flag.String("u", "", "Mongo username")
	password := flag.String("p", "", "Mongo password")
	dropBills := flag.Bool("b", false, "Drop bills")
	dropMembers := flag.Bool("m", false, "Drop members")
	dropCells := flag.Bool("c", false, "Drop cells")
	dropSubjects := flag.Bool("s", false, "Drop subjects")
	flag.Parse()

	var uri string
	if *dev {
		uri = "mongodb://localhost:27017"
	} else {
		uri = fmt.Sprintf("mongodb+srv://%s:%s@cluster0.wht7g.mongodb.net/admin?retryWrites=true&w=majority", *user, *password)
	}

	if err := database.Connect(uri); err != nil {
		panic(err.Error())
	}
	defer database.Disconnect()

	if err := database.Clean(*dropBills, *dropMembers, *dropCells, *dropSubjects); err != nil {
		panic(err.Error())
	}

}
