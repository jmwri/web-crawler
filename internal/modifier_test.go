package internal_test

import (
	"github.com/jmwri/web-crawler/internal"
	"net/url"
	"reflect"
	"testing"
)

func TestModifyURL(t *testing.T) {
	type args struct {
		u         *url.URL
		modifiers []internal.URLModifyFunc
	}
	tests := []struct {
		name       string
		args       args
		wantArgURL *url.URL
		want       *url.URL
	}{
		{
			name: "no mod if none specified",
			args: args{
				u:         &url.URL{Scheme: "https", Host: "test.com", Path: "page/", Fragment: "section"},
				modifiers: []internal.URLModifyFunc{},
			},
			wantArgURL: &url.URL{Scheme: "https", Host: "test.com", Path: "page/", Fragment: "section"},
			want:       &url.URL{Scheme: "https", Host: "test.com", Path: "page/", Fragment: "section"},
		},
		{
			name: "runs combined modifiers",
			args: args{
				u:         &url.URL{Scheme: "https", Host: "test.com", Path: "page/", Fragment: "section"},
				modifiers: []internal.URLModifyFunc{internal.RemoveTrailingSlash, internal.RemoveFragment},
			},
			wantArgURL: &url.URL{Scheme: "https", Host: "test.com", Path: "page/", Fragment: "section"},
			want:       &url.URL{Scheme: "https", Host: "test.com", Path: "page"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := internal.ModifyURL(tt.args.u, tt.args.modifiers...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModifyURL() = %v, want %v", got, tt.want)
			}
			// Ensuring the URL we used as arg hasn't changed
			if !reflect.DeepEqual(tt.args.u, tt.wantArgURL) {
				t.Errorf("ModifyURL() modified URL inplace = %v, want %v", tt.args.u, tt.wantArgURL)
			}
		})
	}
}

func TestModifyURLs(t *testing.T) {
	type args struct {
		urls      []*url.URL
		modifiers []internal.URLModifyFunc
	}
	tests := []struct {
		name        string
		args        args
		wantArgURLs []*url.URL
		want        []*url.URL
	}{
		{
			name: "no mod if none specified",
			args: args{
				urls: []*url.URL{
					{Scheme: "https", Host: "test.com", Path: "page/", Fragment: "section"},
					{Scheme: "https", Host: "test.com", Path: "about"},
				},
				modifiers: []internal.URLModifyFunc{},
			},
			wantArgURLs: []*url.URL{
				{Scheme: "https", Host: "test.com", Path: "page/", Fragment: "section"},
				{Scheme: "https", Host: "test.com", Path: "about"},
			},
			want: []*url.URL{
				{Scheme: "https", Host: "test.com", Path: "page/", Fragment: "section"},
				{Scheme: "https", Host: "test.com", Path: "about"},
			},
		},
		{
			name: "runs combined modifiers",
			args: args{
				urls: []*url.URL{
					{Scheme: "https", Host: "test.com", Path: "page/", Fragment: "section"},
					{Scheme: "https", Host: "test.com", Path: "about"},
					{Scheme: "https", Host: "test.com", Path: "about/"},
				},
				modifiers: []internal.URLModifyFunc{internal.RemoveTrailingSlash, internal.RemoveFragment},
			},
			wantArgURLs: []*url.URL{
				{Scheme: "https", Host: "test.com", Path: "page/", Fragment: "section"},
				{Scheme: "https", Host: "test.com", Path: "about"},
				{Scheme: "https", Host: "test.com", Path: "about/"},
			},
			want: []*url.URL{
				{Scheme: "https", Host: "test.com", Path: "page"},
				{Scheme: "https", Host: "test.com", Path: "about"},
				{Scheme: "https", Host: "test.com", Path: "about"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := internal.ModifyURLs(tt.args.urls, tt.args.modifiers...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModifyURLs() = %v, want %v", got, tt.want)
			}
			// Ensuring the URLs we used as arg hasn't changed
			if !reflect.DeepEqual(tt.args.urls, tt.wantArgURLs) {
				t.Errorf("ModifyURLs() modified URLs inplace = %v, want %v", tt.args.urls, tt.wantArgURLs)
			}
		})
	}
}

func TestRemoveFragment(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name string
		args args
		want *url.URL
	}{
		{
			name: "removes fragment",
			args: args{u: &url.URL{Scheme: "https", Host: "test.com", Fragment: "section"}},
			want: &url.URL{Scheme: "https", Host: "test.com"},
		},
		{
			name: "no effect when no fragment",
			args: args{u: &url.URL{Scheme: "https", Host: "test.com", Fragment: ""}},
			want: &url.URL{Scheme: "https", Host: "test.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			internal.RemoveFragment(tt.args.u)
			if !reflect.DeepEqual(tt.args.u, tt.want) {
				t.Errorf("RemoveFragment() = %v, want %v", tt.args.u, tt.want)
			}
		})
	}
}

func TestRemoveTrailingSlash(t *testing.T) {
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name string
		args args
		want *url.URL
	}{
		{
			name: "removes trailing slash",
			args: args{u: &url.URL{Scheme: "https", Host: "test.com", Path: "/"}},
			want: &url.URL{Scheme: "https", Host: "test.com"},
		},
		{
			name: "no effect when no slash",
			args: args{u: &url.URL{Scheme: "https", Host: "test.com", Path: ""}},
			want: &url.URL{Scheme: "https", Host: "test.com"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			internal.RemoveTrailingSlash(tt.args.u)
			if !reflect.DeepEqual(tt.args.u, tt.want) {
				t.Errorf("RemoveTrailingSlash() = %v, want %v", tt.args.u, tt.want)
			}
		})
	}
}
