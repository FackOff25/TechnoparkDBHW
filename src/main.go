package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"

	"github.com/FackOff25/TechnoparkDBHW/src/server"

	forumDelivery "github.com/FackOff25/TechnoparkDBHW/src/internal/forum/delivery"
	forumRep "github.com/FackOff25/TechnoparkDBHW/src/internal/forum/repository"
	forumUsecase "github.com/FackOff25/TechnoparkDBHW/src/internal/forum/usecase"
	postDelivery "github.com/FackOff25/TechnoparkDBHW/src/internal/post/delivery"
	postRep "github.com/FackOff25/TechnoparkDBHW/src/internal/post/repository"
	postUsecase "github.com/FackOff25/TechnoparkDBHW/src/internal/post/usecase"
	serviceDelivery "github.com/FackOff25/TechnoparkDBHW/src/internal/service/delivery"
	serviceRep "github.com/FackOff25/TechnoparkDBHW/src/internal/service/repository"
	serviceUsecase "github.com/FackOff25/TechnoparkDBHW/src/internal/service/usecase"
	threadDelivery "github.com/FackOff25/TechnoparkDBHW/src/internal/thread/delivery"
	threadRep "github.com/FackOff25/TechnoparkDBHW/src/internal/thread/repository"
	threadUsecase "github.com/FackOff25/TechnoparkDBHW/src/internal/thread/usecase"
	userDelivery "github.com/FackOff25/TechnoparkDBHW/src/internal/user/delivery"
	userRep "github.com/FackOff25/TechnoparkDBHW/src/internal/user/repository"
	userUsecase "github.com/FackOff25/TechnoparkDBHW/src/internal/user/usecase"
)

var postgresConfig = postgres.Config{DSN: "host=localhost user=db_perf_user password=db_perf_password database=db_perf_project port=5432"}

func main() {
	db, err := gorm.Open(postgres.New(postgresConfig),
		&gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	forumDB := forumRep.New(db)
	userDB := userRep.New(db)
	postDB := postRep.New(db)
	threadDB := threadRep.New(db)
	serviceDB := serviceRep.New(db)

	forumUC := forumUsecase.New(forumDB, userDB)
	userUC := userUsecase.New(userDB)
	postUC := postUsecase.New(postDB, userDB, threadDB, forumDB)
	threadUC := threadUsecase.New(threadDB, userDB, forumDB)
	serviceUC := serviceUsecase.New(serviceDB)

	e := echo.New()

	e.Logger.SetHeader(`time=${time_rfc3339} level=${level} prefix=${prefix} ` +
		`file=${short_file} line=${line} message:`)
	e.Logger.SetLevel(elog.INFO)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `time=${time_custom} remote_ip=${remote_ip} ` +
			`host=${host} method=${method} uri=${uri} user_agent=${user_agent} ` +
			`status=${status} error="${error}" ` +
			`bytes_in=${bytes_in} bytes_out=${bytes_out}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	e.Use(middleware.Recover())

	forumDelivery.MakeDelivery(e, forumUC)
	userDelivery.MakeDelivery(e, userUC)
	postDelivery.MakeDelivery(e, postUC)
	threadDelivery.MakeDelivery(e, threadUC)
	serviceDelivery.MakeDelivery(e, serviceUC)

	s := server.NewServer(e)
	if err := s.Start(); err != nil {
		e.Logger.Fatal(err)
	}
}
