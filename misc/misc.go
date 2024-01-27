package misc

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/erikbryant/web"
	"io"
	"os"
	"path"
)

var (
	download = flag.Bool("download", false, "download new data")
)

const piDir = "../datafiles/pi/"

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

// getPi downloads a chunk of digits of pi starting at startDigit
func getPi(startDigit, length int) (map[string]interface{}, error) {
	// https://pi.delivery/

	url := fmt.Sprintf("https://api.pi.delivery/v1/pi?start=%d&numberOfDigits=%d", startDigit, length)

	response, err := webRequest(url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

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

func storePackedPi(file, digits string) error {
	object := path.Join(piDir, file)
	packed := packDigits(digits)
	return os.WriteFile(object, packed, 0644)
}

// downloadPi repeatedly requests more digits of Pi
func downloadPi() {
	length := 1000

	for i := 100001; i <= 100001; i++ {
		if i%10 == 0 {
			fmt.Println(i)
		}
		pi, err := getPi(i*length, length)
		if err != nil {
			panic(err)
		}
		file := fmt.Sprintf("start=%d&numberOfDigits=%d", i*length, length)
		digits := pi["content"].(string)
		err = storePackedPi(file, digits)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	fmt.Printf("Welcome to String Finder!\n\n")

	flag.Parse()

	if *download {
		downloadPi()
		return
	}

	fmt.Println("Invalid combination of arguments")
}
