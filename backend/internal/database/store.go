package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ctx() context.Context {
	return context.Background()
}

func indexOpts() *options.IndexOptions {
	return options.Index().SetUnique(true)
}

// Clean drops collections and recreates indices
func Clean(dropBills, dropMembers, dropCells, dropSubjects bool) error {
	if dropBills {
		if err := billsCollection.Drop(ctx()); err != nil {
			return err
		}
		indices := []mongo.IndexModel{
			{Keys: bson.M{"number": 1}, Options: indexOpts()},
			{Keys: bson.M{"titleLower": 1}},
			{Keys: bson.M{"hasBothParties": 1}},
		}
		if _, err := billsCollection.Indexes().CreateMany(ctx(), indices); err != nil {
			return err
		}
	}

	if dropMembers {
		if err := membersCollection.Drop(ctx()); err != nil {
			return err
		}
		indices := []mongo.IndexModel{
			{Keys: bson.M{"id": 1}, Options: indexOpts()},
			{Keys: bson.M{"name": 1}, Options: indexOpts()},
		}
		if _, err := membersCollection.Indexes().CreateMany(ctx(), indices); err != nil {
			return err
		}
	}

	if dropCells {
		if err := cellsCollection.Drop(ctx()); err != nil {
			return err
		}
		indices := []mongo.IndexModel{
			{Keys: bson.M{"position": 1}, Options: indexOpts()},
		}
		if _, err := cellsCollection.Indexes().CreateMany(ctx(), indices); err != nil {
			return err
		}
	}

	if dropSubjects {
		if err := policyAreasCollection.Drop(ctx()); err != nil {
			return err
		}
		indices := []mongo.IndexModel{
			{Keys: bson.M{"policyArea": 1}, Options: indexOpts()},
		}
		if _, err := policyAreasCollection.Indexes().CreateMany(ctx(), indices); err != nil {
			return err
		}
		if err := subjectsCollection.Drop(ctx()); err != nil {
			return err
		}
		indices = []mongo.IndexModel{
			{Keys: bson.M{"subject": 1}, Options: indexOpts()},
		}
		if _, err := subjectsCollection.Indexes().CreateMany(ctx(), indices); err != nil {
			return err
		}
	}

	return nil
}

// InsertBill inserts a bill into the database
func InsertBill(b *Bill) error {
	_, err := billsCollection.InsertOne(ctx(), b)
	return err
}

// GetBills returns bills matching the supplied filter
func GetBills(filter bson.M) ([]Bill, error) {
	var bills []Bill
	cur, err := billsCollection.Find(ctx(), filter)
	if err != nil {
		return bills, err
	}
	defer cur.Close(ctx())
	err = cur.All(ctx(), &bills)
	return bills, err
}

// GetSponsors passes over bills collection and extracts sponsor data
func GetSponsors() (map[string]bool, error) {
	names := map[string]bool{}
	opts := options.Find()
	opts.SetProjection(bson.M{"sponsors": 1, "cosponsors": 1})
	cur, err := billsCollection.Find(ctx(), bson.M{}, opts)
	if err != nil {
		return names, err
	}
	defer cur.Close(ctx())
	for cur.Next(context.Background()) {
		var bill Bill
		err = cur.Decode(&bill)
		if err != nil {
			return names, err
		}
		for _, name := range append(bill.Sponsors, bill.Cosponsors...) {
			names[name] = true
		}
	}
	return names, nil
}

// InsertMembers inserts a slice of members into the database
func InsertMembers(ms []interface{}) error {
	_, err := membersCollection.InsertMany(ctx(), ms)
	return err
}

// GetMembers returns all members from the database
func GetMembers(filter bson.M) ([]Member, error) {
	var members []Member
	cur, err := membersCollection.Find(ctx(), filter)
	if err != nil {
		return members, err
	}
	defer cur.Close(ctx())
	err = cur.All(ctx(), &members)
	return members, err
}

// UpdateMember updates a member document
func UpdateMember(filter, update bson.M) error {
	_, err := membersCollection.UpdateOne(ctx(), filter, update)
	return err
}

// UpsertCell upserts an adjacency cell
func UpsertCell(filter, update bson.M) error {
	opts := options.Update()
	opts.SetUpsert(true)
	_, err := cellsCollection.UpdateOne(ctx(), filter, update, opts)
	return err
}

// GetCell returns a cell matching the filter
func GetCell(filter bson.M) (Cell, error) {
	var cell Cell
	err := cellsCollection.FindOne(ctx(), filter).Decode(&cell)
	if err != nil {
		return cell, err
	}
	billsFilter := bson.M{
		"number": bson.M{
			"$in": cell.BillNumbers,
		},
	}
	bills, err := GetBills(billsFilter)
	if err != nil {
		return cell, err
	}
	cell.Bills = bills
	return cell, err
}

// GetCells returns cells matching the supplied filter
func GetCells(filter bson.M) ([]Cell, error) {
	var cells []Cell
	cur, err := cellsCollection.Find(ctx(), filter)
	if err != nil {
		return cells, err
	}
	defer cur.Close(ctx())
	err = cur.All(ctx(), &cells)
	return cells, err
}

// InsertSubject inserts a subject into the database
func InsertSubject(subject string, billNumbers []int) error {
	doc := bson.M{
		"subject":     subject,
		"billNumbers": billNumbers,
	}
	_, err := subjectsCollection.InsertOne(ctx(), doc)
	return err
}

// InsertPolicyArea inserts a policy area into the database
func InsertPolicyArea(policyArea string, billNumbers []int) error {
	doc := bson.M{
		"policyArea":  policyArea,
		"billNumbers": billNumbers,
	}
	_, err := policyAreasCollection.InsertOne(ctx(), doc)
	return err
}

// GetPolicyAreas returns all policy areas
func GetPolicyAreas() ([]PolicyArea, error) {
	var policyAreas []PolicyArea
	cur, err := policyAreasCollection.Find(ctx(), bson.M{})
	if err != nil {
		return policyAreas, err
	}
	defer cur.Close(ctx())
	err = cur.All(ctx(), &policyAreas)
	return policyAreas, err
}

// GetSubjects returns all subjects
func GetSubjects() ([]Subject, error) {
	var subjects []Subject
	cur, err := subjectsCollection.Find(ctx(), bson.M{})
	if err != nil {
		return subjects, err
	}
	defer cur.Close(ctx())
	err = cur.All(ctx(), &subjects)
	return subjects, err
}
