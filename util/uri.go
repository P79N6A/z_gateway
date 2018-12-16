package util

import "strings"

func GetUriLastSeg(uri string) string {
	segs := strings.Split(uri, "/")
	return segs[len(segs) -1]
}