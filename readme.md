# ARG Portfolio Tracker

This little project has the purpose of serving as a portfolio tracker for Argentinean investors.

The problem that this tool is trying to solve is the measurement of the results in USD dollars and not in Argentinean pesos, thus providing a reliable way of tracking your hefty profits.

Apart from the total and particular result for each transaction, compound annual growth return (CAGR) and weighted result (having the whole portfolio in mind) are provided.

## CSV Format

Users are supposed to give the portfolio transactions in a CSV file with the following format:

```csv
date,ticker,unit_cost,quantity,dollar_quote
```

So for example, an entry would look like this:

- `dollar_quote` is CCL (Contado con Liqui)
- `unit_cost` is in ARS
- Date format: yyyy-mm-dd

```csv
2022-05-26,spy,4320.11,9,214.08
```

## Report Table Format

```csv
____________________________________
 		RESULTS TABLE
____________________________________
 date	ticker	yield	w_yield	cagr
____________________________________
```

```csv
-----------------------------
TOTAL COST: 	$999.99
TOTAL TODAY: 	$9999.99
TOTAL RESULT:   99.99
TOTAL CAGR: 	99.99
-----------------------------
```

## How to use it

Make shure that you loaded the data in the `test.csv` file and then run the following command:

```bash
go run main.go
```

## Contribution

Just open a PR.

## TODO

- [x] Ticker not found error handling
- [x] Support CEDEARs
- [ ] Support ADRs
- [ ] Support Panel General
- [ ] Support "interactive" mode where program keeps recharging every x seconds

## API URL

### HTTP 1.0

`https://query1.finance.yahoo.com/v7/finance/spark?includePrePost=false&includeTimestamps=true&symbols=aapl&interval=1d&range=1d`

### HTTP 1.1

`https://query2.finance.yahoo.com/v7/finance/spark?includePrePost=false&includeTimestamps=true&symbols=aapl&interval=1d&range=1d`
