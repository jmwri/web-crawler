package internal_test

import (
	"fmt"
	"github.com/jmwri/web-crawler/internal"
	"net/url"
	"os"
	"reflect"
	"testing"
)

// testExtractorFunc is an implementation independent test for an ExtractorFunc
func testExtractorFunc(t *testing.T, f internal.ExtractorFunc) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    []*url.URL
		wantErr bool
	}{
		{
			name: "index",
			args: args{filename: "index.html"},
			want: []*url.URL{
				{
					Path: "index.html",
				},
				{
					Path: "about.html",
				},
				{
					Path: "contact.html",
				},
				{
					Scheme: "https",
					Host:   "github.com",
					Path:   "/jmwri/web-crawler",
				},
			},
			wantErr: false,
		},
		{
			name: "contact",
			args: args{filename: "contact.html"},
			want: []*url.URL{
				{
					Path: "index.html",
				},
				{
					Path: "about.html",
				},
				{
					Path: "contact.html",
				},
				{
					Scheme: "https",
					Host:   "github.com",
					Path:   "/jmwri/web-crawler",
				},
				{
					Scheme: "mailto",
					Opaque: "test@localhost",
				},
				{
					Scheme: "tel",
					Opaque: "00000000000",
				},
			},
			wantErr: false,
		},
		{
			name: "about",
			args: args{filename: "about.html"},
			want: []*url.URL{
				{
					Path: "index.html",
				},
				{
					Path: "about.html",
				},
				{
					Path: "contact.html",
				},
				{
					Scheme: "https",
					Host:   "github.com",
					Path:   "/jmwri/web-crawler",
				},
				{
					Scheme:   "https",
					Host:     "test.com",
					Path:     "/some/page/",
					Fragment: "section",
				},
				{
					Scheme: "https",
					Host:   "test.com",
					Path:   "/some/page/",
				},
				{
					Scheme: "https",
					Host:   "test.com",
					Path:   "/cool-page",
				},
				{
					Scheme: "https",
					Host:   "test.com",
				},
				{
					Host: "test.com",
				},
				{
					Scheme: "https",
					Host:   "test.com",
					Path:   "/duplicate",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(fmt.Sprintf("./testdata/html/%s", tt.args.filename))
			defer file.Close()
			if err != nil {
				t.Errorf("failed to open test fixture: %s", err)
				return
			}
			got, err := f(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extract() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHtmlTokenExtractor(t *testing.T) {
	testExtractorFunc(t, internal.HtmlTokenExtractor)
}
