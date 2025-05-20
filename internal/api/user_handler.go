package api

import (
	"errors"
	"github.com/jdmukiibs/femProject/internal/store"
	"log"
	"regexp"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterRequest(reg *registerUserRequest) error {
	if reg.Username == "" {
		return errors.New("username is required")
	}

	if reg.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(reg.Email) {
		return errors.New("invalid email format")
	}

	if reg.Password == "" {
		return errors.New("password is required")
	}
	if len(reg.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}

// TODO: Implement remaining handlers
