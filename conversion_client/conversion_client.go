package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	pb "conversion.com/currency-conversion/conversion"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcServerAddr = "localhost:50051"
	port           = ":8081"
)

func Convert(client pb.ConversionServiceClient, w http.ResponseWriter, r *http.Request) {
	var req pb.ConversionRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	request := &pb.ConversionRequest{
		BaseCurrency:   req.BaseCurrency,
		SourceCurrency: req.SourceCurrency,
		TransferAmount: req.TransferAmount,
	}

	response, err := client.ConvertCurrency(context.Background(), request)
	if err != nil {
		log.Fatal("failed to convert currency: %v", err)
	}

	res := pb.ConversionResponse{ConvertedAmount: response.ConvertedAmount}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	conn, err := grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewConversionServiceClient(conn)

	router := mux.NewRouter()

	router.HandleFunc("/convert", func(writer http.ResponseWriter, req *http.Request) {
		Convert(client, writer, req)
	}).Methods("POST")

	log.Println("HTTP Server listening on", port)
	http.ListenAndServe(port, router)
}
