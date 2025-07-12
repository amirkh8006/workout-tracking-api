package api

import (
	"encoding/json"
	"femProject/internal/store"
	"femProject/internal/tokens"
	"femProject/internal/utils"
	"log"
	"net/http"
	"time"
)

type TokenHandler struct {
	tokenStore store.TokensStore
	userStore store.UserStore
	logger *log.Logger
}

type CreateTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokensStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore: userStore,
		logger: logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req CreateTokenRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: Create Token Request %v", err)
		utils.WriteJson(w, http.StatusBadRequest , utils.Envlope{"error": "Invalid Request sent"})
		return
	}

	// get the user
	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: Get user by username %v", err)
		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "Internal server error"})
		return
	}

	passwordsDoMatch, err := user.PasswordHash.Matched(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: password hash match %v", err)
		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "Internal server error"})
		return
	}

	if !passwordsDoMatch {
		utils.WriteJson(w, http.StatusUnauthorized , utils.Envlope{"error": "Invalid Credentials"})
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Printf("ERROR: creating token %v", err)
		utils.WriteJson(w, http.StatusInternalServerError , utils.Envlope{"error": "Internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusCreated , utils.Envlope{"auth_token": token})

}