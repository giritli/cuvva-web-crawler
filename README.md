# Cuvva Web Crawler
Time-boxed myself 3 hours for this challenge, ran out of time to write tests for `Crawler.Crawl`.
I am currently using a Windows machine therefore I couldn't create and test a Makefile for this project.
I have listed relevant commands below.

## Building
```
go build -o crawler ./cmd/crawler/crawler.go
```

## Running
```
./crawler -url=https://cuvva.com
```

## Piping output
```
./crawler -url=https://cuvva.com > sitemap.json
```

## Testing
```
go test -count=1 -v ./...
```