package controllers

import (
	"errors"
	"llio-api/models/DTOs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Custom error types
type ValidationError struct {
	Message string
	Code    int
}

func (e *ValidationError) Error() string {
	return e.Message
}

func validateUserIdParameter(c *gin.Context, userIDParam *string) (int, error) {
	// Get the current user from context
	userInfo, isExist := c.Get("current_user")
	if !isExist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non trouvé dans le contexte"})
		return 0, errors.New("utilisateur non trouvé dans le contexte")
	}

	// Check if userIDParam is provided and not empty
	if userIDParam == nil || *userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de l'utilisateur manquant"})
		return 0, errors.New("ID de l'utilisateur manquant")
	}

	// Convert string to int
	userIDInt, err := strconv.Atoi(*userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID d'utilisateur invalide"})
		return 0, errors.New("ID d'utilisateur invalide")
	}

	// Check if user is trying to modify their own account
	if userIDInt == userInfo.(*DTOs.UserDTO).Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Vous ne pouvez pas faire cette action sur votre propre compte"})
		return 0, errors.New("Vous ne pouvez pas faire cette action sur votre propre compte")
	}

	return userIDInt, nil
}
