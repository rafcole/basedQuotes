package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	connectToServer()

}

func connectToServer() {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	// Seems less than ideal for reducing latency
	godotenv.Load()
	var accessStr = os.Getenv("MONGO_CONNECTION_STR")

	fmt.Println("access str => ", accessStr)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(accessStr).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
