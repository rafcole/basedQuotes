package sfox

import (
	"cryptoSnapShot/adapters"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type SFOX struct {
	Query adapters.Query
}

func (sfox SFOX) Authenticate() (int, error) {
	fmt.Println("sFOX candlestick API does not require authentication")
	return 1, nil
}

func (sfox SFOX) ValidatePair(pairStr string) (bool, error) {
	// avaliable pairs accessible through GET https://api.sfox.com/v1/currency-pairs
	// ideally we'd have that data pulled into a file ahead of time
	// for the sake of time I'll support 3x pairs
	// TODO dict of 3x pairs
	return pairStr == "btcusd", nil
}

func (sfox SFOX) FetchOHLCV(q adapters.Query) io.ReadCloser {
	fmt.Println(q)
	var url = fmt.Sprintf("https://chartdata.sfox.com/candlesticks?endTime=%d&pair=%s&period=60&startTime=%d", q.Time_stamp, q.Pair, q.StartTime())

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	ohlcv := sfox.FormatOHLCV(resp.Body)

	fmt.Println(ohlcv)

	// return ohlcv
	return nil
}

// takes slice of candlesticks, returns array of OHLCV structs
func (sfox SFOX) FormatOHLCV(response io.ReadCloser) any {
	var dataArr []adapters.OHLCVData

	err := json.NewDecoder(response).Decode(&dataArr)
	if err != nil {
		fmt.Println("Error in FormatOHLCV")
		log.Fatalln(err)
	}

	fmt.Println(" ====> FormatOHLCV candlestick: ", dataArr[0])

	// This is a hardcode to extract the only element in the array
	// In a more sophisticated program it may be necessary to process
	// the response as a stream as detailed as detailed in encoder pkg
	// docs https://pkg.go.dev/encoding/json@go1.22.2#NewDecoder:~:text=an%20input%20stream.-,Example,-%C2%B6
	candlestick := dataArr[0]

	return candlestick
}

// sFOX returns an array of dictionaries
// extracts and returns the first dictionary
// func extractDictFromSlice(byteArr []byte) map[string]interface{} {
// 	var data []map[string]interface{}
// 	err := json.Unmarshal(byteArr, &data)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return nil
// 	}

// 	// Extract the first dictionary from the slice
// 	if len(data) > 0 {

// 		fmt.Println(" ===> extractDictFromSlice: ", data[0]["close_price"])
// 		return data[0]
// 		dictionary := data[0]
// 		fmt.Println("---- >Dictionary:  ")
// 		fmt.Println(dictionary)

// 		// Convert the dictionary to JSON
// 		jsonDict, err := json.Marshal(dictionary)
// 		if err != nil {
// 			fmt.Println("Error:", err)
// 			return nil
// 		}
// 		fmt.Println("JSON:")
// 		fmt.Println(string(jsonDict))
// 	} else {
// 		fmt.Println("No data found")
// 	}
// 	return nil
// }
