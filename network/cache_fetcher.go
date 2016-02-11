package network

import (
	"compress/gzip"
	"io/ioutil"
	"os"
	"strings"
)

type CacheFileFetcher struct {
	CacheRootPath string
}

func (f *CacheFileFetcher) CacheFilePathIfExist(rawurl string) string {
	seg := ParseUrl(rawurl)
	cachePath := seg.ToCacheFilePath(f.CacheRootPath)
	candidatePaths := []string{
		cachePath + ".gz",
		cachePath,
	}
	for _, filepath := range candidatePaths {
		f, err := os.Open(filepath)
		defer f.Close()

		if err == nil {
			return filepath
		}
	}
	return ""
}

func (f *CacheFileFetcher) CacheExist(rawurl string) bool {
	if f.CacheFilePathIfExist(rawurl) == "" {
		return false
	} else {
		return true
	}
}

func (f *CacheFileFetcher) fetchCompressedCache(rawurl, filepath string) *FetchResult {
	// gzip
	file, _ := os.Open(filepath)
	defer file.Close()

	// decompress
	archive, err := gzip.NewReader(file)
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

func (f *CacheFileFetcher) fetchNormalCache(rawurl, filepath string) *FetchResult {
	// normal
	file, _ := os.Open(filepath)

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

func (f *CacheFileFetcher) Fetch(rawurl string) *FetchResult {
	filepath := f.CacheFilePathIfExist(rawurl)
	if strings.HasSuffix(filepath, ".gz") {
		return f.fetchCompressedCache(rawurl, filepath)
	} else if filepath != "" {
		return f.fetchNormalCache(rawurl, filepath)
	} else {
		// cache not found
		return NewErrorResult(rawurl, FetchCodeErrorCacheNotExist)
	}
}
