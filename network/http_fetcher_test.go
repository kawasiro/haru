package network

import "testing"

// TODO network 테스트하기전에 서버를 어떻게 돌려둘지 생각해보자
// 테스트하기전에 로컬에 테스트 서버를 띄우는게 가능하면
// 127.0.0.1 을 이용해서 테스트가 가능할것이다
func TestFetch_HttpFetcher(t *testing.T) {
	f := HttpFetcher{"_cache"}

	cases := []struct {
		url  string
		code FetchCode
	}{
		// http
		{
			"http://www.google.com/robots.txt",
			FetchCodeSuccessHttp,
		},
		// https
		{
			"https://www.google.com/humans.txt",
			FetchCodeSuccessHttp,
		},
		// 404
		{
			"http://google.com/maybe-404",
			FetchCodeErrorNetwork404,
		},
		{
			"http://127.0.0.1:9999/",
			FetchCodeErrorNetworkNormal,
		},
		{
			"http://192.168.123.123:1234/",
			FetchCodeErrorNetworkTimeout,
		},
	}
	for _, c := range cases {
		r := f.Fetch(c.url)
		if r.Code != c.code {
			t.Errorf("FetchHttp - expected %d, got %d", c.code, r.Code)
		}
	}
}
