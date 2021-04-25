package internal_test

import (
	"errors"
	"fmt"
	"github.com/jmwri/web-crawler/internal"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

// testSiteHandler makes the html in testdata available over a test http client
var testSiteHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(fmt.Sprintf("./testdata/html/%s", r.URL.Path))
	defer f.Close()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
})

func TestHTTPGetLoader_Load(t *testing.T) {
	ts := httptest.NewServer(testSiteHandler)
	defer ts.Close()

	client := ts.Client()

	type fields struct {
		client *http.Client
	}
	type args struct {
		p string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "can load index",
			fields:  fields{client: client},
			args:    args{p: ts.URL + "/index.html"},
			want:    "./testdata/html/index.html",
			wantErr: false,
		},
		{
			name:    "can load about",
			fields:  fields{client: client},
			args:    args{p: ts.URL + "/about.html"},
			want:    "./testdata/html/about.html",
			wantErr: false,
		},
		{
			name:    "errors on 404",
			fields:  fields{client: client},
			args:    args{p: ts.URL + "/nope"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := internal.NewHTTPGetLoader(tt.fields.client)
			got, err := l.Load(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				// Don't try to load any files if we expect an error
				return
			}

			wantReader, err := os.Open(tt.want)
			if err != nil {
				t.Fatalf("failed to read wanted file %s: %s", tt.want, err)
			}

			gotContents, err := ioutil.ReadAll(got)
			if err != nil {
				t.Fatalf("failed to read all of file %s: %s", tt.args.p, err)
			}

			wantContents, err := ioutil.ReadAll(wantReader)
			if err != nil {
				t.Fatalf("failed to read all of wanted file %s: %s", tt.want, err)
			}

			if !reflect.DeepEqual(gotContents, wantContents) {
				t.Errorf("Load() got = %s, want %s", gotContents, wantContents)
			}
		})
	}
}

func TestLoaderWithRetry_RightNumberOfAttempts(t *testing.T) {
	var loaderCalls int

	countingLoader := func(p string) (io.ReadCloser, error) {
		loaderCalls++
		return nil, errors.New("always fail")
	}
	noBackoff := func(attempt int) time.Duration { return 0 }

	loader := internal.LoaderWithRetry(countingLoader, noBackoff, 5)
	_, err := loader("test page")

	if err == nil {
		t.Error("expecting err")
	}
	if loaderCalls != 5 {
		t.Errorf("expecting %d calls, got %d", 5, loaderCalls)
	}
}

func TestLoaderWithRetry_SleepsCorrectly(t *testing.T) {
	failingLoader := func(p string) (io.ReadCloser, error) {
		return nil, errors.New("always fail")
	}
	smallBackoff := func(attempt int) time.Duration { return (time.Millisecond * 100) * time.Duration(attempt) }

	loader := internal.LoaderWithRetry(failingLoader, smallBackoff, 5)

	timeBefore := time.Now()
	_, err := loader("test page")
	duration := time.Since(timeBefore)

	if err == nil {
		t.Error("expecting err")
	}

	// Total duration should be:
	// 100ms * 1 = 100ms
	// 100ms * 2 = 200ms
	// 100ms * 3 = 300ms
	// 100ms * 4 = 400ms
	// 100ms * 5 = 500ms
	// 			 = 1500ms
	expectedDuration := time.Millisecond * 1500
	// Allow an extra 5% of time
	upperBound := expectedDuration + (expectedDuration / 20)
	if duration < expectedDuration || duration > upperBound {
		t.Errorf("expecting to sleep from %s to %s, slept for %s", expectedDuration, upperBound, duration)
	}
}
