package controller

import (
	"backend/internal/database"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getBillsByNumber(w http.ResponseWriter, r *http.Request) {
	numbersString := r.FormValue("billNumbers")
	numbers := []int{}
	for _, s := range strings.Split(numbersString, ",") {
		n, err := strconv.Atoi(s)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "Incorrect numbers param", "")
			return
		}
		numbers = append(numbers, n)
	}
	filter := bson.M{
		"number": bson.M{
			"$in": numbers,
		},
	}
	bills, err := database.GetBills(filter)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Unable to get bills", "")
		return
	}
	WriteResponse(w, bills)
}

func getBillsByTitle(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	bipartisan := r.FormValue("bipartisan")
	strBillNumbers := r.FormValue("billNumbers")

	filter := bson.M{}

	if query != "*" {
		filter["titleLower"] = bson.M{"$regex": strings.ToLower(query)}
	}

	if bipartisan == "true" {
		filter["multiParty"] = true
	}

	if len(strBillNumbers) > 0 {
		// convert from string to int
		billNumbers := []int{}
		for _, s := range strings.Split(strBillNumbers, ",") {
			n, err := strconv.Atoi(s)
			if err != nil {
				WriteError(w, http.StatusBadRequest, "Bill numbers much be integers", "")
				return
			}
			billNumbers = append(billNumbers, n)
		}
		filter["number"] = bson.M{
			"$in": billNumbers,
		}
	}

	bills, err := database.GetBills(filter)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Unable to get bills", "")
		return
	}
	WriteResponse(w, bills)
}

func getBillsBySubjects(w http.ResponseWriter, r *http.Request) {
	subjectsQuery := r.FormValue("subjects")
	if subjectsQuery == "" {
		WriteError(w, http.StatusBadRequest, "Missing subjects", "")
		return
	}
	subjects := strings.Split(subjectsQuery, ",")
	filter := bson.M{
		"subject": bson.M{
			"$in": subjects,
		},
	}
	subjectDocuments, err := database.GetSubjects(filter)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error retreiving subjects", "")
		return
	}
	billNumbers := []int{}
	for _, subjectDocument := range subjectDocuments {
		billNumbers = append(billNumbers, subjectDocument.BillNumbers...)
	}
	filter = bson.M{
		"number": bson.M{
			"$in": billNumbers,
		},
	}
	bipartisan := r.FormValue("bipartisan")
	if bipartisan == "true" {
		filter["multiParty"] = true
	}
	bills, err := database.GetBills(filter)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error retreiving bills", "")
		return
	}
	WriteResponse(w, bills)
}

func getMembers(w http.ResponseWriter, r *http.Request) {
	members, memberMap, err := database.GetMembers(bson.M{})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error retrieving members", "")
		return
	}
	payload := map[string]interface{}{
		"members":   members,
		"memberMap": memberMap,
	}
	WriteResponse(w, payload)
}

func getCell(w http.ResponseWriter, r *http.Request) {
	position, found := mux.Vars(r)["position"]
	if !found {
		WriteError(w, http.StatusBadRequest, "Missing position paramater", "")
		return
	}
	cell, err := database.GetCell(bson.M{"position": position})
	if err == mongo.ErrNoDocuments {
		WriteResponse(w, database.Cell{Bills: []database.Bill{}})
		return
	} else if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error retrieving cell data", err.Error())
		return
	}
	WriteResponse(w, cell)
}

func getCells(w http.ResponseWriter, r *http.Request) {
	subjectsStr, hasSubjects := mux.Vars(r)["subjects"]

	if !hasSubjects || subjectsStr == "" {
		WriteError(w, http.StatusBadRequest, "Missing subject", "")
		return
	}

	subjects := strings.Split(subjectsStr, ",")

	filter := bson.M{
		"subjects": bson.M{
			"$in": subjects,
		},
	}

	cells, err := database.GetCells(filter, subjects)

	if err == mongo.ErrNoDocuments {
		WriteError(w, http.StatusNotFound, "Unable to find any documents", "")
		return
	} else if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error retrieving cell data", err.Error())
		return
	}

	WriteResponse(w, cells)
}

func getSubjects(w http.ResponseWriter, r *http.Request) {
	policyAreas, err := database.GetPolicyAreas()
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error retrieving policy area data", err.Error())
		return
	}
	subjects, err := database.GetSubjects(bson.M{})
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Error retrieving subject data", err.Error())
		return
	}
	WriteResponse(w, map[string]interface{}{
		"policyAreas": policyAreas,
		"subjects":    subjects,
	})
}
