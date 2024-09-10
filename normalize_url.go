package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
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

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	reader := strings.NewReader(htmlBody)
	htmlNodes, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	hrefs := getHrefs(htmlNodes)

	rawBaseURL, _ = normalizeURL(rawBaseURL)
	res, err := convertToAllAbsUrl(hrefs, rawBaseURL)
	return res, err
}

func getHrefs(node *html.Node) []string {
	var hrefs []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					hrefs = append(hrefs, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
	return hrefs
}

func convertToAllAbsUrl(hrefs []string, rawUrl string) ([]string, error) {
	res := make([]string, len(hrefs))
	host := "https://" + getHost(rawUrl)
	fmt.Println(host)
	var err error
	for _, h := range hrefs {
		h, err = convertToAbs(h, host)
		res = append(res, h)
	}
	return res, err
}

func getHost(raw string) string {
	s, _ := url.Parse(raw)
	fmt.Println(s.Host)
	return s.Host
}

func convertToAbs(href string, host string) (string, error) {
	if strings.HasPrefix(href, "http") {
		return normalizeURL(href)
	} else {
		return normalizeURL((host + href))
	}
}
