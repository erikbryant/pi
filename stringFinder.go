package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/erikbryant/stringFinder/cache"
	"github.com/erikbryant/web"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	target = flag.String("target", "pi", "string to find")
)

// webRequest returns the contents of a REST API call
func webRequest(url string) (map[string]interface{}, error) {
	var response *http.Response
	var err error

	response, err = web.Request2(url, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("error fetching symbol data %s", err)
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("got an unexpected StatusCode %v", response)
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var jsonObject map[string]interface{}

	err = json.Unmarshal(contents, &jsonObject)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal json %s", err)
	}

	return jsonObject, nil
}

func fetch(url string) (map[string]interface{}, error) {
	response, err := cache.Read(url)
	if err == nil {
		return response, nil
	}

	time.Sleep(500 * time.Millisecond)
	response, err = webRequest(url)
	if err != nil {
		return nil, err
	}

	cache.Update(url, response)

	return response, nil
}

// getPi returns a chunk of digits of pi starting at startChunk
func getPi(startChunk int) string {
	// https://pi.delivery/

	length := 1000
	url := fmt.Sprintf("https://api.pi.delivery/v1/pi?start=%d&numberOfDigits=%d", startChunk*length, length)

	j, err := fetch(url)
	if err != nil {
		panic(err)
	}

	digits := j["content"].(string)

	return digits
}

// stringToDigits returns a string with each character replaced with its decimal value
func stringToDigits(str string) string {
	s := ""

	for _, ch := range str {
		s += fmt.Sprintf("%d", ch)
	}

	return s
}

func find(s string) {
	digits := stringToDigits(s)
	fmt.Printf("Searching PI for: %s -> %s\n", s, digits)
	for i := 12000; i <= 18000; i++ {
		pi := getPi(i)
		if strings.Contains(pi, digits) {
			offset := i*1000 + strings.Index(pi, digits)
			fmt.Printf("Found! Pattern %s -> %s starts at offset: %d\n", s, digits, offset)
		}
	}
}

func main() {
	fmt.Printf("Welcome to String Finder!\n\n")

	flag.Parse()

	find(*target)
	find(strings.ToUpper(*target))
}
