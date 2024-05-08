package cmd

import (
	"context"
	"cryptoSnapShot/adapters"
	"cryptoSnapShot/adapters/sfox"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func parseCurrencyPair(pairStr string) (string, string) {
	// TODO error handling for improperly formatted args
	parts := strings.Split(pairStr, "/")
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func takeSnapShot(cmd *cobra.Command, args []string) {
	venue_str := args[0]
	pair := args[1]

	// This internal query object requires the most modification between API adapters
	// IE this one was setup for sFox candlesticks, which require a duration arguement
	// looking deeper into the deribiti and cmc api's I see this is not universal

	base, quote := parseCurrencyPair(pair)

	var query = adapters.Query{
		Time_stamp:     int(time.Now().Unix()),
		Venue:          venue_str,
		Currency_Base:  base,
		Currency_Quote: quote,
		Duration:       60,
		Request_ID:     uuid.NewString(),
	}

	// TODO validate venue
	fmt.Println("===> Status: Validating venue")

	venue := AdapterFactory(query)

	// TODO ensure API supports the entered currency pair arg
	// venue.ValidatePair(pair)
	fmt.Println("===> Status: Validating crypto pair")
	// TODO parse base/quote from arg
	// IE [base]/[quote]
	venue.ValidatePair()

	// TODO Send request
	fmt.Printf("===> Status: Requesting %s OHLCV from %s at %d Unix\n", venue.FormattedCurrencyPair(), query.Venue, query.Time_stamp)

	data := venue.FetchOHLCV(query)
	fmt.Println("data in takeSnapshot: ", data)

	// venue.FormatOHLCV(data)

	snapshot := GenerateSnapshot(query, data)

	connectToServer(snapshot)
}

func AdapterFactory(q adapters.Query) adapters.Venue {
	switch q.Venue {
	case "sfox":
		return sfox.SFOX{Query: q}
	case "deribit":
		// return deribit.Deribit{Query: q}
		log.Fatalln("Deribit implementation incomplete. Only s. Please see ./adapters/_deribit/_deribit.go for function signatures/stubs")
		return nil
	case "coinmarketcap":
		// return coinmarketcap.CMC{Query: q}
		log.Fatalln("CoinMarketCap implementation incomplete. Please see ./adapters/ for function signatures/stubs of how this would be implemented")
		return nil
	default:
		log.Fatalln("Unsupported vendor")
		return nil
	}
}

func GenerateSnapshot(q adapters.Query, d adapters.OHLCVData) Snapshot {
	return Snapshot{
		Request_ID:        q.Request_ID,
		Request_Timestamp: q.Time_stamp,
		Venue_Name:        q.Venue,
		Currency_Base:     q.Currency_Base,
		Currency_Quote:    q.Currency_Quote,
		Market_Data:       d,
	}
}

type Snapshot struct {
	Request_ID        string
	Request_Timestamp int
	Venue_Name        string
	Currency_Base     string
	Currency_Quote    string
	Market_Data       adapters.OHLCVData
}

func connectToServer(data any) {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	// Seems less than ideal for reducing latency
	godotenv.Load()
	var accessStr = os.Getenv("MONGO_CONNECTION_STR")

	fmt.Println("AccessStr = ", accessStr)

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

	coll := client.Database(os.Getenv("MONGO_DATABASE_NAME")).Collection(os.Getenv("MONGO_COLLECTION_NAME"))

	result, err := coll.InsertOne(context.TODO(), data)
	if err != nil {
		panic(err)
	}

	fmt.Println("result from mdb: ", result)
}
