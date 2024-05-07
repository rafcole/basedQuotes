package adapters

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func AdapterFactory(q Query) Venue {
	switch q.Venue {
	case "sfox":
		// validation here?
		fmt.Println("Found sfox in AdapterFactory")
		return SFOX{query: q}
	default:
		return nil
	}
}

// func Authenticate() {
// 	// need try/except blocks here
// 	fmt.Println("sFOX candlestick API does not require authentication")
// }

// interface
type Venue interface {
	Authenticate() (int, error)
	ValidatePair(string) (bool, error)
	FetchOHLCV(Query) []byte // Timestamp within method or at time of request?
	FormatOHLCV([]byte) any  // Last touch point
}

type SFOX struct {
	query Query
}

func (sfox SFOX) Authenticate() (int, error) {
	fmt.Println("sFOX candlestick API does not require authentication")
	return 1, nil
}

func (sfox SFOX) ValidatePair(pairStr string) (bool, error) {
	// TODO dict of 3x pairs
	return pairStr == "btcusd", nil
}

type Query struct {
	Time_stamp int
	// requestID string
	Venue    string
	Pair     string
	Duration int
}

func (q Query) StartTime() int {
	var startTime = q.Time_stamp - q.Duration

	return startTime
}

func (sfox SFOX) FetchOHLCV(q Query) []byte {
	fmt.Println(q)
	var url = fmt.Sprintf("https://chartdata.sfox.com/candlesticks?endTime=%d&pair=%s&period=60&startTime=%d", q.Time_stamp, q.Pair, q.StartTime())

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// fmt.Printf("\tResponse from %s: \n%s\n", q.Venue, bodyBytes)
	extractDictFromSlice(bodyBytes)

	sfox.FormatOHLCV(bodyBytes)

	return bodyBytes
}

// sFOX returns an array of dictionaries
// extracts and returns the first dictionary
func extractDictFromSlice(byteArr []byte) map[string]interface{} {
	var data []map[string]interface{}
	err := json.Unmarshal(byteArr, &data)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	// Extract the first dictionary from the slice
	if len(data) > 0 {

		fmt.Println(" ===> extractDictFromSlice: ", data[0]["close_price"])
		return data[0]
		dictionary := data[0]
		fmt.Println("---- >Dictionary:  ")
		fmt.Println(dictionary)

		// Convert the dictionary to JSON
		jsonDict, err := json.Marshal(dictionary)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
		fmt.Println("JSON:")
		fmt.Println(string(jsonDict))
	} else {
		fmt.Println("No data found")
	}
	return nil
}

// takes slice of candlesticks, returns array of OHLCV structs
func (sfox SFOX) FormatOHLCV(data []byte) any {
	// fmt.Println("Hello from formattttttttttting")
	// fmt.Println(data)
	// fmt.Println(string(data))

	var dataArr []OHLCVData
	err := json.Unmarshal(data, &dataArr)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(" ====> FormatOHLCV: ", dataArr)

	// fmt.Println(dataArr[0])
	// Print key-value pairs
	for key, value := range dataArr {
		fmt.Printf("%s: %v\n", key, value)
	}

	return dataArr[0]
}

type OHLCVData struct {
	// Define struct fields corresponding to the JSON data
	// For example:
	Open   string `json:"open_price"`
	High   string `json:"high_price"`
	Low    string `json:"low_price"`
	Close  string `json:"close_price"`
	Volume string `json:"volume"`
}

// func VenueFactory(name string) Venue {
// 	return
// }
