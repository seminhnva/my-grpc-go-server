package main

import (
	"database/sql"
	"log"
	"math/rand/v2"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	dbmigration "github.com/seminhnva/my-grpc-go-server/db"
	mydb "github.com/seminhnva/my-grpc-go-server/internal/adapter/database"
	mygrpc "github.com/seminhnva/my-grpc-go-server/internal/adapter/grpc"
	app "github.com/seminhnva/my-grpc-go-server/internal/application"
	"github.com/seminhnva/my-grpc-go-server/internal/application/domain/bank"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(log.Writer())

	sqlDB, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/grpc?sslmode=disable")

	if err != nil {
		log.Fatalln("Can't connect database :", err)
	}

	if err := dbmigration.Migrate(sqlDB); err != nil {
		log.Fatalln("Database migration failed:", err)
	}

	dbAdapter, err := mydb.NewDatabaseAdapter(sqlDB)
	if err != nil {
		log.Fatalln("Can't create database adapter:", err)
	}

	hs := &app.HelloService{}
	bs := app.NewBankService(dbAdapter)

	if err := createDummyExchangeRate(bs, "USD", "VND"); err != nil {
		log.Fatalln("Failed to seed exchange rate:", err)
	}
	go generateExchangeRates(bs, "USD", "VND", 5*time.Second)
	grpcAdapter := mygrpc.NewGrpcAdapter(hs, bs, 9090)
	grpcAdapter.Run()
}

func createDummyExchangeRate(bs *app.BankService, fromCurrency, toCurrency string) error {
	return bs.CreateDummyExchangeRate(bank.ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         2000 + float64(rand.IntN(300)),
	})
}

func generateExchangeRates(bs *app.BankService, fromCurrency, toCurrency string, duration time.Duration) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for range ticker.C {
		if err := createDummyExchangeRate(bs, fromCurrency, toCurrency); err != nil {
			log.Println("Failed to generate exchange rate:", err)
		}
	}
}
