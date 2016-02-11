package gallery

import (
	"net/url"
	"strings"
)

// url.QueryEscape를 쓰면 될것처럼 보이지만
// ' '의 경우 '%20' 대신 '+'로 바뀌는 문제가 있다
// 이런것만 간단히 예외처리
// https://github.com/golang/go/issues/4013
var simpleUrlEncodeTable = []struct {
	decoded string
	encoded string
}{
	{" ", "%20"},
	{":", url.QueryEscape(":")},
}

func UrlEncode(val string) string {
	result := val
	for _, t := range simpleUrlEncodeTable {
		result = strings.Replace(result, t.decoded, t.encoded, -1)
	}
	return result
}

func UrlDecode(val string) string {
	result := val
	for _, t := range simpleUrlEncodeTable {
		result = strings.Replace(result, t.encoded, t.decoded, -1)
	}
	return result
}
