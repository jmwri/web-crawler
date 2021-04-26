# web-crawler
A simple web crawler

## Usage

```
Usage of crawler <target>:
  -h    show help
  -maxDepth int
        crawl up to this depth - 0 for no limit (default 4)
  -sameDomain
        only crawl the same domain (default true)
  -workers int
        number of workers (default 20)
```

## Improvements

* Add support for `rel="nofollow"`
* Check Content-Type when scraping pages
* Output progress as the crawler is running
* Add different output modes, ie sitemap
* Make it clear when errors impact results

## Tests

```
go test -race -cover ./...
```
