package parse

import (
	"backend/internal/database"
	"backend/pkg/utility"
	"strings"
)

func parsePSD(s string) (string, string, string) {
	tokens := strings.Split(strings.TrimRight(s, "]"), "-")
	return tokens[0], tokens[1], tokens[2]
}

// PopulateMembers populates the members collection from information in bills collection
func PopulateMembers() error {
	names, err := database.GetSponsors()
	if err != nil {
		return err
	}

	m := map[string]*database.Member{}
	id := 1

	for s := range names {
		tokens := strings.Split(s, " [")
		name := tokens[0]
		if member, ok := m[name]; !ok {
			party, state, district := parsePSD(tokens[1])
			m[name] = &database.Member{
				ID:          id,
				Name:        name,
				Parties:     []string{party},
				Districts:   []string{district},
				State:       state,
				FullStrings: []string{s},
			}
			id++
		} else {
			if !utility.Contains(member.FullStrings, s) {
				party, _, district := parsePSD(tokens[1])
				if !utility.Contains(member.Parties, party) {
					member.Parties = append(member.Parties, party)
				}
				if !utility.Contains(member.Districts, district) {
					member.Districts = append(member.Districts, district)
				}
				member.FullStrings = append(member.FullStrings, s)
			}
		}
	}

	members := []interface{}{}
	for _, member := range m {
		members = append(members, member)
	}
	return database.InsertMembers(members)
}
