package cmd

import (
	"context"
	"cryptoSnapShot/adapters"
	"cryptoSnapShot/adapters/sfox"
	"fmt"
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

	var venue_str string = args[0]
	var pair string = args[1]

	fmt.Println("Venue: " + venue_str)
	fmt.Println("Pair: " + pair)

	var query = adapters.Query{
		Time_stamp: unix_timestamp,
		Venue:      venue_str,
		Pair:       pair,
		Duration:   60,
	}

	var venue = AdapterFactory(query)

	// TODO validate venue
	fmt.Println("===> Status: Validating venue")
	// TODO validate pair (can this be done ahead of time?)
	fmt.Println("===> Status: Validating crypto pair")

	// TODO Send request
	fmt.Printf("===> Status: Requesting %s OHLCV from %s at %d Unix\n", query.Pair, query.Venue, query.Time_stamp)

	data := venue.FetchOHLCV(query)
	fmt.Println("data in takeSnapshot: ", data)

	// venue.FormatOHLCV(data)

	// connectToServer()

}

func AdapterFactory(q adapters.Query) adapters.Venue {
	switch q.Venue {
	case "sfox":
		// validation here?
		fmt.Println("Found sfox in AdapterFactory")
		return sfox.SFOX{Query: q}
	default:
		return nil
	}
}

type Snapshot struct {
	Time   int
	Open   float64
	High   float64
	Close  float64
	Volume float64
	Venue  string
	Pair   string
	// RequestID
	// Status
}

func connectToServer() {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	// Seems less than ideal for reducing latency
	godotenv.Load()
	var accessStr = os.Getenv("MONGO_CONNECTION_STR")

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
