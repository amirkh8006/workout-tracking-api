package routes

import (
	"femProject/internal/app"

	"github.com/go-chi/chi/v5"
)

func SetUpRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.MiddleWare.Authenticate)

		r.Get("/workouts/{id}" , app.MiddleWare.RequireUser(app.WorkoutHandler.HandleGetWorkoutById))
		r.Post("/workouts" , app.MiddleWare.RequireUser(app.WorkoutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}" , app.MiddleWare.RequireUser(app.WorkoutHandler.HandleUpdateWorkoutByID))
		r.Delete("/workouts/{id}" , app.MiddleWare.RequireUser(app.WorkoutHandler.HandleDeleteWorkoutById))

	})

	r.Get("/health", app.HealthCheck)
	

	r.Post("/users" , app.UserHandler.HanldeRegisterUser)

	r.Post("/tokens/authentication" , app.TokenHandler.HandleCreateToken)

	return r
}