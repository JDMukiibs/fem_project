package api

import (
	"encoding/json"
	"github.com/jdmukiibs/femProject/internal/store"
	"github.com/jdmukiibs/femProject/internal/utils"
	"log"
	"net/http"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

type Response struct {
	Message string `json:"message"`
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutById: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"workout": workout,
	})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: decodingCreateWorkout: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: createWorkout: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout update id"})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutById: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	// We found an existing workout
	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		wh.logger.Printf("ERROR: decodingUpdateRequest: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
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

	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.logger.Printf("ERROR: updatingWorkout: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkoutById(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Printf("ERROR: readIDParam: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	err = wh.workoutStore.DeleteWorkout(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: deleteWorkout: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "workout deleted"})
}
