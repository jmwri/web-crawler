package internal

import (
	"net/url"
	"sync"
)

// NewCrawlOptions returns new CrawlOptions with sanitised input
func NewCrawlOptions(target *url.URL, sameDomain bool, maxDepth int, workers int) CrawlOptions {
	if maxDepth < 0 {
		maxDepth = 0
	}
	if workers < 1 {
		workers = 10
	}
	return CrawlOptions{
		Target:     target,
		SameDomain: sameDomain,
		MaxDepth:   maxDepth,
		Workers:    workers,
	}
}

// CrawlOptions defines the options for Crawl
type CrawlOptions struct {
	// Target is the URL to start crawling from
	Target *url.URL
	// SameDomain restricts crawling to the same domain
	SameDomain bool
	// MaxDepth is the max depth that the crawler will go to. Set to 0 for no max depth.
	MaxDepth int
	// Workers is the amount of Workers that will be used to crawl. Must be at least 1.
	Workers int
}

// crawlRequest defines a single request that the crawler should perform
type crawlRequest struct {
	// origin is the URL that the the link is from
	origin *url.URL
	// target is the URL that the link is to
	target *url.URL
	// depth the depth that the page was discovered
	depth int
}

// next returns a new crawlRequest to the given target
func (r crawlRequest) next(target *url.URL) crawlRequest {
	return crawlRequest{
		origin: r.target,
		target: target,
		depth:  r.depth + 1,
	}
}

// crawlResponse is the response for a single request that the crawler performed
type crawlResponse struct {
	// request is the request the triggered the page being scraped
	request crawlRequest
	// err is present if the crawler failed to scrape the page
	err error
	// urls are the links found on the page
	urls []*url.URL
}

// Result contains all of the URLs discovered and visited
type Result struct {
	// options are the options that the Crawler was executed against
	options CrawlOptions
	// crawledURLs is a map of URL to URLs discovered on that page
	crawledURLs map[string][]string
	// mu is an internal mutex to ensure routine safe access of crawledURLs
	mu *sync.Mutex
}

// Store discovered URLs against a URL
func (r Result) Store(u *url.URL, urls []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.crawledURLs[u.String()] = urls
}

// Target is the URL that the crawler started at
func (r Result) Target() *url.URL {
	return r.options.Target
}

// SameDomain specifies if the crawler limited to the same domain
func (r Result) SameDomain() bool {
	return r.options.SameDomain
}

// MaxDepth returns the max depth of the crawler
func (r Result) MaxDepth() int {
	return r.options.MaxDepth
}

// URLs returns a map of URLs that the crawler visited, and a list of URLs found on that page
func (r Result) URLs() map[string][]string {
	return r.crawledURLs
}

// requestLog stores URLs that we have previously seen and issued requests for
type requestLog struct {
	// seenURLs contains URLs that we have already requested, or are trying to request
	seenURLs map[string]bool
	// mu is an internal mutex to ensure routine safe access of seenURLs
	mu *sync.Mutex
}

// Seen returns whether or not the given URL has been seen before
func (l requestLog) Seen(u *url.URL) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.seenURLs[u.String()]
}

// MarkAsSeen marks the given URL as seen
func (l requestLog) MarkAsSeen(u *url.URL) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.seenURLs[u.String()] = true
}

// buildFilters according the the provided options
func buildFilters(o CrawlOptions) []URLFilterFunc {
	filters := []URLFilterFunc{
		RemoveNonHTTPURLs,
		DedupeURLs,
	}
	if o.SameDomain {
		filters = append(filters, SameDomainFilter(o.Target))
	}
	return filters
}

// buildModifiers returns the standard modifiers
func buildModifiers() []URLModifyFunc {
	return []URLModifyFunc{
		RemoveTrailingSlash,
		RemoveFragment,
	}
}

// Crawl according to the specified options
func Crawl(loader LoaderFunc, extractor ExtractorFunc, o CrawlOptions) (Result, error) {
	filters := buildFilters(o)
	modifiers := buildModifiers()

	// Initialise our result
	res := Result{
		options:     o,
		crawledURLs: make(map[string][]string),
		mu:          &sync.Mutex{},
	}
	reqLog := requestLog{
		seenURLs: make(map[string]bool),
		mu:       &sync.Mutex{},
	}

	// wg tracks the number of URLs currently being processed
	wg := &sync.WaitGroup{}
	reqCh := make(chan crawlRequest)
	resCh := make(chan crawlResponse)

	for i := 0; i < o.Workers; i++ {
		go requestWorker(loader, extractor, reqCh, resCh, filters, modifiers)
	}
	for i := 0; i < o.Workers; i++ {
		go responseWorker(wg, o.MaxDepth, reqLog, res, reqCh, resCh)
	}

	// Queue the initial request
	wg.Add(1)
	reqLog.MarkAsSeen(o.Target)
	reqCh <- crawlRequest{
		target: o.Target,
		depth:  1,
	}
	// Wait for all URLs to be processed
	wg.Wait()
	// Close channels so Workers exit
	close(reqCh)
	close(resCh)

	return res, nil
}

// requestWorker creates a crawlResponse based on the crawlRequest and sends it to responseWorker
func requestWorker(loader LoaderFunc, extractor ExtractorFunc, reqCh <-chan crawlRequest, resCh chan<- crawlResponse, filters []URLFilterFunc, modifiers []URLModifyFunc) {
	for r := range reqCh {
		// Find the links on the page
		urls, err := scrapeURLs(loader, extractor, r)

		// Normalise and filter the URLs
		urls = ModifyURLs(urls, modifiers...)
		urls = FilterURLs(urls, filters...)

		// Send the URLs to the next worker
		resCh <- crawlResponse{
			request: r,
			err:     err,
			urls:    urls,
		}
	}
}

// responseWorker stores and initiates requests for scraped URLs
func responseWorker(wg *sync.WaitGroup, maxDepth int, reqLog requestLog, res Result, reqCh chan<- crawlRequest, resCh <-chan crawlResponse) {
	for r := range resCh {
		// Store the scraped URLs against the URL they were found on
		res.Store(r.request.target, urlsToString(r.urls))

		// Build the next set of requests
		next := nextRequests(r.request, r.urls)

		for _, n := range next {
			// Keep crawling until we reach max depth
			if maxDepth > 0 && n.depth > maxDepth {
				continue
			}
			// Skip links we've already seen somewhere else
			if reqLog.Seen(n.target) {
				continue
			}
			reqLog.MarkAsSeen(n.target)
			// Send the request back to requestWorker
			// Write to channel in a routine to avoid a deadlock
			wg.Add(1)
			go func(r crawlRequest) {
				reqCh <- r
			}(n)
		}
		wg.Done()
	}
}

// urlsToString converts a slice of url.URL to a slice of string
func urlsToString(urls []*url.URL) []string {
	res := make([]string, len(urls))
	for i, u := range urls {
		res[i] = u.String()
	}
	return res
}

// scrapeURLs from the requests Target
func scrapeURLs(loader LoaderFunc, extractor ExtractorFunc, req crawlRequest) ([]*url.URL, error) {
	// Load the page
	reader, err := loader(req.target.String())
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Extract anchor tags from the page
	urls, err := extractor(reader)
	if err != nil {
		return nil, err
	}

	// Build the URLs as references from the Target
	for i, u := range urls {
		urls[i] = req.target.ResolveReference(u)
	}

	return urls, nil
}

// nextRequests builds the next crawlRequest for each URL based on the current request
func nextRequests(req crawlRequest, urls []*url.URL) []crawlRequest {
	reqs := make([]crawlRequest, 0)
	for _, u := range urls {
		reqs = append(reqs, req.next(u))
	}
	return reqs
}
