package controllers

import (
	"net/http"

	"llio-api/models/DTOs"
	"llio-api/models/enums"
	"llio-api/services"

	"github.com/gin-gonic/gin"
)

var rapportSTR = "rapport"

func ExportCSV(c *gin.Context) {
	currentUser, exists := c.Get("current_user")

	from := c.DefaultQuery("from", "")
	to := c.DefaultQuery("to", "")

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

    buf, err := services.GenerateExportCSV(from, to)
    if err != nil {
        handleError(c, err, rapportSTR)
        return
    }

    c.Header("Content-Type", "text/csv")
    c.Header("Content-Disposition", `attachment; filename="export.csv"`)
    c.Data(http.StatusOK, "text/csv", buf.Bytes())
}