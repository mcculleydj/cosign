package parse

import (
	"backend/internal/database"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

func appendPolicyArea(policyArea database.PolicyArea, cells []database.Cell, throttle chan struct{}, wg *sync.WaitGroup) {
	defer func() {
		throttle <- struct{}{}
		wg.Done()
	}()
	fmt.Printf("Checking cells for membership in policy area: %s...\n", policyArea.PolicyArea)
	for _, cell := range cells {
		for _, billNumber := range policyArea.BillNumbers {
			if _, ok := cell.BillNumbers[billNumber]; ok {
				filter := bson.M{"position": cell.Position}
				update := bson.M{
					"$push": bson.M{
						"policyAreas": policyArea.PolicyArea,
					},
				}
				err := database.UpsertCell(filter, update)
				if err != nil {
					panic(err.Error())
				}
				// one match is sufficient
				break
			}
		}
	}
}

func appendSubject(subject database.Subject, cells []database.Cell, throttle chan struct{}, wg *sync.WaitGroup) {
	defer func() {
		throttle <- struct{}{}
		wg.Done()
	}()
	fmt.Printf("Checking cells for membership in subject: %s...\n", subject.Subject)
	for _, cell := range cells {
		for _, billNumber := range subject.BillNumbers {
			if _, ok := cell.BillNumbers[billNumber]; ok {
				filter := bson.M{"position": cell.Position}
				update := bson.M{
					"$push": bson.M{
						"subjects": subject.Subject,
					},
				}
				err := database.UpsertCell(filter, update)
				if err != nil {
					panic(err.Error())
				}
				// one match is sufficient
				break
			}
		}
	}
}

// iterate over all subjects and policy areas
// iterate over all cells
// iterate over all bill numbers belonging to the subject or policy area
// if a bill number for this subject / policy area appears in the cell's bill number set append this subject / policy area
func updateCellSubjects() error {
	subjects, err := database.GetSubjects()
	if err != nil {
		return err
	}
	policyAreas, err := database.GetPolicyAreas()
	if err != nil {
		return err
	}
	cells, err := database.GetCells(bson.M{})
	if err != nil {
		return err
	}
	fmt.Printf("Updating cell policy areas...")
	throttle := make(chan struct{}, 16)
	for i := 0; i < 16; i++ {
		throttle <- struct{}{}
	}
	wg := sync.WaitGroup{}
	for _, policyArea := range policyAreas {
		<-throttle
		wg.Add(1)
		go appendPolicyArea(policyArea, cells, throttle, &wg)
	}
	wg.Wait()
	fmt.Printf("Updating cell subjects...")
	for _, subject := range subjects {
		<-throttle
		wg.Add(1)
		go appendSubject(subject, cells, throttle, &wg)
	}
	wg.Wait()
	return nil
}

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

	return updateCellSubjects()
}
