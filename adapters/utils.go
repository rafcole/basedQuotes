package adapters

import (
	"io"
)

// interface
type Venue interface {
	Authenticate() (int, error)
	ValidatePair(string) (bool, error)
	FetchOHLCV(Query) OHLCVData          // Timestamp within method or at time of request?
	FormatOHLCV(io.ReadCloser) OHLCVData // Last touch point
}

type Query struct {
	Time_stamp     int
	Venue          string
	Currency_Base  string
	Currency_Quote string
	Duration       int
	Request_ID     string
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
