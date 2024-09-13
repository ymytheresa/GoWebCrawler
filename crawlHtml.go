package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func crawlHtml(u string) {
	fmt.Printf("starting crawl of: %s\n", u)
	visited := make(map[string]int, 0)
	baseUrl, _ := url.Parse(u)
	var wgp sync.WaitGroup
	channel := make(chan struct{}, 10)
	cfg := config{
		pages:              visited,
		baseURL:            baseUrl,
		mu:                 &sync.Mutex{},
		concurrencyControl: channel,
		wg:                 &wgp,
	}
	wgp.Add(1)
	go cfg.crawlPage(u)
	wgp.Wait()
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

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	defer cfg.wg.Done()
	// creat goroutine lock
	cfg.concurrencyControl <- struct{}{}
	// release the lock when exit
	defer func() {
		<-cfg.concurrencyControl
	}()
	norCurrentUrl, _ := normalizeURL(rawCurrentURL)
	if getUrlDomain(cfg.baseURL.String()) != getUrlDomain(rawCurrentURL) {
		return
	}
	if !cfg.addPageVisit(norCurrentUrl) {
		return
	}

	htmlBody, _ := getHTML(rawCurrentURL)
	urls, _ := getURLsFromHTML(htmlBody, rawCurrentURL)
	for _, u := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(u)
	}
}

func printVisied(pages map[string]int) {
	for k, v := range pages {
		fmt.Printf("%s : %d\n", k, v)
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, found := cfg.pages[normalizedURL]; found {
		cfg.pages[normalizedURL]++
		return false
	}
	cfg.pages[normalizedURL]++
	return true
}
