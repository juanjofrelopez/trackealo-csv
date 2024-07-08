package spark

type Spark struct {
	Result []Result `json:"result"`
	Error  *string  `json:"error"`
}

type Result struct {
	Symbol   string     `json:"symbol"`
	Response []Response `json:"response"`
}

type Response struct {
	Meta Meta `json:"meta"`
	//Timestamp  []int64    `json:"timestamp"`
	//Indicators Indicators `json:"indicators"`
}

type Meta struct {
	// Currency string `json:"currency"`
	// Symbol   string `json:"symbol"`
	// ExchangeName         string               `json:"exchangeName"`
	// FullExchangeName     string               `json:"fullExchangeName"`
	// InstrumentType       string               `json:"instrumentType"`
	// FirstTradeDate       int64                `json:"firstTradeDate"`
	// RegularMarketTime    int64                `json:"regularMarketTime"`
	// HasPrePostMarketData bool                 `json:"hasPrePostMarketData"`
	// GmtOffset            int                  `json:"gmtoffset"`
	// Timezone             string               `json:"timezone"`
	// ExchangeTimezoneName string               `json:"exchangeTimezoneName"`
	RegularMarketPrice float64 `json:"regularMarketPrice"`
	// FiftyTwoWeekHigh     float64              `json:"fiftyTwoWeekHigh"`
	// FiftyTwoWeekLow      float64              `json:"fiftyTwoWeekLow"`
	// RegularMarketDayHigh float64              `json:"regularMarketDayHigh"`
	// RegularMarketDayLow  float64              `json:"regularMarketDayLow"`
	// RegularMarketVolume  int64                `json:"regularMarketVolume"`
	// ChartPreviousClose   float64              `json:"chartPreviousClose"`
	// PriceHint            int                  `json:"priceHint"`
	// CurrentTradingPeriod CurrentTradingPeriod `json:"currentTradingPeriod"`
	// DataGranularity      string               `json:"dataGranularity"`
	// Range                string               `json:"range"`
	// ValidRanges          []string             `json:"validRanges"`
}

type CurrentTradingPeriod struct {
	Pre     TradingPeriod `json:"pre"`
	Regular TradingPeriod `json:"regular"`
	Post    TradingPeriod `json:"post"`
}

type TradingPeriod struct {
	Timezone  string `json:"timezone"`
	End       int64  `json:"end"`
	Start     int64  `json:"start"`
	GmtOffset int    `json:"gmtoffset"`
}

type Indicators struct {
	Quote []Quote `json:"quote"`
}

type Quote struct {
	Close []float64 `json:"close"`
}

type Root struct {
	Spark Spark `json:"spark"`
}
