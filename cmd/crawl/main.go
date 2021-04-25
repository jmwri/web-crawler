package main

import (
	"flag"
	"fmt"
	webcrawler "github.com/jmwri/web-crawler"
	"net/url"
	"os"
)

func main() {
	sameDomainPtr := flag.Bool("sameDomain", true, "only crawl the same domain")
	maxDepthPtr := flag.Int("maxDepth", 2, "crawl up to this depth")
	workersPtr := flag.Int("workers", 10, "number of workers")
	helpPtr := flag.Bool("h", false, "show help")

	flag.Parse()

	if *helpPtr {
		flag.Usage()
		os.Exit(0)
	}
	if len(os.Args) != 2 {
		flag.Usage()
		os.Exit(1)
	}
	target := os.Args[1]

	parsedUrl, _ := url.Parse(target)
	res, err := webcrawler.DefaultCrawler.Crawl(parsedUrl, *sameDomainPtr, *maxDepthPtr, *workersPtr)
	if err != nil {
		panic(err)
	}

	for page, links := range res.URLs() {
		if links == nil {
			continue
		}
		fmt.Printf("%s links to:\n", page)
		for _, link := range links {
			fmt.Printf("- %s\n", link)
		}
	}

	fmt.Printf("crawled %d pages\n", len(res.URLs()))
}
