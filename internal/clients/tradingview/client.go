package tradingview

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/pkg/errors"
)

const (
	apiURL = "https://scanner.tradingview.com/forex/scan"
)

const (
	UsdTicker = "USDRUB"
	CnyTicker = "CNYRUB"
	EurTicker = "EURRUB"
)

type Client struct {
	httpClient *http.Client
}

func New() *Client {
	return &Client{
		httpClient: http.DefaultClient,
	}
}

var errCodes = []string{"400", "404", "402", "500", "502", "505"}

func (c *Client) GetQuote(ticker string) (float64, error) {
	var (
		resp *http.Response
		err  error
	)

	ticker = strings.ToUpper(ticker)
	body := []byte(`{"symbols":{"tickers":["FX_IDC:` + ticker + `"],"query":{"types":[]}},"columns":["close", "change_abs", "change"]}`)

	resp, err = http.Post(apiURL, "multipart/form-data", bytes.NewReader(body))
	if err != nil || checkResStatus(resp.Status) {
		return 0, errors.New("invalid response status code")
	}

	defer resp.Body.Close()

	htmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("unable to decode response body: %w", err)
	}

	val, _, _, err := jsonparser.Get(htmlData, "data", "[0]", "d")
	if err != nil {
		return 0, err
	}

	str := strings.Replace(string(val), "[", "", -1)
	str1 := strings.Replace(str, "]", "", -1)
	strArr := strings.Split(str1, ",")

	quoteVal := float64(0)
	quoteVal, err = strconv.ParseFloat(strArr[0], 64)
	if err != nil {
		return 0, err
	}

	return quoteVal, nil
}

func checkResStatus(s string) bool {
	for _, code := range errCodes {
		if code == s {
			return true
		}
	}
	return false
}
