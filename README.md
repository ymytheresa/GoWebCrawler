# GoWebCrawler

GoWebCrawler is a command-line interface (CLI) application I'm writing in Go that generates an "internal links" report for any website on the internet. It crawls each page of the site, analyzing the internal linking structure.

## Project Overview

In this project, I'm building a Web Crawler in Golang to help website owners understand and optimize their internal linking structure. Good internal linking is crucial for SEO (Search Engine Optimization) as it helps search engines understand the structure and hierarchy of a website's content.

## Features

- Crawls websites and analyzes their internal link structure
- Generates a report of internal links for each page
- Handles relative and absolute URLs
- Respects `robots.txt` rules (to be implemented)
- Provides a user-friendly CLI interface

## Learning Goals

Through this project, I aim to:

- Gain hands-on practice with local Go development and tooling
- Learn to make HTTP requests in Go
- Parse HTML with Golang
- Practice unit testing in Go
- Understand web crawling concepts and best practices

## Getting Started

### Prerequisites

- Go 1.16 or higher

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/GoWebCrawler.git
   ```
2. Navigate to the project directory:
   ```
   cd GoWebCrawler
   ```
3. Build the project:
   ```
   go build
   ```

