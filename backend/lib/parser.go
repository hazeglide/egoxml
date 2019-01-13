package egoxml

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Entry struct {
	Time     float64 `xml:"time,attr"`
	Title    string  `xml:"title,attr"`
	AttrText string  `xml:"text,attr"`
	Faction  string  `xml:"faction,attr"`
	Money    float64 `xml:"money,attr"`
	Ship     string
}

func Parse(file string) []Entry {

	// Open our xmlFile
	xmlFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	entries := make([]Entry, 0)
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		var e Entry

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "entry" {
				decoder.DecodeElement(&e, &se)
				if e.AttrText != "" {
					entries = append(entries, e)
				}
			}
		}
	}

	shipRegex := regexp.MustCompile("([A-Z]{3}\\-[0-9]{3})")
	for i := 0; i < len(entries); i++ {
		entries[i].Time = entries[i].Time / 60 / 60
		if entries[i].Money != 0 {
			entries[i].Money = entries[i].Money / 100
			ships := shipRegex.FindAll([]byte(entries[i].AttrText), -1)
			if ships != nil && len(ships) > GetConfig().ShipIndex {
				entries[i].Ship = string(ships[GetConfig().ShipIndex])
			} else {
				entries[i].Ship = ""
			}
		}
	}

	return entries
}

func containsTrigger(input string, triggers []string) bool {
	for _, trigger := range triggers {
		if strings.Contains(input, trigger) {
			return true
		}
	}
	return false
}
