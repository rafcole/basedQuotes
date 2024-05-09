package sfox

import (
	"cryptoSnapShot/adapters"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type SFOX struct {
	Query adapters.Query
}

func (sfox SFOX) Authenticate() (int, error) {
	// sFox API is public, no need to authenticate
	fmt.Println("\n===> Status: Authenticating")
	return 1, nil
}

func (sfox SFOX) ValidatePair() (bool, error) {
	// avaliable pairs accessible through GET https://api.sfox.com/v1/currency-pairs
	// ideally we'd have that data pulled into a file ahead of time for validation

	formattedPair := sfox.FormattedCurrencyPair()

	req, err := http.NewRequest("GET", "https://api.sfox.com/v1/markets/currency-pairs", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	godotenv.Load()
	// fmt.Println("API key -> ", os.Getenv("SFOX_API_KEY"))

	// Add custom headers to the request
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SFOX_API_KEY")))

	// fmt.Println(req.Header)

	// Make the HTTP request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	// Parse the response body into a map[string]interface{}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	// Check if the pairStr exists in the response
	if pair, ok := data[formattedPair]; ok {
		fmt.Println("Pair found:", pair)
		return true, nil
	}

	// If pairStr not found in response
	fmt.Println("Pair not valid for sFox API:", formattedPair)
	return false, nil
}

// for sfox the API needs a simple concat of basequote -> "btc/usd" -> "btcusd"
func (sfox SFOX) FormattedCurrencyPair() string {
	return sfox.Query.Currency_Base + sfox.Query.Currency_Quote
}

func (sfox SFOX) FetchOHLCV(q adapters.Query) adapters.OHLCVData {
	sfox.Authenticate()

	var url = fmt.Sprintf("https://chartdata.sfox.com/candlesticks?endTime=%d&pair=%s&period=60&startTime=%d", sfox.Query.Time_stamp, sfox.FormattedCurrencyPair(), sfox.Query.StartTime())

	fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	fmt.Println("\n===> API Response status: ", resp.Status)
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Body to string: ", string(body))

	ohlcv, err := sfox.FormatOHLCV(resp.Body)
	if err != nil {
		log.Fatalln("Could not format OHLCV: ", err)
	}

	fmt.Println(ohlcv)

	return ohlcv
}

// takes slice of candlesticks, returns array of OHLCV structs
func (sfox SFOX) FormatOHLCV(response io.ReadCloser) (adapters.OHLCVData, error) {
	fmt.Println("===> Response recieved, parsing data")
	var dataArr []adapters.OHLCVData

	err := json.NewDecoder(response).Decode(&dataArr)
	if err != nil {
		fmt.Println("Error in FormatOHLCV")
		log.Fatalln(err)
	}

	fmt.Println(" ========> FormatOHLCV candlestick: ", dataArr)

	if len(dataArr) == 0 {
		return adapters.OHLCVData{}, errors.New("no data available for this period")
	}

	// This is a hardcode to extract the only element in the array
	// In a more sophisticated program it may be necessary to process
	// the response as a stream as detailed as detailed in encoder pkg
	// docs https://pkg.go.dev/encoding/json@go1.22.2#NewDecoder:~:text=an%20input%20stream.-,Example,-%C2%B6
	candlestick := dataArr[0]

	return candlestick, nil
}
