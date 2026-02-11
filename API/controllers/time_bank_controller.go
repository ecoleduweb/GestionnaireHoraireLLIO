package controllers

import (
	"llio-api/models/DTOs"
	"llio-api/services"
	"net/http"

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

	response, err := services.CalculateTimeBank(user.Id, requestDTO)
	if err != nil {
		handleError(c, err, "banque d'heures")
		return
	}

	c.JSON(http.StatusOK, response)
}
