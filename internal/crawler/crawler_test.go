package crawler

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		u       string
		workers int
	}
	tests := []struct {
		name    string
		args    args
		want    *Crawler
		wantErr bool
	}{
		{
			"no workers",
			args{
				"https://test",
				0,
			},
			&Crawler{
				u:       &url.URL{
					Scheme: "https",
					Host: "test",
				},
				workers: 1,
				client: http.DefaultClient,
			},
			false,
		},
		{
			"parsed url",
			args{
				"https://test",
				3,
			},
			&Crawler{
				u:       &url.URL{
					Scheme:     "https",
					Host:       "test",
				},
				workers: 3,
				client: http.DefaultClient,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.u, tt.args.workers)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

// Quick and dirty response modifier
type rtf func(r *http.Request) *http.Response

func (f rtf) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r), nil
}

func TestCrawler_scan(t *testing.T) {
	type fields struct {
		u       *url.URL
		workers int
		client  *http.Client
	}
	type args struct {
		u *url.URL
	}
	tests := []struct {
		name    string
		html    string
		fields  fields
		args    args
		want    []*url.URL
		wantErr bool
	}{
		{
			"links",
			`<html><head></head><body><a href="https://test/helloworld">Hello World</a><a href="https://test/helloworld2">Hello World 2</a></body></html>`,
			fields{
				u:       &url.URL{
					Scheme: "https",
					Host: "test",
				},
				workers: 1,
			},
			args{
				u: &url.URL{
					Scheme:     "https",
					Host:       "test",
				},
			},
			[]*url.URL{
				{
					Scheme:     "https",
					Host:       "test",
					Path:       "/helloworld",
				},
				{
					Scheme:     "https",
					Host:       "test",
					Path:       "/helloworld2",
				},
			},
			false,
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				u:       tt.fields.u,
				workers: tt.fields.workers,
				client:  &http.Client{
					Transport: rtf(func(r *http.Request) *http.Response {
						return &http.Response{
							StatusCode: 200,
							Body: ioutil.NopCloser(bytes.NewBufferString(tt.html)),
						}
					}),
				},
			}
			got, err := c.scan(tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("got length != want length")
			}

			for i, j := 0, len(got); i < j; i++ {
				if !reflect.DeepEqual(*got[i], *tt.want[i]) {
					t.Errorf("scan() got = %v, want %v", *got[i], *tt.want[i])
				}
			}
		})
	}
}