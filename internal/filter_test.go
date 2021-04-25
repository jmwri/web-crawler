package internal_test

import (
	"github.com/jmwri/web-crawler/internal"
	"net/url"
	"reflect"
	"testing"
)

func TestDedupeURLs(t *testing.T) {
	type args struct {
		urls []*url.URL
	}
	tests := []struct {
		name string
		args args
		want []*url.URL
	}{
		{
			name: "removes dupes",
			args: args{urls: []*url.URL{
				{Scheme: "https", Host: "test.com"},
				{Scheme: "https", Host: "test.com"},
				{Scheme: "https", Host: "test.com", Path: "about"},
				{Scheme: "https", Host: "test.com", Path: "contact"},
			}},
			want: []*url.URL{
				{Scheme: "https", Host: "test.com"},
				{Scheme: "https", Host: "test.com", Path: "about"},
				{Scheme: "https", Host: "test.com", Path: "contact"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := internal.DedupeURLs(tt.args.urls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DedupeURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterURLs(t *testing.T) {
	type args struct {
		urls    []*url.URL
		filters []internal.URLFilterFunc
	}
	tests := []struct {
		name string
		args args
		want []*url.URL
	}{
		{
			name: "dedupe and nonhttp",
			args: args{
				urls: []*url.URL{
					{Scheme: "https", Host: "test.com"},
					{Scheme: "https", Host: "test.com"},
					{Scheme: "mailto", Opaque: "test@localhost"},
				},
				filters: []internal.URLFilterFunc{
					internal.DedupeURLs, internal.RemoveNonHTTPURLs,
				},
			},
			want: []*url.URL{
				{Scheme: "https", Host: "test.com"},
			},
		},
		{
			name: "dedupe and same domain",
			args: args{
				urls: []*url.URL{
					{Scheme: "https", Host: "test.com"},
					{Scheme: "https", Host: "test.com"},
					{Scheme: "https", Host: "sub.test.com"},
				},
				filters: []internal.URLFilterFunc{
					internal.DedupeURLs, internal.SameDomainFilter(&url.URL{Scheme: "https", Host: "test.com"}),
				},
			},
			want: []*url.URL{
				{Scheme: "https", Host: "test.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := internal.FilterURLs(tt.args.urls, tt.args.filters...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveNonHTTPURLs(t *testing.T) {
	type args struct {
		urls []*url.URL
	}
	tests := []struct {
		name string
		args args
		want []*url.URL
	}{
		{
			name: "removes mailto",
			args: args{urls: []*url.URL{
				{Scheme: "mailto", Opaque: "test@localhost"},
				{Scheme: "https", Host: "test.com"},
				{Scheme: "http", Host: "test.com"},
			}},
			want: []*url.URL{
				{Scheme: "https", Host: "test.com"},
				{Scheme: "http", Host: "test.com"},
			},
		},
		{
			name: "removes tel",
			args: args{urls: []*url.URL{
				{Scheme: "tel", Opaque: "00000000000"},
				{Scheme: "https", Host: "test.com"},
				{Scheme: "http", Host: "test.com"},
			}},
			want: []*url.URL{
				{Scheme: "https", Host: "test.com"},
				{Scheme: "http", Host: "test.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := internal.RemoveNonHTTPURLs(tt.args.urls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveNonHTTPURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSameDomainFilter(t *testing.T) {
	type params struct {
		target *url.URL
	}
	type args struct {
		urls []*url.URL
	}
	tests := []struct {
		name   string
		params params
		args   args
		want   []*url.URL
	}{
		{
			name:   "removes different scheme",
			params: params{target: &url.URL{Scheme: "https", Host: "test.com"}},
			args: args{urls: []*url.URL{
				{Scheme: "https", Host: "test.com"},
				{Scheme: "http", Host: "test.com"},
			}},
			want: []*url.URL{
				{Scheme: "https", Host: "test.com"},
			},
		},
		{
			name:   "includes different paths",
			params: params{target: &url.URL{Scheme: "https", Host: "test.com"}},
			args: args{urls: []*url.URL{
				{Scheme: "https", Host: "test.com"},
				{Scheme: "https", Host: "test.com", Path: "contact"},
			}},
			want: []*url.URL{
				{Scheme: "https", Host: "test.com"},
				{Scheme: "https", Host: "test.com", Path: "contact"},
			},
		},
		{
			name:   "removes different host",
			params: params{target: &url.URL{Scheme: "https", Host: "test.com"}},
			args: args{urls: []*url.URL{
				{Scheme: "https", Host: "test.com"},
				{Scheme: "https", Host: "test.net"},
			}},
			want: []*url.URL{
				{Scheme: "https", Host: "test.com"},
			},
		},
		{
			name:   "removes different subdomain",
			params: params{target: &url.URL{Scheme: "https", Host: "test.com"}},
			args: args{urls: []*url.URL{
				{Scheme: "https", Host: "test.com"},
				{Scheme: "https", Host: "sub.test.com"},
			}},
			want: []*url.URL{
				{Scheme: "https", Host: "test.com"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := internal.SameDomainFilter(tt.params.target)
			if got := f(tt.args.urls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SameDomainFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
