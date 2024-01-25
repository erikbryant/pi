package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/erikbryant/stringFinder/cache"
	"github.com/erikbryant/web"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

var (
	download = flag.Bool("download", false, "download new data")
	pack     = flag.Bool("pack", false, "pack any downloaded data")
	target   = flag.String("target", "", "string to find")
)

// v-------------------------------- Downloader --------------------------------v

// webRequest returns the contents of a REST API call
func webRequest(url string) (map[string]interface{}, error) {
	response, err := web.Request2(url, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("error fetching data %s", err)
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

// getPi downloads a chunk of digits of pi starting at startChunk
func getPi(startChunk int) {
	// https://pi.delivery/

	length := 1000
	url := fmt.Sprintf("https://api.pi.delivery/v1/pi?start=%d&numberOfDigits=%d", startChunk*length, length)

	_, err := fetch(url)
	if err != nil {
		panic(err)
	}
}

// downloadPi repeatedly requests more digits of Pi
func downloadPi() {
	for i := 70000; i <= 80000; i++ {
		if i%10 == 0 {
			fmt.Println(i)
		}
		getPi(i)
	}
}

// v--------------------------------  Packer  --------------------------------v

// packDigits returns a packed representation of the given digits
func packDigits(digits string) []byte {
	count := len(digits)

	if count == 0 || count%2 != 0 {
		panic("packDigits failure! a non-zero, even number of digits is required!")
	}

	packed := make([]byte, count/2)

	for i := 0; i < count/2; i++ {
		packed[i] = (digits[i*2]-'0')<<4 | (digits[1+i*2] - '0')
	}

	return packed
}

const cacheDir = "./web-request-cache/"
const piDir = "./pi/"
const archiveDir = "./archive/"

func inFiles() []os.FileInfo {

	file, err := os.Open(cacheDir)
	if err != nil {
		panic(err)
	}

	files, err := file.Readdir(-1)
	if err != nil {
		panic(err)
	}

	return files
}

func readPi(file string) (string, error) {
	response, err := cache.Read(file)
	if err != nil {
		return "", fmt.Errorf("no data found for %s", file)
	}
	return response["content"].(string), nil
}

func storePi(file, digits string) error {
	object := path.Join(piDir, file)
	packed := packDigits(digits)
	return os.WriteFile(object, packed, 0644)
}

// packPi converts any downloaded cache files to packed files
func packPi() {
	for _, file := range inFiles() {
		inFile := file.Name()
		digits, err := readPi(inFile)
		if err != nil {
			panic(err)
		}
		// https:--api.pi.delivery-v1-pi-start=0-numberOfDigits=1000
		outFile := strings.Replace(inFile, "https:--api.pi.delivery-v1-pi-", "", 1)
		err = storePi(outFile, digits)
		if err != nil {
			panic(err)
		}
		archiveFile := path.Join(archiveDir, inFile)
		cacheFile := path.Join(cacheDir, inFile)
		err = os.Rename(cacheFile, archiveFile)
		if err != nil {
			panic(err)
		}
	}
}

// v--------------------------------  Searcher  --------------------------------v

// stringToDigits returns a string with each character replaced with its decimal value
func stringToDigits(str string) string {
	s := ""

	for _, ch := range str {
		s += fmt.Sprintf("%d", ch)
	}

	return s
}

func find(s string) {
	//digits := stringToDigits(s)
	//fmt.Printf("Searching PI for: %s -> %s\n", s, digits)
	//for i := 0; i <= 63000; i++ {
	//	pi, err := readPi(i)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	if strings.Contains(pi, digits) {
	//		offset := i*1000 + strings.Index(pi, digits)
	//		fmt.Printf("Found! Pattern %s -> %s starts at offset: %d\n", s, digits, offset)
	//	}
	//}
}

func main() {
	fmt.Printf("Welcome to String Finder!\n\n")

	flag.Parse()

	if *download {
		downloadPi()
		return
	}

	if *pack {
		packPi()
		return
	}

	if *target != "" {
		find(*target)
		find(strings.ToUpper(*target))
		return
	}

	fmt.Println("Invalid combination of arguments")
}
