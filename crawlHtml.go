package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func crawlHtml(url string) {
	fmt.Printf("starting crawl of: %s\n", url)
	visited := make(map[string]int, 0)
	crawlPage(url, url, visited)
	printVisied(visited)
	os.Exit(0)
}

func getHTML(rawURL string) (string, error) {
	res, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("got Network error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode > 399 {
		return "", fmt.Errorf("got HTTP error: %s", res.Status)
	}

	contentType := res.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("got non-HTML response: %s", contentType)
	}

	htmlBodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("couldn't read response body: %v", err)
	}

	htmlBody := string(htmlBodyBytes)

	return htmlBody, nil
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	norCurrentUrl, _ := normalizeURL(rawCurrentURL)
	if _, found := pages[norCurrentUrl]; found {
		pages[norCurrentUrl]++
		return
	}
	if getUrlDomain(rawBaseURL) != getUrlDomain(rawCurrentURL) {
		return
	}
	pages[norCurrentUrl]++
	htmlBody, _ := getHTML(rawCurrentURL)
	urls, _ := getURLsFromHTML(htmlBody, rawCurrentURL)
	for _, u := range urls {
		crawlPage(rawBaseURL, u, pages)
	}
}

func printVisied(pages map[string]int) {
	for k, v := range pages {
		fmt.Printf("%s : %d\n", k, v)
	}
}
