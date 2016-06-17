# PCRE Matcher Web Service

## Requirements

1. Go environment setup https://golang.org/doc/install
2. gcc `apt-get install gcc`
3. libpcre++-dev `apt-get install pcre++-dev`
4. Text file called *pcre.txt* with regexes. The format should be the following: "$regex\t$Context"
Example: `^(\"|').+?\\1$\tI am a context string`

## Installing

1. `go get github.com/ren-/pcre-matcher-service`
2. Locate `src/github.com/ren-/pcre-matcher-service` in your *$GOPATH*
3. `go install`
4. `go build`

Now you can run the executable `./pcre-matcher-service` to start the server. It will bind on port 5000.

## Example

Send a HTTP POST to http://localhost:5000/matchUrls with the following data:

```
{
  "urls":["http://google.com", "http://yahoo.com"]
}
```

## License

Licensed under GNU GENERAL PUBLIC LICENSE, please read the LICENSE file.
