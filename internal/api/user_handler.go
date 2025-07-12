package api

import (
	"encoding/json"
	"errors"
	"femProject/internal/store"
	"femProject/internal/utils"
	"log"
	"net/http"
	"regexp"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Bio string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger: logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if len(req.Username) > 50 {
		return errors.New("username cannot be greater than 50")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("email is not valid")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	return nil
}


func (h *UserHandler) HanldeRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: Decoding Body %v", err)
		utils.WriteJson(w, http.StatusBadRequest , utils.Envlope{"error": "Invalid Request sent"})
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest , utils.Envlope{"error": err.Error()})
		return
	}

	user := &store.User{
		UserName: req.Username,
		Email: req.Email,
	}

	if req.Bio != "" {
		user.Bio = req.Bio
	}

	// hash password
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: hashing password %v", err)
		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "internal server error"})
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("ERROR: registering user %v", err)
		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusCreated, utils.Envlope{"user": user})
}