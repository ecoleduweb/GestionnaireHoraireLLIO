package controllers

import (
	"net/http"

	"llio-api/models/DTOs"
	"llio-api/models/enums"
	"llio-api/services"

	"github.com/gin-gonic/gin"
)

var rapportSTR = "rapport"

func ExportExcel(c *gin.Context) {
	currentUser, exists := c.Get("current_user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	user, ok := currentUser.(*DTOs.UserDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne du serveur"})
		c.Abort()
		return
	}

	if user.Role != enums.Administrator{
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur n'est pas Administrateur'"})
		return
	}

	// Headers download
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=activities.xlsx")
	c.Header("Transfer-Encoding", "chunked") // 🔥 important pour gros fichiers

    err := services.GenerateExcel(c.Writer)
    if err != nil {
        handleError(c, err, rapportSTR)
        return
    }
}