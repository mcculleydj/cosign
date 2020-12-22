package controller

import "github.com/gorilla/mux"

// Router constructor function
func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/bills/number", getBillsByNumber).Methods("GET").Queries("billNumbers", "{billNumbers}")
	router.HandleFunc("/api/bills/title", getBillsByTitle).Methods("GET").Queries(
		"query", "{query}",
		"bipartisan", "{bipartisan}",
		"billNumbers", "{billNumbers}",
	)
	router.HandleFunc("/api/bills/title", getBillsByTitle).Methods("GET").Queries(
		"query", "{query}",
		"bipartisan", "{bipartisan}",
	)
	router.HandleFunc("/api/bills/subject", getBillsBySubjects).Methods("GET").Queries(
		"subjects", "{subjects}",
		"bipartisan", "{bipartisan}",
	)
	router.HandleFunc("/api/members", getMembers).Methods("GET")
	router.HandleFunc("/api/cell/{position}", getCell).Methods("GET")
	router.HandleFunc("/api/cells", getCells).Methods("GET").Queries("subjects", "{subjects}")
	router.HandleFunc("/api/subjects", getSubjects).Methods("GET")
	return router
}
