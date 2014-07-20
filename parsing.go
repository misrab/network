package network

import (
	"bytes"
	"strings"

	"html/template"
	"encoding/xml"
)




// ParseRss takes a []byte assumed to be RSS XML, 
// and decodes it into an object.
func ParseRss(body []byte) (Query, error) {
	var q Query
	q = new(Rss)

	reader := bytes.NewReader(body)
	d := xml.NewDecoder(reader)
	err := d.Decode(q)

	return q, err
}

type Query interface {
	Items() []RssItem
}

func (r *Rss) Items() []RssItem {
	return r.Channel.Items
}

type Rss struct {
	Rss xml.Name 			`xml:"rss"`
	Channel struct {
		Items []RssItem		`xml:"item"`
	} 						`xml:"channel"`
}

type RssItem struct {
	Title string 			`xml:"title"`
	Link string 			`xml:"link"`
	Description string 		`xml:"description"`
	PubDate string 			`xml:"pubDate"`
}



// from https://github.com/kennygrant/sanitize/blob/master/sanitize.go
// Strip html tags, replace common entities, and escape <>&;'" in the result.
// Note the returned text may contain entities as it is escaped by HTMLEscapeString, and most entities are not translated.
func RemoveHtml(s string) (output string) {

	output = ""

	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		output = s
	} else {

		// First remove line breaks etc as these have no meaning outside html tags (except pre)
		// this means pre sections will lose formatting... but will result in less uninentional paras.
		s = strings.Replace(s, "\n", "", -1)

		// Then replace line breaks with newlines, to preserve that formatting
		s = strings.Replace(s, "</p>", "\n", -1)
		s = strings.Replace(s, "<br>", "\n", -1)
		s = strings.Replace(s, "</br>", "\n", -1)
		s = strings.Replace(s, "<br/>", "\n", -1)

		// Walk through the string removing all tags
		b := bytes.NewBufferString("")
		inTag := false
		for _, r := range s {
			switch r {
			case '<':
				inTag = true
			case '>':
				inTag = false
			default:
				if !inTag {
					b.WriteRune(r)
				}
			}
		}
		output = b.String()
	}

	// In case we have missed any tags above, escape the text - removes <, >, &, ' and ". 
	output = template.HTMLEscapeString(output)

	// Remove a few common harmless entities, to arrive at something more like plain text
	// This relies on having removed *all* tags above
	output = strings.Replace(output, "&nbsp;", " ", -1)
	output = strings.Replace(output, "&quot;", "\"", -1)
	output = strings.Replace(output, "&apos;", "'", -1)
	output = strings.Replace(output, "&#34;", "\"", -1)
	output = strings.Replace(output, "&#39;", "'", -1)
	// NB spaces here are significant - we only allow & not part of entity
	output = strings.Replace(output, "&amp; ", "& ", -1)
	output = strings.Replace(output, "&amp;amp; ", "& ", -1)

	return output
}



