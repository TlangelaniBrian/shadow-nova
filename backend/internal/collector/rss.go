package collector

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

type RSSFeed struct {
	XMLName xml.Name  `xml:"rss"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title string    `xml:"title"`
	Items []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	// YouTube specific
	VideoID     string `xml:"videoId"` 
}

type AtomFeed struct {
	XMLName xml.Name   `xml:"feed"`
	Title   string     `xml:"title"`
	Entries []AtomEntry `xml:"entry"`
}

type AtomEntry struct {
	Title     string    `xml:"title"`
	Link      AtomLink  `xml:"link"`
	Summary   string    `xml:"summary"`
	Published string    `xml:"published"`
	Updated   string    `xml:"updated"`
	// YouTube specific
	VideoID   string    `xml:"videoId"`
	Group     MediaGroup `xml:"group"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
}

type MediaGroup struct {
	Description string `xml:"description"`
	Thumbnail   MediaThumbnail `xml:"thumbnail"`
}

type MediaThumbnail struct {
	URL string `xml:"url,attr"`
}

// FetchFeed fetches and parses an RSS or Atom feed
var FetchFeed = func(url string) ([]ContentMetadata, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d", resp.StatusCode)
	}

	// Try parsing as RSS first (most common)
	// We need to read body once, so we might need to buffer it if we want to try multiple parsers
	// Or just try to decode.
	// For simplicity, let's decode into a struct that handles both or check root element.
	// Actually, let's just try RSS, if empty try Atom.
	
	// Better approach: Decode into a generic map or check the first few bytes.
	// But `xml.Decoder` is stream based.
	
	// Let's assume we know the type or try one.
	// YouTube feeds are Atom. Blog feeds are often RSS.
	
	var items []ContentMetadata
	
	// We'll use a decoder and look at the root element
	decoder := xml.NewDecoder(resp.Body)
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "rss" {
				var feed RSSFeed
				decoder.DecodeElement(&feed, &se)
				for _, item := range feed.Channel.Items {
					pubDate, _ := parseDate(item.PubDate)
					items = append(items, ContentMetadata{
						Title:       item.Title,
						Description: item.Description,
						URL:         item.Link,
						ImageURL:    "", // RSS usually puts image in description HTML or enclosure
						PublishedAt: pubDate,
					})
				}
				return items, nil
			} else if se.Name.Local == "feed" {
				var feed AtomFeed
				decoder.DecodeElement(&feed, &se)
				for _, entry := range feed.Entries {
					pubDate, _ := parseDate(entry.Published)
					if pubDate.IsZero() {
						pubDate, _ = parseDate(entry.Updated)
					}
					
					desc := entry.Summary
					img := ""
					if entry.Group.Description != "" {
						desc = entry.Group.Description
					}
					if entry.Group.Thumbnail.URL != "" {
						img = entry.Group.Thumbnail.URL
					}

					items = append(items, ContentMetadata{
						Title:       entry.Title,
						Description: desc,
						URL:         entry.Link.Href,
						ImageURL:    img,
						PublishedAt: pubDate,
					})
				}
				return items, nil
			}
		}
	}

	return nil, fmt.Errorf("unknown feed format")
}

func parseDate(dateStr string) (time.Time, error) {
	// Try various formats
	formats := []string{
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		"Mon, 02 Jan 2006 15:04:05 -0700",
	}
	
	for _, f := range formats {
		if t, err := time.Parse(f, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Now(), fmt.Errorf("could not parse date")
}
