package main

// import (
// 	"google.golang.org/grpc"
// )

// type ExchangeRates struct {
// 	rates map[string]float64
// }

// func NewExchangeRates() *ExchangeRates {
// 	return &ExchangeRates{
// 		rates: make(map[string]float64),
// 	}
// }

// func (e *ExchangeRates) SetRate(fromCurrency, toCurrency string, rate float64) {
// 	key := fromCurrency + toCurrency
// 	e.rates[key] = rate

// 	reciprocalKey := toCurrency + fromCurrency
// 	reciprocalRate := 1.0 / rate
// 	e.rates[reciprocalKey] = reciprocalRate
// }

// func (e *ExchangeRates) GetRate(fromCurrency, toCurrency string) (float64, bool) {
// 	key := fromCurrency + toCurrency
// 	rate, ok := e.rates[key]
// 	return rate, ok
// }
