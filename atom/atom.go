package atom

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

// Indices indicates which capture group in the
// regex-parsed HTML to reference for each Atom field.
// An empty slice will use a blank/default value.
// For the slices of strings, these will be appended between
// captured group text.
type IndicesT struct {
	TitleC       []string
	TitleI       []int8
	LinkC        []string
	LinkI        []int8
	DescriptionC []string
	DescriptionI []int8
	AuthorNameC  []string
	AuthorNameI  []int8
	AuthorEmailC []string
	AuthorEmailI []int8
}

func GetFromString(feed *feeds.Feed, input *string, regexString string, infoIndices IndicesT, maxFeedItems int) (*string, error) {
	now := time.Now()
	r := regexp.MustCompile(regexString)
	newsItems := r.FindAllString(*input, maxFeedItems)
	for i, item := range newsItems {
		info := r.FindStringSubmatch(item)
		link := stitchFields(info, infoIndices.LinkC, infoIndices.LinkI)
		feedItem := &feeds.Item{
			Title:       stitchFields(info, infoIndices.TitleC, infoIndices.TitleI),
			Link:        &feeds.Link{Href: link},
			Id:          "urn:backtoneid:" + strconv.Itoa(i) + ":" + link,
			Description: stitchFields(info, infoIndices.DescriptionC, infoIndices.DescriptionI),
			Author: &feeds.Author{
				Name:  stitchFields(info, infoIndices.AuthorNameC, infoIndices.AuthorNameI),
				Email: stitchFields(info, infoIndices.AuthorEmailC, infoIndices.AuthorEmailI),
			},
			Created: now,
		}
		feed.Items = append(feed.Items, feedItem)
	}
	atom, err := feed.ToAtom()
	if err != nil {
		return nil, errors.New("unable to create Atom RSS feed: " + err.Error())
	}
	return &atom, nil
}

func stitchFields(info []string, connectors []string, indices []int8) string {
	var output strings.Builder
	if len(connectors) > 0 {
		for i := range connectors {
			output.WriteString(connectors[i])
			if len(indices) >= i+1 {
				output.WriteString(info[indices[i]])
			}
		}
	} else {
		for i := range indices {
			if i > 0 {
				output.WriteString(" | ")
			}
			output.WriteString(info[indices[i]])
		}
	}
	return interpretUnicode(output.String())
}

func interpretUnicode(s string) string {
	// match \uXXXX patterns
	r := regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)
	result := r.ReplaceAllStringFunc(s, func(match string) string {
		// extract the hex value
		hexStr := match[2:] // skip the \u part
		code, err := strconv.ParseInt(hexStr, 16, 32)
		if err != nil {
			return match // if parsing fails, return original
		}
		return string(rune(code))
	})
	return result
}
