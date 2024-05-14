package snapshot

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
	parts := strings.Split(pairStr, "/")
	if len(parts) != 2 {
		log.Fatalln("Improperly formatted currency pair. Please enter \"[base]/[quote]\"")
		return "", ""
	}
	return parts[0], parts[1]
}

func TakeSnapShot(cmd *cobra.Command, args []string) {
	venue_str := args[0]
	pair := args[1]
	base, quote := parseCurrencyPair(pair)

	query := adapters.Query{
		Time_stamp:     int(time.Now().Unix()),
		Venue:          venue_str,
		Currency_Base:  base,
		Currency_Quote: quote,
		Duration:       60,
		Request_ID:     uuid.NewString(),
	}

	printStatus("Validating venue")
	venue, err := AdapterFactory(query)
	if err != nil {
		log.Fatalln(err)
	}

	printStatus("Validating crypto pair")
	err = venue.ValidatePair()
	if err != nil {
		log.Fatalln(err)
	}

	printStatus(fmt.Sprintf("Requesting %s OHLCV from %s at %d Unix\n", venue.FormattedCurrencyPair(), query.Venue, query.Time_stamp))
	data, err := venue.FetchOHLCV(query)
	if err != nil {
		log.Fatalln(err)
	}
	// fmt.Println("\ndata in takeSnapshot: ", data)

	printStatus("Generating snapshot")
	snapshot := GenerateSnapshot(query, data)

	printStatus("Sending snapshot to database")
	serverResponse, err := pushSnapshotToDB(snapshot)
	if err != nil {
		log.Fatalln(err)
	}
	printStatus(fmt.Sprintf("Snapshot succesfully sent to database: %v", serverResponse))
}

func printStatus(str string) {
	fmt.Println("\n===> Status: " + str)
}

func AdapterFactory(q adapters.Query) (adapters.Venue, error) {
	switch q.Venue {
	case "sfox":
		return sfox.SFOX{Query: q}, nil
	case "deribit":
		// return deribit.Deribit{Query: q}
		log.Fatalln("Deribit implementation incomplete. Only s. Please see ./adapters/_deribit/_deribit.go for function signatures/stubs")
		return nil, nil
	case "coinmarketcap":
		// return coinmarketcap.CMC{Query: q}
		log.Fatalln("CoinMarketCap implementation incomplete. Please see ./adapters/ for function signatures/stubs of how this would be implemented")
		return nil, nil
	default:
		log.Fatalln("Unsupported vendor")

		return nil, nil
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

func pushSnapshotToDB(data any) (interface{}, error) {
	godotenv.Load()
	accessStr := os.Getenv("MONGO_CONNECTION_STR")

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

	return result.InsertedID, nil
}
