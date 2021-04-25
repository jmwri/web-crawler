package web_crawler

import (
	"github.com/jmwri/web-crawler/internal"
	"net/http"
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

// DefaultCrawler is the default Crawler
var DefaultCrawler Crawler

func init() {
	loader := internal.NewHTTPGetLoader(http.DefaultClient)
	DefaultCrawler = crawler{
		loader:    internal.LoaderWithRetry(loader.Load, internal.SimpleBackoff, 5),
		extractor: internal.HtmlTokenExtractor,
	}
}

// crawler is a wrapper around the internal Crawler. It means we can use public interfaces.
type crawler struct {
	loader    internal.LoaderFunc
	extractor internal.ExtractorFunc
}

// Crawl according to the specified options
func (c crawler) Crawl(target *url.URL, sameDomain bool, maxDepth int, workers int) (Result, error) {
	o := internal.NewCrawlOptions(target, sameDomain, maxDepth, workers)
	return internal.Crawl(c.loader, c.extractor, o)
}
