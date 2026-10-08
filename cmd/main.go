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

	dbmigration.Migrate(sqlDB)

	dbAdapter, err := mydb.NewDatabaseAdapter(sqlDB)
	if err != nil {
		log.Fatalln("Can't create database adapter:", err)
	}

	hs := &app.HelloService{}
	bs := app.NewBankService(dbAdapter)

	// go generateExchangeRates(bs, "USD", "VND", 5*time.Second)
	grpcAdapter := mygrpc.NewGrpcAdapter(hs, bs, 9090)
	grpcAdapter.Run()

}

func generateExchangeRates(bs *app.BankService, fromCurrency, toCurrency string, duration time.Duration) {
	ticker := time.NewTicker(duration)

	for range ticker.C {
		now := time.Now()
		validFrom := now.Truncate(time.Second).Add(3 * time.Second)
		validTo := validFrom.Add(duration).Add(-1 * time.Millisecond)

		dummyRate := bank.ExchangeRate{
			FromCurrency:       fromCurrency,
			ToCurrency:         toCurrency,
			ValidFromTimeStamp: validFrom,
			ValidToTimeStamp:   validTo,
			Rate:               2000 + float64(rand.IntN(300)),
		}

		bs.CreateDummyExchangeRate(dummyRate)
	}
}
