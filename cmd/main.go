package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	dbmigration "github.com/seminhnva/my-grpc-go-server/db"
	mydb "github.com/seminhnva/my-grpc-go-server/internal/adapter/database"
	mygrpc "github.com/seminhnva/my-grpc-go-server/internal/adapter/grpc"
	app "github.com/seminhnva/my-grpc-go-server/internal/application"
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

	grpcAdapter := mygrpc.NewGrpcAdapter(hs, bs, 9090)
	grpcAdapter.Run()

}
