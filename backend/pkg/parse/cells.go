package parse

import (
	"backend/internal/database"
	"fmt"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

// PartyID stores a member's party and ID in the context of their appearance
// on a particular bill (e.g. Justin Amash's party will appear as R, I, and L)
type PartyID struct {
	Party byte
	ID    int
}

func buildNameToIDMap() (map[string]int, error) {
	nameToID := map[string]int{}
	members, _, err := database.GetMembers(bson.M{})
	if err != nil {
		return nameToID, err
	}
	for _, m := range members {
		for _, name := range m.FullStrings {
			nameToID[name] = m.ID
		}
	}
	return nameToID, nil
}

func updateCells(members []PartyID, billNumber int, throttle chan struct{}, wg *sync.WaitGroup) error {
	defer func() {
		throttle <- struct{}{}
		wg.Done()
	}()
	for _, i := range members {
		for _, j := range members {
			if i.ID >= j.ID || i.Party == j.Party {
				continue
			}
			position := fmt.Sprintf("%d_%d", i.ID, j.ID)
			filter := bson.M{"position": position}
			update := bson.M{
				"$setOnInsert": bson.M{
					"position": position,
				},
				"$inc":  bson.M{"count": 1},
				"$push": bson.M{"billNumbers": billNumber},
			}
			if err := database.UpsertCell(filter, update); err != nil {
				return err
			}
		}
	}
	return nil
}

// PopulateCells populates the cells of the adjacency matrix
func PopulateCells() error {
	nameToID, err := buildNameToIDMap()
	if err != nil {
		return err
	}
	bills, err := database.GetBills(bson.M{"multiParty": true})
	if err != nil {
		return err
	}
	throttle := make(chan struct{}, 16)
	for i := 0; i < 16; i++ {
		throttle <- struct{}{}
	}
	wg := sync.WaitGroup{}
	for _, b := range bills {
		members := []PartyID{}
		for _, s := range append(b.Sponsors, b.Cosponsors...) {
			party := strings.Split(s, "[")[1][0]
			members = append(members, PartyID{party, nameToID[s]})
		}
		<-throttle
		wg.Add(1)
		go updateCells(members, b.Number, throttle, &wg)
	}
	wg.Wait()
	return nil
}
