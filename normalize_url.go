package main

import (
	"errors"
	"net/url"
	"strings"
)

func normalizeURL(rawUrl string) (string, error) {
	//[scheme:][//[userinfo@]host][/]path[?query][#fragment]
	urlStruct, err := url.Parse(rawUrl)
	if err != nil {
		return "", errors.New("bad url")
	}
	fullPath := urlStruct.Host + urlStruct.Path
	fullPath = strings.ToLower(fullPath)
	fullPath = strings.TrimSuffix(fullPath, "/")
	return fullPath, nil
}
