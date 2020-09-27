package parse

import (
	"backend/internal/database"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Node represents a generic XML node
type Node struct {
	XMLName xml.Name
	Parent  string
	Attrs   []xml.Attr `xml:",any,attr"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`
}

func walk(parent *Node, nodes []Node, f func(Node) bool) {
	for _, n := range nodes {
		if parent != nil {
			n.Parent = parent.XMLName.Local
		}
		// ensure that f returns false for any desired stopping condition
		if f(n) {
			walk(&n, n.Nodes, f)
		}
	}
}

func parseNames(n Node) []string {
	fullNames := []string{}
	for _, child := range n.Nodes {
		for _, grandchild := range child.Nodes {
			if grandchild.XMLName.Local == "fullName" {
				fullNames = append(fullNames, strings.Replace(string(grandchild.Content), "Rep. ", "", 1))
			}
		}
	}
	return fullNames
}

func parseSubjects(n Node) []string {
	subjects := []string{}
	for _, item := range n.Nodes {
		if len(item.Nodes) != 1 {
			panic("Expected a single name node under legislative subjects item")
		}
		subjects = append(subjects, string(item.Nodes[0].Content))
	}
	return subjects
}

func aggregate(bill *database.Bill) {
	var d, r, l, i int
	for _, s := range append(bill.Sponsors, bill.Cosponsors...) {
		party := strings.Split(s, "[")[1][0]
		if party == 'D' {
			bill.NumDems++
			d = 1
		} else if party == 'R' {
			bill.NumReps++
			r = 1
		} else if party == 'I' {
			bill.NumInds++
			i = 1
		} else if party == 'L' {
			bill.NumLibs++
			l = 1
		} else {
			panic("Need to account for more party affiliations than D, R, L, and I")
		}
	}
	bill.Score = bill.NumDems - bill.NumReps
	bill.MultiParty = d+r+i+l > 1
}

func populateBill(path string, throttle chan struct{}, wg *sync.WaitGroup) {
	defer func() {
		throttle <- struct{}{}
		wg.Done()
	}()

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}
	r := bytes.NewReader(bs)
	dec := xml.NewDecoder(r)
	var n Node
	err = dec.Decode(&n)
	if err != nil {
		panic(err.Error())
	}

	bill := new(database.Bill)

	walk(nil, []Node{n}, func(n Node) bool {
		switch n.XMLName.Local {
		case "billNumber":
			n, err := strconv.Atoi(string(n.Content))
			if err != nil {
				panic(err.Error())
			}
			bill.Number = int(n)
		case "sponsors":
			if n.Parent == "bill" {
				sponsors := parseNames(n)
				if err != nil {
					panic(err.Error())
				}
				bill.Sponsors = sponsors
			}
		case "cosponsors":
			if n.Parent == "bill" {
				cosponsors := parseNames(n)
				if err != nil {
					panic(err.Error())
				}
				bill.Cosponsors = cosponsors
			}
		case "title":
			if n.Parent == "bill" {
				bill.Title = string(n.Content)
				bill.TitleLower = strings.ToLower(bill.Title)
			}
		case "policyArea":
			if len(n.Nodes) > 0 {
				if n.Nodes[0].XMLName.Local != "name" {
					panic("bad assumption on policy area")
				}
				bill.PolicyArea = string(n.Nodes[0].Content)
			}
		case "legislativeSubjects":
			bill.Subjects = parseSubjects(n)
		}
		// TODO: need a default case?
		return true
	})

	aggregate(bill)

	bill.Link = fmt.Sprintf("https://www.congress.gov/bill/116th-congress/house-bill/%d", bill.Number)

	if err = database.InsertBill(bill); err != nil {
		panic(err.Error())
	}
}

// PopulateBills parses XML into bill documents and populates the collection in Mongo
func PopulateBills() error {
	matches, err := filepath.Glob("../../bills/*.xml")
	if err != nil {
		return err
	}
	throttle := make(chan struct{}, 16)
	for i := 0; i < 16; i++ {
		throttle <- struct{}{}
	}
	wg := sync.WaitGroup{}
	for _, path := range matches {
		<-throttle
		wg.Add(1)
		go populateBill(path, throttle, &wg)
	}
	wg.Wait()
	return nil
}
