package handlers

import (
	"fmt"
	"net/http"

	"github.com/cjnghn/db-shard-example/internal/models"
	"github.com/labstack/echo/v4"
)

type createUserRequestBody struct {
	Name string `json:"name" validate:"required"`
}

type createUserOKResponse struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

type getUserOKResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func CreateUserHandler(c echo.Context) error {
	var body createUserRequestBody
	if err := c.Bind(&body); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request format")
	}

	if body.Name == "" {
		return c.String(http.StatusBadRequest, "Name is required")
	}

	userID, err := models.CreateUser(body.Name)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
	}

	return c.JSON(http.StatusOK, &createUserOKResponse{
		UserID:  userID,
		Message: fmt.Sprintf("User %s created with ID %s", body.Name, userID),
	})
}

func GetUserHandler(c echo.Context) error {
	userID := c.Param("id")

	name, err := models.GetUser(userID)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get user: %v", err))
	}

	return c.JSON(http.StatusOK, &getUserOKResponse{
		UserID: userID,
		Name:   name,
	})
}

func GetAllUsersHandler(c echo.Context) error {
	users, err := models.GetAllUsers()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get all users: %v", err))
	}

	return c.JSON(http.StatusOK, users)
}
