package main

import (
	"context"
	"testing"

	pb "conversion.com/currency-conversion/conversion"
	"github.com/stretchr/testify/assert"
)

func TestCurrencyConversionServer_SuccessfulConversionFromUSDtoINR(t *testing.T) {
	exchangeRates := NewExchangeRates()
	exchangeRates.SetRate("USD", "INR", 82.87)

	server := &CurrencyConversionServer{exchangeRates: exchangeRates}

	request := &pb.ConversionRequest{
		BaseCurrency:   "INR",
		SourceCurrency: "USD",
		TransferAmount: 100,
	}

	expectedResponse := &pb.ConversionResponse{
		ConvertedAmount: 8287.0,
	}

	response, err := server.ConvertCurrency(context.Background(), request)

	assert.Nil(t, err, "Unexpected error: %v", err)
	assert.NotNil(t, response)
	assert.Equal(t, response.ConvertedAmount, expectedResponse.ConvertedAmount, "Expected converted amount %f, but got %f", expectedResponse.ConvertedAmount, response.ConvertedAmount)
}

func TestCurrencyConversionServer_UnsupportedCurrency(t *testing.T) {
	exchangeRates := NewExchangeRates()

	server := &CurrencyConversionServer{exchangeRates: exchangeRates}

	request := &pb.ConversionRequest{
		BaseCurrency:   "INR",
		SourceCurrency: "EUR",
		TransferAmount: 100,
	}

	expectedErrorMsg := "currency not supported: EUR to INR"

	response, err := server.ConvertCurrency(context.Background(), request)

	assert.NotNil(t, err, "Expected an error, but got none.")
	assert.Equal(t, err.Error(), expectedErrorMsg, "Expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
	assert.Nil(t, response, "Expected nil response, but got %+v", response)
}
