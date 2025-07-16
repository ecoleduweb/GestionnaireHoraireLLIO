package controllers

import (
	"llio-api/models/DTOs"
	"llio-api/models/enums"
	"llio-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var userSTR = "utilisateur"

func GetUserInfo(c *gin.Context) {
	userInfo, isExist := c.Get("current_user")
	if !isExist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non trouvé dans le contexte"})
		return
	}

	currentUser, ok := userInfo.(*DTOs.UserDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Type d'utilisateur invalide"})
		return
	}

	user, err := services.GetUserByEmail(currentUser.Email)
	if err != nil {
		handleError(c, err, userSTR)
		return
	}

	c.JSON(http.StatusOK, gin.H{"firstName": user.FirstName, "lastName": user.LastName, "role": user.Role})
}

func GetAllUsers(c *gin.Context) {
	roleParams := c.QueryArray("role")

	var userRoles []enums.UserRole
	for _, roleStr := range roleParams {
		roleValue, err := strconv.Atoi(roleStr)
		if err != nil || roleValue < 0 || roleValue > 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role invalide"})
			return
		}
		userRoles = append(userRoles, enums.UserRole(roleValue))
	}

	users, err := services.GetAllUsers(userRoles)
	if err != nil {
		handleError(c, err, userSTR)
		return
	}

	if users == nil {
		users = []*DTOs.UserDTO{}
	}

	c.JSON(http.StatusOK, users)
}

func UpdateUserRole(c *gin.Context) {
	// Get the current user from context
	userInfo, isExist := c.Get("current_user")
	if !isExist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non trouvé dans le contexte"})
		return
	}

	// Get the user ID from the URL parameters
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de l'utilisateur manquant"})
		return
	}

	userID_int, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID d'utilisateur invalide"})
		return
	}

	if userID_int == userInfo.(*DTOs.UserDTO).Id {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vous ne pouvez pas changer votre propre rôle"})
		return
	}

	// Parse the role from the request body
	var roleRequest struct {
		Role *int `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&roleRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rôle invalide: " + err.Error()})
		return
	}

	// Import the enums package and convert string to UserRole
	userRole, err := enums.ParseUserRole(*roleRequest.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rôle invalide: " + err.Error()})
		return
	}

	userDTO := DTOs.UserDTO{
		Id:   userID_int,
		Role: userRole,
	}

	// Update the user's role
	user, err := services.UpdateUserRole(&userDTO)
	if err != nil {
		handleError(c, err, userSTR)
		return
	}

	c.JSON(http.StatusOK, user)
}

func DeleteUserById(c *gin.Context) {
	// Get the current user from context
	userInfo, isExist := c.Get("current_user")
	if !isExist {
		c.JSON(http.StatusForbidden, gin.H{"error": "Utilisateur non trouvé dans le contexte"})
		return
	}

	// Get the user ID from the URL parameters
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "ID de l'utilisateur manquant"})
		return
	}

	userID_int, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "ID d'utilisateur invalide"})
		return
	}

	if userID_int == userInfo.(*DTOs.UserDTO).Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "Vous ne pouvez pas supprimer votre propre compte"})
		return
	}

	userHasActivities, err := services.UserHasActivities(userID_int)
	if err != nil {
		handleError(c, err, userSTR)
		return
	}

	if userHasActivities {
		c.JSON(http.StatusForbidden, gin.H{"error": "L'utilisateur ne peut pas être supprimé car il a des activités"})
		return
	}
	userHasProjects, err := services.UserHasProjects(userID_int)
	if err != nil {
		handleError(c, err, userSTR)
		return
	}

	if userHasProjects {
		c.JSON(http.StatusForbidden, gin.H{"error": "L'utilisateur ne peut pas être supprimé car il est gestionnaire d'un ou plusieurs projets"})
		return
	}

	userDTO, err := services.DeleteUserById(userID)
	if err != nil {
		handleError(c, err, userSTR)
		return
	}
	if userDTO == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Utilisateur supprimé avec succès", "user": userDTO})
}
