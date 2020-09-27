package controller

import "github.com/gorilla/mux"

// Router constructor function
func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/bills/number", getBillsByNumber).Methods("GET").Queries("numbers", "{numbers}")
	router.HandleFunc("/api/bills/title", getBillsByTitle).Methods("GET").Queries(
		"query", "{query}",
		"bipartisan", "{bipartisan}",
		"billNumbers", "{billNumbers}",
	)
	router.HandleFunc("/api/bills/title", getBillsByTitle).Methods("GET").Queries(
		"query", "{query}",
		"bipartisan", "{bipartisan}",
	)
	router.HandleFunc("/api/members", getMembers).Methods("GET")
	router.HandleFunc("/api/cell/{position}", getCell).Methods("GET")
	router.HandleFunc("/api/cells/{id}", getCells).Methods("GET")
	router.HandleFunc("/api/subjects", getSubjects).Methods("GET")
	return router
}