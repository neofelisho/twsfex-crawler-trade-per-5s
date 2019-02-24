package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/neofelisho/twsfex-model"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	csvUrl      string
	dateString  string
	date        time.Time
	mongoApiUrl string
)

func main() {
	initFlags()

	ioReader := getCsvDataFromUrl()
	//ioReader := inputFromFile()
	defer func() {
		if err := ioReader.Close(); err != nil {
			panic(err)
		}
	}()

	r := parseCsv(ioReader)
	orderBooks := parseOrderBook(r)
	saveToDb(orderBooks)
}

func getCsvDataFromUrl() io.ReadCloser {
	s := csvUrl + dateString
	resp, err := http.Get(s)
	if err != nil {
		panic(err)
	}

	return resp.Body
}

func inputFromFile() *os.File {
	f, err := os.Open("MI_5MINs_20190215.csv")
	if err != nil {
		panic(err)
	}
	return f
}

func saveToDb(orderBooks []model.OrderBook) {
	doc := model.TWSEDaily{
		Date:       date,
		Opening:    orderBooks[0],
		Close:      orderBooks[1],
		UpdateTime: time.Now(),
	}
	jsonValue, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}

	apiUrl, err := url.Parse(mongoApiUrl)
	if err != nil {
		panic(err)
	}
	response, err := http.Post(apiUrl.String(), "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		panic(err)
	}

	fmt.Println(response.Status)
}

func parseOrderBook(r *csv.Reader) []model.OrderBook {
	ss, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	if len(ss) < 3 {
		panic("can not parse daily trade per 5s record from Taiwan stock exchange")
	}
	ss = ss[1:] //ignore the header
	results := make([]model.OrderBook, len(ss))
	for i := 0; i < len(ss); i++ {
		row := ss[i]
		results[i] = model.OrderBook{
			TimeStamp:   getTime(row[0]).Unix(),
			BidOrders:   getNumbers(strings.Replace(row[1], ",", "", -1)),
			BidVolume:   getNumbers(strings.Replace(row[2], ",", "", -1)),
			AskOrders:   getNumbers(strings.Replace(row[3], ",", "", -1)),
			AskVolume:   getNumbers(strings.Replace(row[4], ",", "", -1)),
			Transaction: getNumbers(strings.Replace(row[5], ",", "", -1)),
			TradeVolume: getNumbers(strings.Replace(row[6], ",", "", -1)),
			TradeValue:  getNumbers(strings.Replace(row[7], ",", "", -1)),
		}
	}
	return results
}

func getNumbers(ns string) uint64 {
	numbers, _ := strconv.ParseUint(ns, 10, 64)
	return numbers
}

func getTime(ts string) time.Time {
	hh, _ := strconv.ParseInt(ts[:2], 10, 8)
	mm, _ := strconv.ParseInt(ts[3:5], 10, 8)
	ss, _ := strconv.ParseInt(ts[6:], 10, 8)
	t := date.Add(time.Hour*time.Duration(hh) + time.Minute*time.Duration(mm) + time.Second*time.Duration(ss))
	return t
}

func parseCsv(r io.Reader) *csv.Reader {
	s := ""
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		t = strings.Replace(t, "=", "", -1)
		if strings.HasPrefix(t, "\"Time\"") || strings.HasPrefix(t, "\"09:00:00\"") || strings.HasPrefix(t, "\"13:30:00\"") {
			s += t + "\n"
		}
	}

	return csv.NewReader(strings.NewReader(s))
}

func initFlags() {
	flag.StringVar(&csvUrl, "csvUrl", "http://www.twse.com.tw/en/exchangeReport/MI_5MINS?response=csv&date=", "the url to get csv file")
	flag.StringVar(&dateString, "date", "20190218", "which day's data should be parsed")
	flag.StringVar(&mongoApiUrl, "mongoApiUrl", "http://localhost:8080/daily", "backend api for accessing database")
	flag.Parse()

	getDate()
}

func getDate() {
	local, _ := time.LoadLocation("Asia/Taipei")
	year, _ := strconv.ParseInt(dateString[:4], 10, 16)
	month, _ := strconv.ParseInt(dateString[4:6], 10, 8)
	day, _ := strconv.ParseInt(dateString[6:], 10, 8)
	date = time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, local)
}
