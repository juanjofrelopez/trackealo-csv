package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/juanjofrelopez/csv-cli-tool/colors"
	"github.com/juanjofrelopez/csv-cli-tool/ratios"
	"github.com/juanjofrelopez/csv-cli-tool/spark"
)

const MAX_ENTRIES = 30
const layoutISO = "2006-01-02"
const hoursInAYear = 8760

var quotes map[string]float64 = make(map[string]float64)
var tickers string

type Entry struct {
	date        time.Time
	ticker      string
	price       float64
	quantity    uint64
	dollarQuote float64
}

type ResultsEntry struct {
	date        time.Time
	ticker      string
	cost        float64
	quantity    int
	dollarQuote float64
	result      float64
	wResult     float64
	cagr        float64
}

type ResultsTable struct {
	entries     []ResultsEntry
	totalToday  float64
	totalCost   float64
	totalResult float64
	years       float64
}

// regularMarketPrice float64
func ParseRecord(record []string) Entry {
	entry := Entry{}
	for i, v := range record {
		switch i {
		case 0:
			parsed, err := time.Parse(layoutISO, v)
			if err != nil {
				log.Fatal(err)
			}
			entry.date = parsed
		case 1:
			entry.ticker = v
		case 2:
			val, err := strconv.ParseFloat(v, 32)
			if err != nil {
				log.Fatal(err)
			}
			entry.price = val
		case 3:
			val, err := strconv.ParseUint(v, 10, 16)
			if err != nil {
				log.Fatal(err)
			}
			entry.quantity = val
		case 4:
			val, err := strconv.ParseFloat(v, 32)
			if err != nil {
				log.Fatal(err)
			}
			entry.dollarQuote = val
		default:
			log.Fatal("error parsing csv")
		}
	}
	return entry
}

func getUSPrice() float64 {
	res, err := http.Get("https://query2.finance.yahoo.com/v7/finance/spark?includePrePost=false&includeTimestamps=true&symbols=aapl&interval=1d&range=1d")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var data spark.Root
	json.NewDecoder(res.Body).Decode(&data)
	usPrice := data.Spark.Result[0].Response[0].Meta.RegularMarketPrice
	return usPrice
}

func getARSPrice() float64 {
	res, err := http.Get("https://query2.finance.yahoo.com/v7/finance/spark?includePrePost=false&includeTimestamps=true&symbols=aapl.ba&interval=1d&range=1d")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	var data spark.Root
	json.NewDecoder(res.Body).Decode(&data)
	arsPrice := data.Spark.Result[0].Response[0].Meta.RegularMarketPrice
	return arsPrice
}

func CalculateDollarQuote() float64 {
	// aapl US Price
	usPrice := getUSPrice()
	// const aaplARGPrice
	arsPrice := getARSPrice()
	ratio := ratios.GetTickerRatio("aapl")
	return (arsPrice * ratio) / usPrice
}

func printResultTable(t ResultsTable) {
	fmt.Println("_________________________________________________")
	fmt.Println(colors.BoldYellow, "\t\tRESULTS TABLE", colors.Reset)
	fmt.Println("_________________________________________________")
	fmt.Println(colors.BoldYellow, "date\t\tticker\tyield\tw yield\tcagr", colors.Reset)
	fmt.Println("_________________________________________________")

	for _, v := range t.entries {
		yield := v.result * 100
		wYield := v.wResult * 100
		cagr := v.cagr * 100
		var yieldColor string
		if yield > 0 {
			yieldColor = colors.BoldGreen
		} else {
			yieldColor = colors.BoldRed
		}
		fmt.Printf("%s\t%s\t%s%.2f\t%.2f\t%.2f %s\n", v.date.Format(layoutISO), v.ticker, yieldColor, yield, wYield, cagr, colors.Reset)
	}
	var yieldColor string
	if t.totalResult > 0 {
		yieldColor = colors.BoldGreen
	} else {
		yieldColor = colors.BoldRed
	}
	fmt.Println("_________________________________________________")
	fmt.Println("-----------------------------")

	fmt.Printf("%sTOTAL COST: \t$%.2f%s\n", colors.BoldYellow, t.totalCost, colors.Reset)
	fmt.Printf("%sTOTAL TODAY: \t$%.2f%s\n", colors.BoldYellow, t.totalToday, colors.Reset)
	fmt.Printf("%sTOTAL RESULT: \t%.2f%s\n", yieldColor, t.totalResult*100, colors.Reset)

	totalCagr := math.Pow((t.totalResult+1), 1/t.years) - 1

	fmt.Printf("%sTOTAL CAGR: \t%.2f%s\n", yieldColor, totalCagr*100, colors.Reset)
	fmt.Println("-----------------------------")

}

func main() {
	file, err := os.Open("test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	entries := make([]Entry, 0, MAX_ENTRIES)

	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		entry := ParseRecord(record)
		entries = append(entries, entry)
	}

	for _, entry := range entries {
		if _, err := quotes[entry.ticker]; !err {
			quotes[entry.ticker] = 0
			tickers = fmt.Sprintf("%s,%s", tickers, entry.ticker)
		}
	}

	tickers, _ = strings.CutPrefix(tickers, ",")
	url := fmt.Sprintf("https://query2.finance.yahoo.com/v7/finance/spark?includePrePost=false&includeTimestamps=true&symbols=%s&interval=1d&range=1d", tickers)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var data spark.Root
	json.NewDecoder(res.Body).Decode(&data)

	for _, v := range data.Spark.Result {
		price := v.Response[0].Meta.RegularMarketPrice
		quotes[v.Symbol] = price
	}

	// dollarQuote := CalculateDollarQuote()
	resultsTable := ResultsTable{years: 0}

	for _, v := range entries {
		boughtPriceUSD := v.price * ratios.GetTickerRatio(v.ticker) / v.dollarQuote

		var newResult float64
		if quotes[v.ticker] == 0 {
			newResult = 0
		} else {
			newResult = (quotes[v.ticker] / boughtPriceUSD) - 1
			resultsTable.totalCost += (v.price * float64(v.quantity)) / v.dollarQuote
		}

		// for merval
		// newResult := ( ( quotes[v.ticker] / todayUSDRate ) / boughtPriceUSD) - 1

		years := time.Since(v.date).Hours() / hoursInAYear
		cagr := math.Pow((newResult+1), 1/years) - 1
		resultsTable.years = math.Max(years, resultsTable.years)

		newEntry := ResultsEntry{
			date:        v.date,
			ticker:      v.ticker,
			cost:        v.price * float64(v.quantity),
			quantity:    int(v.quantity),
			dollarQuote: v.dollarQuote,
			result:      newResult,
			cagr:        cagr,
		}
		resultsTable.entries = append(resultsTable.entries, newEntry)
		todaysTotal := (quotes[v.ticker] / ratios.GetTickerRatio(v.ticker)) * float64(v.quantity)
		resultsTable.totalToday += todaysTotal
	}

	for i, v := range resultsTable.entries {
		var todayTotal float64
		if quotes[v.ticker] == 0 {
			todayTotal = 0
		} else {
			todayTotal = (quotes[v.ticker] / ratios.GetTickerRatio(v.ticker)) * float64(v.quantity)
		}
		resultsTable.entries[i].wResult = (todayTotal / resultsTable.totalToday) * v.result
	}

	resultsTable.totalResult = ((resultsTable.totalToday / resultsTable.totalCost) - 1)
	printResultTable(resultsTable)

}
