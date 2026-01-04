package main

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/gorilla/feeds"
	"github.com/rwinkhart/backtone/atom"
	"github.com/rwinkhart/backtone/bthtml"
	"github.com/rwinkhart/go-boilerplate/other"
)

const cfgPath = "/etc/backtone.xml"

// feedFieldInfo provides the string value for a given field in the parent feed
// and specifies which index in the regex capture group to extract the same field
// from per individual news item.
// An index of -1 can be used to use the same value as the parent feed.
// An index of -2 can be used to use a default (placeholder) value.
type feedFieldInfoT struct {
	Parent                string   `xml:"parent,omitempty"`
	CaptureGroupConnector []string `xml:"captureGroupConnector,omitempty"`
	CaptureGroupIndex     []int8   `xml:"captureGroupIndex,omitempty"`
}
type feedT struct {
	Title       feedFieldInfoT `xml:"title"`
	Link        feedFieldInfoT `xml:"link"`
	Description feedFieldInfoT `xml:"description"`
	AuthorName  feedFieldInfoT `xml:"authorName"`
	AuthorEmail feedFieldInfoT `xml:"authorEmail"`
}
type inputT struct {
	FlareSolverrURL string  `xml:"flareSolverrURL"`
	Regex           string  `xml:"regex"`
	LoadSeconds     float32 `xml:"loadSeconds"`
	MaxFeedItems    int     `xml:"maxFeedItems"`
	Feed            feedT   `xml:"feed"`
}

func main() {
	// handle input argument (XML)
	if len(os.Args) < 2 {
		demoInputBytes, _ := xml.MarshalIndent(
			inputT{
				FlareSolverrURL: "http://127.0.0.1:8191",
				Regex:           "PLACEHOLDER",
				LoadSeconds:     2.5,
				MaxFeedItems:    10,
				Feed: feedT{
					Title:       feedFieldInfoT{Parent: "Demo Title", CaptureGroupIndex: []int8{0}},
					Link:        feedFieldInfoT{Parent: "https://example.com/", CaptureGroupIndex: []int8{1}},
					Description: feedFieldInfoT{Parent: "Demo description", CaptureGroupConnector: []string{"Prefix - ", " - Middle - ", " - Suffix"}, CaptureGroupIndex: []int8{1, 3}},
					AuthorName:  feedFieldInfoT{Parent: "John Doe", CaptureGroupIndex: []int8{3}},
					AuthorEmail: feedFieldInfoT{Parent: "johndoe@example.com", CaptureGroupIndex: []int8{4}},
				},
			},
			"",
			"    ",
		)
		other.PrintError(os.Args[0]+" must be called with valid base64-encoded XML input using the following schema:\n\n"+string(demoInputBytes), 1)
	}

	// decode base64 input XML
	decodedXML, err := base64.StdEncoding.DecodeString(os.Args[1])
	if err != nil {
		other.PrintError("Failed to decode base64 input XML: "+err.Error(), 1)
	}

	var inputInfo inputT
	if err := xml.Unmarshal(decodedXML, &inputInfo); err != nil {
		other.PrintError("Failed to unmarshal input XML: "+err.Error(), 1)
	}

	rawHTML, err := bthtml.GetFromURL(inputInfo.FlareSolverrURL, inputInfo.Feed.Link.Parent, inputInfo.LoadSeconds)
	if err != nil {
		other.PrintError("Failed to get HTML content: "+err.Error(), 1)
	}

	atomXML, err := atom.GetFromHTML(
		&feeds.Feed{
			Title:       inputInfo.Feed.Title.Parent,
			Link:        &feeds.Link{Href: inputInfo.Feed.Link.Parent},
			Description: inputInfo.Feed.Description.Parent,
			Author:      &feeds.Author{Name: inputInfo.Feed.AuthorName.Parent, Email: inputInfo.Feed.AuthorEmail.Parent},
		},
		&rawHTML,
		inputInfo.Regex,
		atom.IndicesT{
			TitleC:       inputInfo.Feed.Title.CaptureGroupConnector,
			TitleI:       inputInfo.Feed.Title.CaptureGroupIndex,
			LinkC:        inputInfo.Feed.Link.CaptureGroupConnector,
			LinkI:        inputInfo.Feed.Link.CaptureGroupIndex,
			DescriptionC: inputInfo.Feed.Description.CaptureGroupConnector,
			DescriptionI: inputInfo.Feed.Description.CaptureGroupIndex,
			AuthorNameC:  inputInfo.Feed.AuthorName.CaptureGroupConnector,
			AuthorNameI:  inputInfo.Feed.AuthorName.CaptureGroupIndex,
			AuthorEmailC: inputInfo.Feed.AuthorEmail.CaptureGroupConnector,
			AuthorEmailI: inputInfo.Feed.AuthorEmail.CaptureGroupIndex,
		},
		inputInfo.MaxFeedItems,
	)
	if err != nil {
		other.PrintError("Failed to generate Atom RSS feed from HTML:\n\n"+rawHTML+"\n\n"+err.Error(), 1)
	}
	fmt.Println(*atomXML)
}
