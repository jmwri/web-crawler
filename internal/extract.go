package internal

import (
	"errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"io"
	"net/url"
)

// ExtractorFunc extracts links from the given io.Reader
type ExtractorFunc func(r io.Reader) ([]*url.URL, error)

// HtmlTokenExtractor uses html.Tokenizer to extract links
func HtmlTokenExtractor(r io.Reader) ([]*url.URL, error) {
	t := html.NewTokenizer(r)

	links := make([]*url.URL, 0)
	seenLinks := map[string]bool{}

	for {
		tokenType := t.Next()
		if tokenType == html.ErrorToken {
			if errors.Is(t.Err(), io.EOF) {
				return links, nil
			}
			return nil, t.Err()
		}
		token := t.Token()

		// Only searching for anchor tags, which aren't usually self closing
		if token.Type != html.StartTagToken || token.DataAtom != atom.A {
			continue
		}
		link := hrefValue(token)
		if link == "" {
			continue
		}
		// Skip the link if we've already seen it
		if seen, ok := seenLinks[link]; ok && seen {
			continue
		}
		parsedUrl, err := url.Parse(link)
		if err != nil {
			// Skip links we're unable to parse
			continue
		}
		links = append(links, parsedUrl)
		seenLinks[link] = true
	}
}

// hrefValue returns the value of the href attribute
func hrefValue(t html.Token) string {
	for _, attr := range t.Attr {
		attrAtom := atom.Lookup([]byte(attr.Key))
		if attrAtom == atom.Href {
			return attr.Val
		}
	}
	return ""
}
