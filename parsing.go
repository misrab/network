package network

import (
	"bytes"

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