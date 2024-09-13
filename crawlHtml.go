package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
)

func crawlHtml(u string, maxConcurrent int, maxPages int) {
	fmt.Printf("starting crawl of: %s\n", u)
	visited := make(map[string]int, 0)
	baseUrl, _ := url.Parse(u)
	var wgp sync.WaitGroup
	channel := make(chan struct{}, maxConcurrent)
	cfg := config{
		pages:              visited,
		baseURL:            baseUrl,
		mu:                 &sync.Mutex{},
		concurrencyControl: channel,
		wg:                 &wgp,
		maxPages:           maxPages,
	}
	wgp.Add(1)
	go cfg.crawlPage(u)
	wgp.Wait()
	printVisied(visited, baseUrl.String())
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
	maxPages           int
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	// creat goroutine lock
	cfg.concurrencyControl <- struct{}{}
	// release the lock when exit
	defer func() {
		cfg.wg.Done()
		<-cfg.concurrencyControl
	}()
	reachedMax := cfg.checkLength()
	if reachedMax {
		return
	}
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

func printVisied(pages map[string]int, u string) {
	// Define the header and footer
	headerFooter := strings.Repeat("=", 29) // Repeat "=" 29 times
	reportTitle := fmt.Sprintf("REPORT for %s", u)

	// Print the header
	fmt.Println(headerFooter)
	fmt.Println(reportTitle)
	fmt.Println(headerFooter)

	sortedPages := sortMap(pages)
	for _, page := range sortedPages {
		fmt.Printf("Found %d internal links to %s\n", page.Value, page.Key)
	}
}

type KeyValue struct {
	Key   string
	Value int
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if _, found := cfg.pages[normalizedURL]; found {
		cfg.pages[normalizedURL]++
		return false
	}
	fmt.Printf("crawling : %s \n", normalizedURL)
	cfg.pages[normalizedURL]++
	return true
}

func (cfg *config) checkLength() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	return len(cfg.pages) >= cfg.maxPages
}

func sortMap(p map[string]int) []KeyValue {

	var res []KeyValue

	for k, v := range p {
		res = append(res, KeyValue{
			Key:   k,
			Value: v,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		if res[i].Value == res[j].Value {
			// If Values are equal, sort alphabetically by URL
			return res[i].Key < res[j].Key
		}
		// Sort by Value in descending order
		return res[i].Value > res[j].Value
	})

	//sort by
	return res
}
