# book-search

## Web Crawler

### Environment Variables
```
CRAWLER_THREADS=8
CRAWLER_USER_AGENT=BOOK-SEARCH-CRAWLER

REDIS_HOST=localhost:6379
REDIS_PASSWORD=1q2w3e4r
REDIS_DB=0
```

### Run
```
go mod tidy
go run webcrawler/main.go
```
