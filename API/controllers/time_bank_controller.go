package controllers

import (
	"llio-api/models/DTOs"
	"llio-api/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CalculateTimeBank(c *gin.Context) {

	user, shouldReturn := getUserFromContext(c)
	if shouldReturn {
		return
	}

	var requestDTO DTOs.TimeBankRequestDTO
	if err := c.ShouldBindJSON(&requestDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides : startDate et hoursPerWeek sont requis"})
		return
	}

	messageErrsJSON := services.VerifyJSON(c, &requestDTO)
	if len(messageErrsJSON) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": messageErrsJSON})
		return
	}

	if _, err := time.Parse("2006-01-02", requestDTO.StartDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format de date invalide (attendu: YYYY-MM-DD)"})
		return
	}

	response, err := services.CalculateTimeBank(user.Id, requestDTO)
	if err != nil {
		handleError(c, err, "erreur lors du calcul de banque d'heures")
		return
	}

	c.JSON(http.StatusOK, response)
}
