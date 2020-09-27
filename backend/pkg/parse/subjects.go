package parse

import (
	"backend/internal/database"

	"go.mongodb.org/mongo-driver/bson"
)

// PopulateSubjects populates the policy areas and subjects collection from information in bills collection
func PopulateSubjects() error {
	bills, err := database.GetBills(bson.M{})
	if err != nil {
		return err
	}

	policyAreaMap := map[string][]int{}
	subjectMap := map[string][]int{}

	for _, b := range bills {
		if b.PolicyArea == "" {
			continue
		}
		if billNumbers, ok := policyAreaMap[b.PolicyArea]; !ok {
			policyAreaMap[b.PolicyArea] = []int{b.Number}
		} else {
			policyAreaMap[b.PolicyArea] = append(billNumbers, b.Number)
		}
		for _, subject := range b.Subjects {
			if subject == "" {
				continue
			}
			if billNumbers, ok := subjectMap[subject]; !ok {
				subjectMap[subject] = []int{b.Number}
			} else {
				subjectMap[subject] = append(billNumbers, b.Number)
			}
		}
	}

	for policyArea, billNumbers := range policyAreaMap {
		err := database.InsertPolicyArea(policyArea, billNumbers)
		if err != nil {
			return err
		}
	}

	for subject, billNumbers := range subjectMap {
		err := database.InsertSubject(subject, billNumbers)
		if err != nil {
			return err
		}
	}

	return nil
}
