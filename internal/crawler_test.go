package internal_test

import (
	"fmt"
	"github.com/jmwri/web-crawler/internal"
	"io"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestCrawl(t *testing.T) {
	var loader internal.LoaderFunc = func(p string) (io.ReadCloser, error) {
		u, err := url.Parse(p)
		if err != nil {
			return nil, err
		}

		filename := u.Path
		if u.Host == "github.com" {
			filename = strings.ReplaceAll(u.Path, "/", "_")
			filename = "gh" + filename + ".html"
		}

		f, err := os.Open(fmt.Sprintf("./testdata/html/%s", filename))
		return f, err
	}

	type args struct {
		loader    internal.LoaderFunc
		extractor internal.ExtractorFunc
		o         internal.CrawlOptions
	}
	tests := []struct {
		name           string
		args           args
		wantTarget     *url.URL
		wantSameDomain bool
		wantMaxDepth   int
		wantURLs       map[string][]string
		wantErr        bool
	}{
		{
			name: "same domain",
			args: args{
				loader:    loader,
				extractor: internal.HtmlTokenExtractor,
				o:         internal.NewCrawlOptions(&url.URL{Scheme: "https", Host: "localhost", Path: "index.html"}, true, 0, 0),
			},
			wantTarget:     &url.URL{Scheme: "https", Host: "localhost", Path: "index.html"},
			wantSameDomain: true,
			wantMaxDepth:   0,
			wantURLs: map[string][]string{
				"https://localhost/index.html": {
					"https://localhost/index.html",
					"https://localhost/about.html",
					"https://localhost/contact.html",
				},
				"https://localhost/about.html": {
					"https://localhost/index.html",
					"https://localhost/about.html",
					"https://localhost/contact.html",
				},
				"https://localhost/contact.html": {
					"https://localhost/index.html",
					"https://localhost/about.html",
					"https://localhost/contact.html",
				},
			},
			wantErr: false,
		},
		{
			name: "same domain max depth",
			args: args{
				loader:    loader,
				extractor: internal.HtmlTokenExtractor,
				o:         internal.NewCrawlOptions(&url.URL{Scheme: "https", Host: "localhost", Path: "index.html"}, true, 1, 0),
			},
			wantTarget:     &url.URL{Scheme: "https", Host: "localhost", Path: "index.html"},
			wantSameDomain: true,
			wantMaxDepth:   1,
			wantURLs: map[string][]string{
				"https://localhost/index.html": {
					"https://localhost/index.html",
					"https://localhost/about.html",
					"https://localhost/contact.html",
				},
			},
			wantErr: false,
		},
		{
			name: "all domains",
			args: args{
				loader:    loader,
				extractor: internal.HtmlTokenExtractor,
				o:         internal.NewCrawlOptions(&url.URL{Scheme: "https", Host: "localhost", Path: "index.html"}, false, 0, 0),
			},
			wantTarget:     &url.URL{Scheme: "https", Host: "localhost", Path: "index.html"},
			wantSameDomain: false,
			wantMaxDepth:   0,
			wantURLs: map[string][]string{
				"https://localhost/index.html": {
					"https://localhost/index.html",
					"https://localhost/about.html",
					"https://localhost/contact.html",
					"https://github.com/jmwri/web-crawler",
				},
				"https://localhost/about.html": {
					"https://localhost/index.html",
					"https://localhost/about.html",
					"https://localhost/contact.html",
					"https://github.com/jmwri/web-crawler",
					"https://test.com/some/page",
					"https://test.com/cool-page",
					"https://test.com",
					"https://test.com/duplicate",
				},
				"https://localhost/contact.html": {
					"https://localhost/index.html",
					"https://localhost/about.html",
					"https://localhost/contact.html",
					"https://github.com/jmwri/web-crawler",
				},
				"https://github.com/jmwri/web-crawler": {
					"https://github.com/jmwri/web-crawler/commits",
					"https://github.com/jmwri/web-crawler/releases",
					"https://github.com/jmwri/web-crawler/pulls",
				},
				"https://github.com/jmwri/web-crawler/commits": {
					"https://github.com/jmwri/web-crawler/commits",
					"https://github.com/jmwri/web-crawler/releases",
					"https://github.com/jmwri/web-crawler/pulls",
					"https://github.com/jmwri/web-crawler/commits/aaaa",
					"https://github.com/jmwri/web-crawler/commits/bbbb",
					"https://github.com/jmwri/web-crawler/commits/cccc",
				},
				"https://github.com/jmwri/web-crawler/commits/aaaa": {},
				"https://github.com/jmwri/web-crawler/commits/bbbb": {},
				"https://github.com/jmwri/web-crawler/commits/cccc": {},
				"https://github.com/jmwri/web-crawler/pulls": {
					"https://github.com/jmwri/web-crawler/commits",
					"https://github.com/jmwri/web-crawler/releases",
					"https://github.com/jmwri/web-crawler/pulls",
					"https://github.com/jmwri/web-crawler/pulls/1",
					"https://github.com/jmwri/web-crawler/pulls/2",
					"https://github.com/jmwri/web-crawler/pulls/3",
				},
				"https://github.com/jmwri/web-crawler/pulls/1": {},
				"https://github.com/jmwri/web-crawler/pulls/2": {},
				"https://github.com/jmwri/web-crawler/pulls/3": {},
				"https://github.com/jmwri/web-crawler/releases": {
					"https://github.com/jmwri/web-crawler/commits",
					"https://github.com/jmwri/web-crawler/releases",
					"https://github.com/jmwri/web-crawler/pulls",
					"https://github.com/jmwri/web-crawler/releases/1.0.0",
					"https://github.com/jmwri/web-crawler/releases/2.0.0",
					"https://github.com/jmwri/web-crawler/releases/3.0.0",
				},
				"https://github.com/jmwri/web-crawler/releases/1.0.0": {},
				"https://github.com/jmwri/web-crawler/releases/2.0.0": {},
				"https://github.com/jmwri/web-crawler/releases/3.0.0": {},
				"https://test.com/some/page":                          {},
				"https://test.com/cool-page":                          {},
				"https://test.com":                                    {},
				"https://test.com/duplicate":                          {},
			},
			wantErr: false,
		},
		{
			name: "all domains max depth",
			args: args{
				loader:    loader,
				extractor: internal.HtmlTokenExtractor,
				o:         internal.NewCrawlOptions(&url.URL{Scheme: "https", Host: "localhost", Path: "index.html"}, false, 2, 0),
			},
			wantTarget:     &url.URL{Scheme: "https", Host: "localhost", Path: "index.html"},
			wantSameDomain: false,
			wantMaxDepth:   2,
			wantURLs: map[string][]string{
				"https://localhost/index.html": {
					"https://localhost/index.html",
					"https://localhost/about.html",
					"https://localhost/contact.html",
					"https://github.com/jmwri/web-crawler",
				},
				"https://localhost/about.html": {
					"https://localhost/index.html",
					"https://localhost/about.html",
					"https://localhost/contact.html",
					"https://github.com/jmwri/web-crawler",
					"https://test.com/some/page",
					"https://test.com/cool-page",
					"https://test.com",
					"https://test.com/duplicate",
				},
				"https://localhost/contact.html": {
					"https://localhost/index.html",
					"https://localhost/about.html",
					"https://localhost/contact.html",
					"https://github.com/jmwri/web-crawler",
				},
				"https://github.com/jmwri/web-crawler": {
					"https://github.com/jmwri/web-crawler/commits",
					"https://github.com/jmwri/web-crawler/releases",
					"https://github.com/jmwri/web-crawler/pulls",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := internal.Crawl(tt.args.loader, tt.args.extractor, tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crawl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				// Don't check result if we got an error
				return
			}
			if !reflect.DeepEqual(got.Target(), tt.wantTarget) {
				t.Errorf("Crawl().Target() got = %v, want %v", got.Target(), tt.wantTarget)
			}
			if !reflect.DeepEqual(got.SameDomain(), tt.wantSameDomain) {
				t.Errorf("Crawl().SameDomain() got = %v, want %v", got.SameDomain(), tt.wantSameDomain)
			}
			if !reflect.DeepEqual(got.MaxDepth(), tt.wantMaxDepth) {
				t.Errorf("Crawl().MaxDepth() got = %v, want %v", got.MaxDepth(), tt.wantMaxDepth)
			}
			if !reflect.DeepEqual(got.URLs(), tt.wantURLs) {
				t.Errorf("Crawl().URLs() got = %v, want %v", got.URLs(), tt.wantURLs)
			}
		})
	}
}

func TestNewCrawlOptions(t *testing.T) {
	type args struct {
		target     *url.URL
		sameDomain bool
		maxDepth   int
		workers    int
	}
	tests := []struct {
		name string
		args args
		want internal.CrawlOptions
	}{
		{
			name: "max depth at least 0",
			args: args{
				target:     &url.URL{Path: "/"},
				sameDomain: true,
				maxDepth:   -1,
				workers:    1,
			},
			want: internal.CrawlOptions{
				Target:     &url.URL{Path: "/"},
				SameDomain: true,
				MaxDepth:   0,
				Workers:    1,
			},
		},
		{
			name: "workers defaults to 10 if < 1",
			args: args{
				target:     &url.URL{Path: "/"},
				sameDomain: true,
				maxDepth:   0,
				workers:    0,
			},
			want: internal.CrawlOptions{
				Target:     &url.URL{Path: "/"},
				SameDomain: true,
				MaxDepth:   0,
				Workers:    10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := internal.NewCrawlOptions(tt.args.target, tt.args.sameDomain, tt.args.maxDepth, tt.args.workers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCrawlOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
