package internal

import (
	"net/url"
	"strings"
)

// URLModifyFunc modifies the given URL in place
type URLModifyFunc func(u *url.URL)

// ModifyURL returns a modified copy of the given URL
func ModifyURL(u *url.URL, modifiers ...URLModifyFunc) *url.URL {
	// Don't want to modify the given URL, so we assign the value to a local var
	uValue := *u
	// We then use the pointer of the local var to modify the URL in place
	localURL := &uValue
	for _, f := range modifiers {
		f(localURL)
	}
	return localURL
}

// ModifyURLs returns modified copies of the given URLs
func ModifyURLs(urls []*url.URL, modifiers ...URLModifyFunc) []*url.URL {
	modified := make([]*url.URL, len(urls))
	for i, u := range urls {
		modified[i] = ModifyURL(u, modifiers...)
	}
	return modified
}

// RemoveTrailingSlash removes the trailing slash from the URLs path
func RemoveTrailingSlash(u *url.URL) {
	u.Path = strings.TrimSuffix(u.Path, "/")
}

// RemoveFragment removes the fragment from the URL
func RemoveFragment(u *url.URL) {
	u.Fragment = ""
}
