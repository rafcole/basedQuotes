package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

func takeSnapShot(cmd *cobra.Command, args []string) {
	const duration int = 60
	var unix_timestamp = int(time.Now().Unix())
	var startTime = unix_timestamp - duration

	var venue = args[0]
	var pair = args[1]

	fmt.Println("Venue: " + venue)
	fmt.Println("Pair: " + pair)

	// TODO validate venue
	fmt.Println("Status: Validating venue")
	// TODO validate pair (can this be done ahead of time?)
	fmt.Println("Status: Validating crypto pair")

	// TODO Send request
	fmt.Printf("Status: Requesting %s OHLCV from %s at %d Unix\n", pair, venue, unix_timestamp)

	var url = fmt.Sprintf("https://chartdata.sfox.com/candlesticks?endTime=%d&pair=%s&period=60&startTime=%d", unix_timestamp, pair, startTime)

	// var url = "https://chartdata.sfox.com/candlesticks?endTime=1665165809&pair=btcusd&period=86400&startTime=1657217002"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	fmt.Printf("\tResponse from %s: \n%s\n", venue, string(responseBody))

}
