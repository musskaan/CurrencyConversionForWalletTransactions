package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	pb "conversion.com/currency-conversion/conversion"
	"google.golang.org/grpc"
)

const (
	port              = ":50051"
	exchangeRatesFile = "exchange_rates.json"
)

type ExchangeRates struct {
	rates map[string]float64
}

func NewExchangeRates() *ExchangeRates {
	return &ExchangeRates{
		rates: make(map[string]float64),
	}
}

func (e *ExchangeRates) SetRate(fromCurrency, toCurrency string, rate float64) {
	key := fromCurrency + toCurrency
	e.rates[key] = rate

	reciprocalKey := toCurrency + fromCurrency
	reciprocalRate := 1.0 / rate
	e.rates[reciprocalKey] = reciprocalRate
}

func (e *ExchangeRates) GetRate(fromCurrency, toCurrency string) (float64, bool) {
	key := fromCurrency + toCurrency
	rate, ok := e.rates[key]
	return rate, ok
}

type CurrencyConversionServer struct {
	pb.UnimplementedConversionServiceServer
	exchangeRates *ExchangeRates
}

func convertToBaseCurrency(transferAmount float64, sourceCurrency, baseCurrency string, exchangeRates *ExchangeRates) (float64, error) {
	rate, ok := exchangeRates.GetRate(sourceCurrency, baseCurrency)
	if !ok {
		return 0, fmt.Errorf("currency not supported: %s to %s", sourceCurrency, baseCurrency)
	}

	convertedAmount := transferAmount * rate
	return convertedAmount, nil
}

func (s *CurrencyConversionServer) ConvertCurrency(ctx context.Context, request *pb.ConversionRequest) (*pb.ConversionResponse, error) {
	baseCurrency := request.GetBaseCurrency()
	sourceCurrency := request.GetSourceCurrency()
	trasferAmount := request.GetTransferAmount()

	convertedAmount, err := convertToBaseCurrency(trasferAmount, sourceCurrency, baseCurrency, s.exchangeRates)
	if err != nil {
		return nil, err
	}

	response := &pb.ConversionResponse{
		ConvertedAmount: convertedAmount,
	}

	return response, nil
}

func (e *ExchangeRates) LoadFromFile() error {
	file, err := os.Open(exchangeRatesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var rates map[string]float64
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&rates)
	if err != nil {
		return err
	}
	e.rates = rates
	return nil
}

func (e *ExchangeRates) SaveToFile() error {
	file, err := os.Create(exchangeRatesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(e.rates)
	return err
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	exchangeRates := NewExchangeRates()
	if err := exchangeRates.LoadFromFile(); err != nil {
		log.Printf("Error loading exchange rates from file: %v", err)
	}

	// exchangeRates.SetRate("USD", "CAD", 1.35)

	pb.RegisterConversionServiceServer(s, &CurrencyConversionServer{exchangeRates: exchangeRates})

	log.Printf("Go GRPC server started on port %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	if err := exchangeRates.SaveToFile(); err != nil {
		log.Printf("Error saving exchange rates to file: %v", err)
	}
}
