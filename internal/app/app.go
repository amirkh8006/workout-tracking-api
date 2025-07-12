package app

import (
	"database/sql"
	"femProject/internal/api"
	"femProject/internal/middleware"
	"femProject/internal/store"
	"femProject/migrations"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Application struct {
	Logger *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler *api.UserHandler
	TokenHandler *api.TokenHandler
	MiddleWare middleware.UserMiddleware
	DB *sql.DB
}

func NewApplication() (*Application, error) {
	pgdb, err := store.Open()
	if err != nil {
		return nil, err
 	}

	err = store.MigrateFs(pgdb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout , "" , log.Ldate|log.Ltime)

	// our store
	workoutStore := store.NewPostgresWorkoutStore(pgdb)
	userStore := store.NewPostgresUserStore(pgdb)
	tokenStore := store.NewPostgresTokenStore(pgdb)

	// our hanlder
	workoutHanlder := api.NewWorkoutHandler(workoutStore , logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	middleWareHandler := middleware.UserMiddleware{UserStore: userStore}

	app := &Application{
		Logger: logger,
		WorkoutHandler: workoutHanlder,
		UserHandler: userHandler,
		TokenHandler: tokenHandler,
		MiddleWare: middleWareHandler,
		DB: pgdb,
	}

	return app, nil
}


func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}