package main

import (
	"backend/internal/database"
	"flag"
)

func main() {
	dropBills := flag.Bool("b", false, "Drop bills")
	dropMembers := flag.Bool("m", false, "Drop members")
	dropCells := flag.Bool("c", false, "Drop cells")
	dropSubjects := flag.Bool("s", false, "Drop subjects")
	flag.Parse()

	if err := database.Connect(); err != nil {
		panic(err.Error())
	}
	defer database.Disconnect()

	if err := database.Clean(*dropBills, *dropMembers, *dropCells, *dropSubjects); err != nil {
		panic(err.Error())
	}

}
