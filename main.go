package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
)

type Hit struct {
	Line  int
	URLID int
}

// urls to be processed
type Urls struct {
	Urls []string `json:"urls"`
}

type FormattedHit struct {
	MatchedRegex string `json:"regex"`
	Url          string `json:"url"`
}

var lines []string
var compiledRegexes []pcre.Matcher

func main() {
	// 200-300 regexes
	file, err := os.Open("pcre.txt")
	if err != nil {
		panic("Regexes should be available in the pcre.txt file, format: regex\\tcontext")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	// Pre-allocate space for regexes to save some time
	compiledRegexes = make([]pcre.Matcher, 0, len(lines))

	// generate regexp instances for each regex
	for i := 0; i < len(lines); i++ {
		// Get the the regex
		parts := strings.Split(lines[i], "\t")
		reg, err := pcre.Compile(parts[0], pcre.UTF8)
		if err != nil {
			fmt.Printf("Invalid regexp [%s]\n", parts[0])
			fmt.Println(err)
			continue
		}
		matcher := reg.MatcherString("", pcre.NOTEMPTY)
		compiledRegexes = append(compiledRegexes, *matcher)
	}
	http.HandleFunc("/matchUrls", parseUrls)
	log.Fatal(http.ListenAndServe(":5000", nil))
	return
}

func findMatches(urlList Urls) []Hit {
	var results []Hit
	for i := 0; i < len(compiledRegexes); i++ {
		for j := 0; j < len(urlList.Urls); j++ {

			// Test if the URL matches the regex
			result := compiledRegexes[i].MatchString(urlList.Urls[j], pcre.NOTEMPTY)

			//If it matches, return the result over the channel
			if result {
				results = append(results, Hit{
					Line:  i,
					URLID: j,
				})
			}
		}
	}
	return results

}

func parseUrls(rw http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var data Urls
	err := decoder.Decode(&data)

	if err != nil {
		fmt.Println(err)
	}

	hits := findMatches(data)
	results := make([]FormattedHit, 0, len(hits))
	for _, element := range hits {
		results = append(results, FormattedHit{lines[element.Line], data.Urls[element.URLID]})
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(results)
	fmt.Printf("%+v", results)
}
