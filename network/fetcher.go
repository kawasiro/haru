package network

import (
	"bytes"
	"net/http"
	"os"
	"time"
)

const (
	// http code와 값을 공유하려고 2xx 코드에 겹침
	FetchCodeSuccessHttp        = 299
	FetchCodeSuccessUseCompress = 298
	FetchCodeSuccessUseCache    = 297

	// custom error code
	FetchCodeErrorCacheNotExist  = 1001
	FetchCodeErrorUnknown        = 1002
	FetchCodeErrorNetworkNormal  = 1003
	FetchCodeErrorNetworkTimeout = 1004

	// from http status code
	FetchCodeErrorNetwork404 = http.StatusNotFound
)

type FetchCode int32

const (
	FetcherTypeProxy = 1
	FetcherTypeHttp  = 2
	FetcherTypeCache = 3
)

type FetcherType int32

type Fetcher interface {
	Fetch(rawurl string) *FetchResult
}

type FetchResult struct {
	Url  string
	Data []byte
	Date time.Time
	Code FetchCode
}

func (r *FetchResult) IsSuccess() bool {
	return 200 <= r.Code && r.Code < 300
}

func (r *FetchResult) String() string {
	buf := bytes.NewBuffer(r.Data)
	return buf.String()
}

func (r *FetchResult) SaveToFile(dstFilePath string) {
	// TODO: check file existence first with io.IsExist
	output, err := os.Create(dstFilePath)
	if err != nil {
		panic(err)
	}
	defer output.Close()
	output.Write(r.Data)
}

const defaultCacheRootPath = "_cache"

func NewHttpFetcher() Fetcher {
	return NewFetcher(FetcherTypeHttp, "")
}

func NewDefaultFetcher() Fetcher {
	return NewFetcher(FetcherTypeProxy, defaultCacheRootPath)
}
func NewCacheFetcher() CacheFileFetcher {
	return CacheFileFetcher{defaultCacheRootPath}
}

func NewFetcher(fetcherType FetcherType, cacheRootPath string) Fetcher {
	switch fetcherType {
	case FetcherTypeProxy:
		return &ProxyFetcher{cacheRootPath}
	case FetcherTypeCache:
		return &CacheFileFetcher{cacheRootPath}
	case FetcherTypeHttp:
		return &HttpFetcher{cacheRootPath}
	default:
		panic("unknown fetcher type")
	}
}

func NewErrorResult(rawurl string, code FetchCode) *FetchResult {
	return &FetchResult{
		Url:  rawurl,
		Data: []byte{},
		Date: time.Now(),
		Code: code,
	}
}

type ProxyFetcher struct {
	CacheRootPath string
}

func (r *FetchResult) Size() int64 {
	return int64(len(r.Data))
}

func (f *ProxyFetcher) Fetch(rawurl string) *FetchResult {
	cacheFetcher := CacheFileFetcher{f.CacheRootPath}
	result := cacheFetcher.Fetch(rawurl)
	if result.IsSuccess() {
		return result
	}

	httpFetcher := HttpFetcher{f.CacheRootPath}
	return httpFetcher.Fetch(rawurl)
}
