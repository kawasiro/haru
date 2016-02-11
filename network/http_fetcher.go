package network

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"
)

type HttpFetcher struct {
	CacheRootPath string
}

func (f *HttpFetcher) saveCompressedCacheFile(result *FetchResult) {
	// Save to cache
	seg := ParseUrl(result.Url)

	cacheDir := seg.ToCacheDir(f.CacheRootPath)
	os.MkdirAll(cacheDir, 0755)

	// compress
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(result.Data)
	w.Close()

	// write file
	cacheFile := seg.ToCacheFilePath(f.CacheRootPath)
	cacheFile += ".gz"

	file, err := os.OpenFile(cacheFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	file.Write(b.Bytes())
	file.Close()
}

func (f *HttpFetcher) saveNormalCacheFile(result *FetchResult) {
	// Save to cache
	seg := ParseUrl(result.Url)

	cacheDir := seg.ToCacheDir(f.CacheRootPath)
	os.MkdirAll(cacheDir, 0755)

	// write file
	cacheFile := seg.ToCacheFilePath(f.CacheRootPath)
	file, err := os.OpenFile(cacheFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	file.Write(result.Data)
	file.Close()
}

func (f *HttpFetcher) handleError(err error, rawurl string) *FetchResult {
	// http://stackoverflow.com/questions/22761562/portable-way-to-detect-different-kinds-of-network-error-in-golang
	switch err.(type) {
	case *url.Error:
		switch err.(*url.Error).Err.(type) {
		case *net.OpError:
			return NewErrorResult(rawurl, FetchCodeErrorNetworkNormal)

		default:
			// err.Err가 *http.httpError 인 경우가 있는데
			// private라서 쓸수없는 클래스다. 그래서 이름으로 분기
			if reflect.TypeOf(err.(*url.Error).Err).String() == "*http.httpError" {
				return NewErrorResult(rawurl, FetchCodeErrorNetworkTimeout)
			}
		}
	}
	fmt.Printf("unknown error type: %Q\n", err)
	return NewErrorResult(rawurl, FetchCodeErrorUnknown)
}

func (f *HttpFetcher) fetch(client *http.Client, rawurl string) *FetchResult {
	response, err := client.Get(rawurl)
	if err != nil {
		return f.handleError(err, rawurl)
	}
	defer response.Body.Close()

	createSuccessFunc := func() *FetchResult {
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		return &FetchResult{
			Url:  rawurl,
			Data: buf.Bytes(),
			Date: time.Now(),
			Code: FetchCodeSuccessHttp,
		}
	}

	switch response.StatusCode {
	case 404:
		return NewErrorResult(rawurl, FetchCodeErrorNetwork404)
	case 200:
		return createSuccessFunc()
	default:
		return createSuccessFunc()
	}
}

func (f *HttpFetcher) Fetch(rawurl string) *FetchResult {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}

	// http://stackoverflow.com/questions/16895294/how-to-set-timeout-for-http-get-requests-in-golang
	// 오래 걸리는 작업은 워커로 넘길거니까 타임아웃을 굳이 설정할 필요는 없을듯

	// http://stackoverflow.com/questions/12122159/golang-how-to-do-a-https-request-with-bad-certificate
	client := &http.Client{}
	switch u.Scheme {
	case "http":
		client = &http.Client{}
	case "https":
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{
			Transport: tr,
		}
	default:
		// 없으면 http로 취급
		panic("unknown scheme")
	}
	result := f.fetch(client, rawurl)

	// 실패한 요청의 캐시는 저장할 필요 없다
	if result.IsSuccess() && f.CacheRootPath != "" {
		// 이미지 파일은 이미 압축된 상태일테니 또 압축할 필요없다
		if strings.HasSuffix(result.Url, ".jpg") || strings.HasSuffix(result.Url, ".png") {
			f.saveNormalCacheFile(result)
		} else {
			f.saveCompressedCacheFile(result)
		}
	}

	if result.IsSuccess() {
		log.Printf("HttpFetcher: %s -> success\n", rawurl)
	} else {
		log.Printf("HttpFetcher: %s -> fail\n", rawurl)
	}
	return result
}
