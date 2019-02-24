package main

import (
	"encoding/csv"
	"github.com/neofelisho/twsfex-model"
	"io"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_parse(t *testing.T) {
	type args struct {
		r io.Reader
	}
	f, _ := os.Open("./MI_5MINs_20190215.csv")

	tests := []struct {
		name string
		args args
		want [][]string
	}{{
		"test csv file",
		args{f},
		[][]string{
			{"Time", "Acc. Bid Orders", "Acc. Bid Volume", "Acc. Ask Orders", "Acc. Ask Volume", "Acc. Transaction", "Acc. Trade Volume", "Acc. Trade Value (NT$M)", ""},
			{"09:00:00", "197,267", "3,646,951", "177,784", "1,069,239", "0", "0", "0", ""},
			{"13:30:00", "5,069,709", "11,612,878", "6,275,984", "7,339,558", "1,024,038", "4,523,597", "112,550", ""},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCsv(tt.args.r)
			actual, err := got.ReadAll()
			if err != nil {
				t.Error(err)
			}
			if len(actual) != len(tt.want) {
				t.Errorf("parseCsv %v records, want %v records", len(actual), len(tt.want))
			}
			if !reflect.DeepEqual(actual, tt.want) {
				t.Errorf("actual results = %v, want %v", actual, tt.want)
			}
		})
	}
}

func Test_parseOrderBook(t *testing.T) {
	type args struct {
		r *csv.Reader
	}

	s := `
"Time","Acc. Bid Orders","Acc. Bid Volume","Acc. Ask Orders","Acc. Ask Volume","Acc. Transaction","Acc. Trade Volume","Acc. Trade Value (NT$M)",
"09:00:00","197,267","3,646,951","177,784","1,069,239","0","0","0",
"13:30:00","5,069,709","11,612,878","6,275,984","7,339,558","1,024,038","4,523,597","112,550",
`
	loc, _ := time.LoadLocation("UTC")
	openTime := time.Date(1, 1, 1, 9, 0, 0, 0, loc)
	closeTime := time.Date(1, 1, 1, 13, 30, 0, 0, loc)

	tests := []struct {
		name string
		args args
		want []model.OrderBook
	}{{
		"test order book parser",
		args{csv.NewReader(strings.NewReader(s))},
		[]model.OrderBook{{
			TimeStamp:   openTime.Unix(),
			BidOrders:   197267,
			BidVolume:   3646951,
			AskOrders:   177784,
			AskVolume:   1069239,
			Transaction: 0,
			TradeVolume: 0,
			TradeValue:  0,
		}, {
			TimeStamp:   closeTime.Unix(),
			BidOrders:   5069709,
			BidVolume:   11612878,
			AskOrders:   6275984,
			AskVolume:   7339558,
			Transaction: 1024038,
			TradeVolume: 4523597,
			TradeValue:  112550,
		}},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseOrderBook(tt.args.r)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOrderBook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNumbers(t *testing.T) {
	type args struct {
		ns string
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{{
		"number string",
		args{"123456789"},
		123456789,
	}, {
		"max uint",
		args{"18446744073709551615"},
		math.MaxUint64,
	}, {
		"negative number",
		args{"-2"},
		0,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNumbers(tt.args.ns); got != tt.want {
				t.Errorf("getNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTime(t *testing.T) {
	type args struct {
		ts string
	}
	loc, _ := time.LoadLocation("UTC")
	tests := []struct {
		name string
		args args
		want time.Time
	}{{
		name: "12:34:56",
		args: args{"12:34:56"},
		want: time.Date(1, 1, 1, 12, 34, 56, 0, loc),
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTime(tt.args.ts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
