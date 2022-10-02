package datahandler

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var spendSlice = make([]SpendData, 0)

type Handler struct {
}

type SpendData struct {
	summa    int
	category string
	dateTime time.Time
}

func New() (*Handler, error) {
	return &Handler{}, nil
}

func AddSpend(data string) error {
	strSlice := strings.Split(data, "/")
	var spendData SpendData
	var summa int
	var dateSpend time.Time

	summa, errSumm := strconv.Atoi(strSlice[0])
	if errSumm != nil {
		return errors.Wrap(errSumm, "")
	}
	spendData.summa = summa

	spendData.category = strSlice[1]

	dateSpend, errDateSpend := time.Parse("2006-01-02", strSlice[2])
	if errDateSpend != nil {
		return errors.Wrap(errDateSpend, "")
	}
	spendData.dateTime = dateSpend
	fmt.Print(strSlice[2])
	fmt.Print(dateSpend)

	spendSlice = append(spendSlice, spendData)

	return errors.Wrap(nil, "")
}

func Report(code int) string {
	var str string
	var categoryMap = make(map[string]int, 0)
	currentDate := time.Now()
	var nextSDate time.Time
	switch code {
	case 1:
		nextSDate = currentDate.AddDate(0, 0, -7)
	case 2:
		nextSDate = currentDate.AddDate(0, -1, 0)
	case 3:
		nextSDate = currentDate.AddDate(-1, 0, 0)
	}

	sort.SliceStable(spendSlice, func(i, j int) bool {
		return spendSlice[i].dateTime.After(spendSlice[j].dateTime)
	})

	for _, v := range spendSlice {
		if v.dateTime.After(nextSDate) {
			categoryMap[v.category] += v.summa
		} else if v.dateTime.Before(nextSDate) {
			break
		}
	}

	keys := make([]string, 0, len(categoryMap))
	for k := range categoryMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		str = str + k + " - " + strconv.Itoa(categoryMap[k]) + "\n"
	}

	return str
}
