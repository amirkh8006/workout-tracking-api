package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"femProject/internal/middleware"
	"femProject/internal/store"
	"femProject/internal/utils"
	"log"
	"net/http"
)

type WorkoutHandler struct{
	workoutStore store.WorkoutStore
	logger *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger: logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam %v", err)
		utils.WriteJson(w, http.StatusBadRequest , utils.Envlope{"error": "Invalid WorkoutId"})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutByID %v", err)
		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "Internal Server Error"})
		return
	}

	if workout == nil {
		http.NotFound(w, r)
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envlope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: Decoding Body %v", err)
		utils.WriteJson(w, http.StatusBadRequest , utils.Envlope{"error": "Invalid Request sent"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJson(w, http.StatusBadRequest , utils.Envlope{"error": "You must be logged in"})
		return
	}

	workout.UserID = currentUser.ID

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: Create Workout %v", err)
		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "Failed to create workout"})
		return
	}

	utils.WriteJson(w , http.StatusCreated, utils.Envlope{"Workout": createdWorkout})
}


func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam %v", err)
		utils.WriteJson(w, http.StatusBadRequest , utils.Envlope{"error": "Invalid WorkoutId"})
		return
	}


	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Printf("ERROR: Update workout By ID %v", err)
		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "Failed to fetch workout"})
		return
	}

	if existingWorkout == nil {
		wh.logger.Printf("ERROR: update Workout %v", err)
		utils.WriteJson(w, http.StatusNotFound , utils.Envlope{"error": "Workout Not Found"})
		return
	}

	var updateWorkoutRequest struct {
		ID *int `json:"id"`
		Title *string `json:"title"`
		Description *string `json:"description"`
		DurationMinutes *int `json:"duration_minutes"`
		CaloriesBurned *int `json:"calories_burned"`
		Entries []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)

	if err != nil {
		wh.logger.Printf("ERROR: Update Workout Decode Json %v", err)
		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "Internal Server Error"})
		return
	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}

	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}

	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}

	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}

	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJson(w, http.StatusBadRequest , utils.Envlope{"error": "You must be logged in"})
		return
	}	

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJson(w, http.StatusNotFound , utils.Envlope{"error": "Workout does not exist"})
			return
		}

		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "Internal server error"})
		return
	}

	if workoutOwner != currentUser.ID {
		utils.WriteJson(w, http.StatusForbidden , utils.Envlope{"error": "You are not authorized to update this workout"})
		return
	}

	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.logger.Printf("ERROR: Update Workout Error%v", err)
		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "Failed to update workout"})
		return
	}

	utils.WriteJson(w , http.StatusOK, utils.Envlope{"Workout": existingWorkout})
}


func (wh *WorkoutHandler) HandleDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIdParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam Delete Workout%v", err)
		utils.WriteJson(w, http.StatusBadRequest , utils.Envlope{"error": "Invalid WorkoutId"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJson(w, http.StatusBadRequest , utils.Envlope{"error": "You must be logged in"})
		return
	}	

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJson(w, http.StatusNotFound , utils.Envlope{"error": "Workout does not exist"})
			return
		}

		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "Internal server error"})
		return
	}

	if workoutOwner != currentUser.ID {
		utils.WriteJson(w, http.StatusForbidden , utils.Envlope{"error": "You are not authorized to delete this workout"})
		return
	}

	err = wh.workoutStore.DeleteWorkoutByID(workoutID)

	if err == sql.ErrNoRows {
		wh.logger.Printf("ERROR: Delete Workout Error%v", err)
		utils.WriteJson(w, http.StatusNotFound , utils.Envlope{"error": "Workout Not Found"})
		return
	}

	
	if err != nil {
		wh.logger.Printf("ERROR: Delete Workout Error%v", err)
		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "Internal Server Error"})
		return
	}

	utils.WriteJson(w , http.StatusNoContent, utils.Envlope{})

}
