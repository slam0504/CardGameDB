package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"CardGameDB/internal/application"
	"CardGameDB/internal/domain/card"
	"CardGameDB/internal/infrastructure/eventbus"
	"CardGameDB/internal/infrastructure/eventstore"
	mysqlRepo "CardGameDB/internal/infrastructure/repository/mysql"
	httpInterface "CardGameDB/internal/interface/http"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "user:password@tcp(localhost:3306)/carddb"
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := mysqlRepo.New(db)

	store := eventstore.NewMySQL(db)
	bus := eventbus.New(store)
	searchUC := application.NewSearchUseCase(repo)
	manageUC := application.NewManageUseCase(repo)
	bus.Subscribe("card.search", func(e interface{}) {
		if evt, ok := e.(card.SearchRequested); ok {
			searchUC.Handle(evt)
		}
	})
	bus.Subscribe("card.create", func(e interface{}) {
		if evt, ok := e.(card.CreateRequested); ok {
			manageUC.HandleCreate(evt)
		}
	})
	bus.Subscribe("card.update", func(e interface{}) {
		if evt, ok := e.(card.UpdateRequested); ok {
			manageUC.HandleUpdate(evt)
		}
	})

	server := httpInterface.NewServer(bus)
	log.Println("Listening on :8080")
	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
