package adapters

import (
	"io"
)

// interface
type Venue interface {
	Authenticate() (int, error)
	ValidatePair(string) (bool, error)
	FetchOHLCV(Query) io.ReadCloser // Timestamp within method or at time of request?
	FormatOHLCV(io.ReadCloser) any  // Last touch point
}

type Query struct {
	Time_stamp int
	Venue      string
	Pair       string
	Duration   int
}

func (q Query) StartTime() int {
	var startTime = q.Time_stamp - q.Duration

	return startTime
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
