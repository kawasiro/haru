package network

import (
	"reflect"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	f := NewFetcher(FetcherTypeCache, "testdata")

	cases := []struct {
		url    string
		result FetchResult
	}{
		// simple
		{
			"http://a.com/foo/bar/file.html",
			FetchResult{
				"http://a.com/foo/bar/file.html",
				[]byte{'h', 'e', 'l', 'l', 'o', '\n'},
				time.Now(),
				FetchCodeSuccessUseCache,
			},
		},
		// use gzip
		{
			"http://a.com/foo/bar/compress.html",
			FetchResult{
				"http://a.com/foo/bar/compress.html",
				[]byte{'h', 'e', 'l', 'l', 'o', '\n'},
				time.Now(),
				FetchCodeSuccessUseCompress,
			},
		},
		// cache file not exists
		{
			"http://a.com/foo/bar/not-exist.html",
			FetchResult{
				"http://a.com/foo/bar/not-exist.html",
				[]byte{},
				time.Now(),
				FetchCodeErrorCacheNotExist,
			},
		},
	}

	for _, c := range cases {
		r := f.Fetch(c.url)
		// 날짜는 비교할 필요 없다
		r.Date = c.result.Date
		if !reflect.DeepEqual(*r, c.result) {
			t.Errorf("Fetch - expected %q, got %q", c.result, r)
		}
	}
}
