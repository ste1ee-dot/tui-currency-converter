package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/charmbracelet/huh"
)



func ShowOptions(baseCurrency *string, newCurrency *string, howMuch *int) {
	var (
		baseCur string
		newCur string
		howM string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
			Title("Choose base currency").
			Options(
				huh.NewOption("EUR", "EUR"),
				huh.NewOption("USD", "USD"),
				huh.NewOption("GBP", "GBP"),
				huh.NewOption("RSD", "RSD"),
			).
			Value(&baseCur),

			huh.NewSelect[string]().
			Title("Choose currency to convert to").
			Options(
				huh.NewOption("EUR", "EUR"),
				huh.NewOption("USD", "USD"),
				huh.NewOption("GBP", "GBP"),
				huh.NewOption("RSD", "RSD"),
			).
			Value(&newCur),

			huh.NewInput().
			Title("How much to convert").
			Value(&howM),

		),
	)

	err := form.Run()
	if err != nil { panic(err) }

	howMuchInt, err := strconv.Atoi(howM)
	if err != nil { panic(err) }

	*baseCurrency = baseCur
	*newCurrency = newCur
	*howMuch = howMuchInt
}

type Response struct {
	Disclaimer string
	License string
	Timestamp int
	Base string
	Rates map[string]float64 `json:"rates"`
}

func main() {

	var (
		baseCurrency string
		newCurrency string
		howMuch int
	)

	ShowOptions(&baseCurrency, &newCurrency, &howMuch)

	resp, err :=	http.Get("https://openexchangerates.org/api/latest.json" + 
	"?app_id=" + apiKey + "&symbols=" + 
	baseCurrency + "," + newCurrency)
	if err != nil { panic(err) }
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil { panic(err) }

	var response Response
	json.Unmarshal(body, &response)

	baseToUSD := response.Rates[baseCurrency]
	newToUSD := response.Rates[newCurrency]

	rate := float64(howMuch) * (newToUSD / baseToUSD)

	fmt.Println(rate)
}
