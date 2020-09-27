package database

// Bill describes a piece of legislation
type Bill struct {
	Number     int      `json:"number" bson:"number"`
	Title      string   `json:"title" bson:"title"`
	TitleLower string   `json:"-" bson:"titleLower"`
	Sponsors   []string `json:"sponsors" bson:"sponsors"`
	Cosponsors []string `json:"cosponsors" bson:"cosponsors"`
	Score      int      `json:"score" bson:"score"`
	NumDems    int      `json:"numDems" bson:"numDems"`
	NumReps    int      `json:"numReps" bson:"numReps"`
	NumInds    int      `json:"numInds" bson:"numInds"`
	NumLibs    int      `json:"numLibs" bson:"numLibs"`
	MultiParty bool     `json:"multiParty" bson:"multiParty"`
	Link       string   `json:"link" bson:"link"`
	PolicyArea string   `json:"policyArea" bson:"policyArea"`
	Subjects   []string `json:"subjects" bson:"subjects"`
}

// Member describes a member of the House
type Member struct {
	ID          int            `json:"id" bson:"id"`
	Name        string         `json:"name" bson:"name"`
	Parties     []string       `json:"parties" bson:"parties"`
	Districts   []string       `json:"districts" bson:"districts"`
	State       string         `json:"state" bson:"state"`
	FullStrings []string       `json:"-" bson:"fullStrings"`
	Counts      map[string]int `json:"counts" bson:"counts"`
}

// Cell describes the adjacency matrix cell data
type Cell struct {
	Position    string `json:"position" bson:"position"`
	Count       int    `json:"count" bson:"count"`
	BillNumbers []int  `json:"-" bson:"billNumbers"`
	Bills       []Bill `json:"bills" bson:"-"`
}

// PolicyArea describes a policy area category on a bill
type PolicyArea struct {
	PolicyArea  string `json:"policyArea" bson:"policyArea"`
	BillNumbers []int  `json:"billNumbers" bson:"billNumbers"`
}

// Subject describes a subject category on a bill
type Subject struct {
	Subject     string `json:"subject" bson:"subject"`
	BillNumbers []int  `json:"billNumbers" bson:"billNumbers"`
}
