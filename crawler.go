package web_crawler

import (
	"net/url"
)

// Crawler is a basic web crawler
type Crawler interface {
	// Crawl according to the specified options
	Crawl(target *url.URL, sameDomain bool, maxDepth int, workers int) (Result, error)
}

// Result is the output of Crawler
type Result interface {
	// Target is the URL that the crawler started on
	Target() *url.URL
	// SameDomain specifies if the crawler was limited to the same domain
	SameDomain() bool
	// MaxDepth returns the depth that the crawler was limited to
	MaxDepth() int
	// URLs returns a map of URLs that the crawler visited, and a list of URLs found on that page
	URLs() map[string][]string
}
