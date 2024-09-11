package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	//convert raw base url to a url struct
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	//create io.Reader
	htmlReader := strings.NewReader(htmlBody)
	//get html modes from html reader
	doc, err := html.Parse(htmlReader)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse HTML: %v", err)
	}

	//recursive function to loop through the doc nodes and get all the hrefs from tag a
	var urls []string
	var traverseNodes func(*html.Node)
	traverseNodes = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, anchor := range node.Attr {
				if anchor.Key == "href" {
					href, err := url.Parse(anchor.Val)
					if err != nil {
						fmt.Printf("couldn't parse href '%v': %v\n", anchor.Val, err)
						continue
					}

					resolvedURL := baseURL.ResolveReference(href)
					//will add base host to relative paths, if href is abs link it will not change and just return the abs link
					urls = append(urls, resolvedURL.String())
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverseNodes(child)
		}
	}
	traverseNodes(doc)

	return urls, nil
}
