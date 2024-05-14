// // this package is non-functional, it's stubbed out to mimic how the implementation would be different for different API's

// package deribit

// import (
// 	"cryptoSnapShot/adapters"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// )

// type Deribit struct {
// 	Query adapters.Query
// }

// func (deribit Deribit) Authenticate() (int, error) {
// 	return 1, nil
// }

// func (deribit Deribit) ValidatePair(pairStr string) (bool, error) {
// 	// deribit doesn't list btc/usd directly, probably have to use something like btc-perpetual?
// 	return pairStr == "btcusd", nil
// }

// func (deribit Deribit) FetchOHLCV(q adapters.Query) adapters.OHLCVData {
// 	fmt.Println(q)
// 	var url = fmt.Sprintf("https://test.deribit.com/api/v2/public/ticker?instrument_name=BTC-PERPETUAL")

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer resp.Body.Close()

// 	ohlcv := deribit.FormatOHLCV(resp.Body)

// 	// return ohlcv
// 	return ohlcv
// }

// // TODO - change for interrogating deribit return
// func (deribit Deribit) FormatOHLCV(response io.ReadCloser) adapters.OHLCVData {
// 	var dataArr []adapters.OHLCVData

// 	err := json.NewDecoder(response).Decode(&dataArr)
// 	if err != nil {
// 		fmt.Println("Error in FormatOHLCV")
// 		log.Fatalln(err)
// 	}

// 	fmt.Println(" ====> FormatOHLCV candlestick: ", dataArr[0])

// 	candlestick := dataArr[0]

// 	return candlestick
// }
