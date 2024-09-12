package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func crawlHtml(url string) {
	fmt.Printf("starting crawl of: %s", url)
	htmlbody, _ := getHtml(url)
	getURLsFromHTML(htmlbody, url)
	os.Exit(0)
}

func getHtml(url string) (string, error) {

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	resType := res.Header.Get("Content-Type")
	if res.StatusCode >= 400 {
		return "", fmt.Errorf("code > 299")
	} else if res.StatusCode > 200 {
		return "", fmt.Errorf(res.Status)
	}
	if resType == "text" || resType == "html" {
		return "", fmt.Errorf("not html or text")
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)

	return string(body), nil
}
