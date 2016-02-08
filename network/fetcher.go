package network

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	FetchCodeSuccessHttp        = 1
	FetchCodeSuccessUseCompress = 2
	FetchCodeSuccessUseCache    = 3

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
	return r.Code < 100
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

type CacheFileFetcher struct {
	CacheRootPath string
}
type ProxyFetcher struct {
	CacheRootPath string
}

func (r *FetchResult) Size() int64 {
	return int64(len(r.Data))
}

func (f *CacheFileFetcher) Fetch(rawurl string) *FetchResult {
	seg := ParseUrl(rawurl)
	cachePath := seg.ToCacheFilePath(f.CacheRootPath)

	// gzip 압축된것이 있는지 확인
	compressedCachePath := cachePath + ".gz"
	compressedFile, err := os.Open(compressedCachePath)
	if err == nil {
		// decompress
		archive, err := gzip.NewReader(compressedFile)
		if err != nil {
			panic(err)
		}
		defer archive.Close()

		data, err := ioutil.ReadAll(archive)
		if err != nil {
			panic(err)
		}
		return &FetchResult{
			Url:  rawurl,
			Data: data,
			Date: archive.Header.ModTime,
			Code: FetchCodeSuccessUseCompress,
		}
	}

	file, err := os.Open(cachePath)
	if err != nil {
		// 에러에 따라서 분기
		return NewErrorResult(rawurl, FetchCodeErrorCacheNotExist)
	}

	// normal state
	fileInfo, _ := file.Stat()
	data := make([]byte, fileInfo.Size())
	file.Read(data)

	return &FetchResult{
		Url:  rawurl,
		Data: data,
		Date: fileInfo.ModTime(),
		Code: FetchCodeSuccessUseCache,
	}
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
