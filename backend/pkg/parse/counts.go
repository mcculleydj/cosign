package parse

import (
	"backend/internal/database"
	"fmt"
	"strings"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setCounts(m database.Member, throttle chan struct{}, wg *sync.WaitGroup) {
	defer func() {
		throttle <- struct{}{}
		wg.Done()
	}()

	counts := map[string]int{}
	patterns := []string{
		fmt.Sprintf("^%d_", m.ID),
		fmt.Sprintf("_%d$", m.ID),
	}

	for i, pattern := range patterns {
		filter := bson.M{
			"position": bson.M{
				"$regex": primitive.Regex{Pattern: pattern, Options: "i"},
			},
		}
		cells, err := database.GetCells(filter)
		if err != nil {
			panic(err.Error())
		}

		for _, c := range cells {
			var id string
			if i == 0 {
				id = strings.Replace(c.Position, fmt.Sprintf("%d_", m.ID), "", 1)
			} else {
				id = strings.Replace(c.Position, fmt.Sprintf("_%d", m.ID), "", 1)
			}
			counts[id] = c.Count
		}
	}

	filter := bson.M{"id": m.ID}
	update := bson.M{
		"$set": bson.M{
			"counts": counts,
		},
	}
	if err := database.UpdateMember(filter, update); err != nil {
		panic(err.Error())
	}
}

// PopulateCounts maps member IDs to number of bills cosponsored
func PopulateCounts() error {
	members, err := database.GetMembers(bson.M{})
	if err != nil {
		return err
	}
	throttle := make(chan struct{}, 16)
	for i := 0; i < 16; i++ {
		throttle <- struct{}{}
	}
	wg := sync.WaitGroup{}
	for _, m := range members {
		<-throttle
		wg.Add(1)
		go setCounts(m, throttle, &wg)
	}
	wg.Wait()
	return nil
}
