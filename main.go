package main

import (
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"CardGameDB/internal/application"
	"CardGameDB/internal/domain/card"
	"CardGameDB/internal/infrastructure/commandbus"
	"CardGameDB/internal/infrastructure/eventstore"
	mysqlRepo "CardGameDB/internal/infrastructure/repository/mysql"
	httpInterface "CardGameDB/internal/interface/http"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "user:password@tcp(localhost:3306)/carddb?parseTime=true"
	}
	gdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	repo := mysqlRepo.New(gdb)

	store := eventstore.NewMySQL(gdb)

	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		brokersEnv = "localhost:9092"
	}
	brokers := strings.Split(brokersEnv, ",")
	bus, err := commandbus.NewKafka(store, brokers)
	if err != nil {
		log.Fatal(err)
	}
	searchUC := application.NewSearchUseCase(repo)
	manageUC := application.NewManageUseCase(repo)
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

	server := httpInterface.NewServer(bus, searchUC)
	log.Println("Listening on :8080")
	if err := server.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
