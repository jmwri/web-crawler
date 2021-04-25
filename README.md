# web-crawler
A simple web crawler

## Usage

```
Usage of crawler <target>:
  -h    show help
  -maxDepth int
        crawl up to this depth (default 2)
  -sameDomain
        only crawl the same domain (default true)
  -workers int
        number of workers (default 10)
```

## Improvements

* Add support for `rel="nofollow"`
* Output results as the program runs rather than at the end
* Make it clear when errors impact results
* Check Content-Type when scraping pages

## Tests

```
go test ./...
```
