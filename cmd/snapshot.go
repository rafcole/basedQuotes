package cmd

import (
	"context"
	"cryptoSnapShot/adapters"
	"cryptoSnapShot/adapters/sfox"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func takeSnapShot(cmd *cobra.Command, args []string) {
	// default duration appears to be 60s
	const duration int = 60
	// unix_timestamp := int(time.Now().Unix())

	venue_str := args[0]
	pair := args[1]

	// fmt.Println("Venue: " + venue_str)
	// fmt.Println("Pair: " + pair)

	var query = adapters.Query{
		Time_stamp:    int(time.Now().Unix()),
		Venue:         venue_str,
		Currency_Pair: pair,
		Duration:      duration,
		Request_ID:    uuid.NewString(),
	}

	var venue = AdapterFactory(query)

	// TODO validate venue
	fmt.Println("===> Status: Validating venue")
	// TODO validate pair (can this be done ahead of time?)
	fmt.Println("===> Status: Validating crypto pair")

	// TODO Send request
	fmt.Printf("===> Status: Requesting %s OHLCV from %s at %d Unix\n", query.Currency_Pair, query.Venue, query.Time_stamp)

	data := venue.FetchOHLCV(query)
	fmt.Println("data in takeSnapshot: ", data)

	// venue.FormatOHLCV(data)

	snapshot := GenerateSnapshot(query, data)

	connectToServer(snapshot)
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

func GenerateSnapshot(q adapters.Query, d adapters.OHLCVData) Snapshot {
	return Snapshot{
		Request_ID:        q.Request_ID,
		Request_Timestamp: q.Time_stamp,
		Venue_Name:        q.Venue,
		Currency_Pair:     q.Currency_Pair,
		Market_Data:       d,
	}
}

type Snapshot struct {
	Request_ID        string
	Request_Timestamp int
	Venue_Name        string
	Currency_Pair     string
	Market_Data       adapters.OHLCVData
	// RequestID
	// Status
}

func connectToServer(data any) {
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

	coll := client.Database("crypto_thp").Collection("ohlcv_shapshots")

	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		panic(err)
	}

	fmt.Println("result from mdb: ", result)
	// Send a ping to confirm a successful connection
	// if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
