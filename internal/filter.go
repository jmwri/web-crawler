package internal

import (
	"net/url"
)

// URLFilterFunc returns filtered URLs
type URLFilterFunc func(urls []*url.URL) []*url.URL

// FilterURLs filters the given URLs through the given filters
func FilterURLs(urls []*url.URL, filters ...URLFilterFunc) []*url.URL {
	for _, f := range filters {
		urls = f(urls)
	}
	return urls
}

// DedupeURLs compares the string representation of the given URLs and removes duplicates
func DedupeURLs(urls []*url.URL) []*url.URL {
	seenURLs := make(map[string]bool)
	filtered := make([]*url.URL, 0)
	for _, u := range urls {
		if seen, ok := seenURLs[u.String()]; ok && seen {
			continue
		}
		filtered = append(filtered, u)
		seenURLs[u.String()] = true
	}
	return filtered
}

// RemoveNonHTTPURLs filters out URLs that do not have an HTTP scheme
func RemoveNonHTTPURLs(urls []*url.URL) []*url.URL {
	filtered := make([]*url.URL, 0)
	for _, u := range urls {
		if u.Scheme != "http" && u.Scheme != "https" {
			continue
		}
		filtered = append(filtered, u)
	}
	return filtered
}

// SameDomainFilter returns a URLFilterFunc that filters out URLs that do not have the same scheme and host
// as the given target
func SameDomainFilter(target *url.URL) URLFilterFunc {
	return func(urls []*url.URL) []*url.URL {
		filtered := make([]*url.URL, 0)
		for _, u := range urls {
			if u.Host != target.Host {
				continue
			}
			if u.Scheme != target.Scheme {
				continue
			}
			filtered = append(filtered, u)
		}
		return filtered
	}
}
