# Go FIM scraper

Example scripts to scrape GPS locations from your own private FIMTrack account.

This project is not related to the FIM imaging field.

## Disclaimer

The software is distributed as is, and is created with learning purposes, it is not directly related to the original
publisher of the webpage, thus, you may use it without any guarantee and at your own risk.

## Running it

You need to have a set of credentials from the service, as it is the only information you will be able to have access to.

Replace them in the `main.go` definition and execute the following commands to build and run.

```
go mod download
go build main.go
./main
```