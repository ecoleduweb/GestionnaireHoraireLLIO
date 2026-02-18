package controllers

import (
	"llio-api/models/DTOs"
	"llio-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

const timeBankSubject = "banque d'heures"

// GET /user/time-bank
func GetTimeBankBalance(c *gin.Context) {
	user, shouldReturn := getUserFromContext(c)
	if shouldReturn {
		return
	}

	balance, err := services.GetTimeBankBalance(user.Id)
	if err != nil {
		handleError(c, err, timeBankSubject)
		return
	}

	c.JSON(http.StatusOK, balance)
}

// GET /user/time-bank/config
func GetTimeBankConfig(c *gin.Context) {
	user, shouldReturn := getUserFromContext(c)
	if shouldReturn {
		return
	}

	config, err := services.GetTimeBankConfig(user.Id)
	if err != nil {
		handleError(c, err, timeBankSubject)
		return
	}

	// Si pas de config, on renvoie null explicitement (code 200)
	c.JSON(http.StatusOK, config)
}

// POST /user/time-bank/config
func SaveTimeBankConfig(c *gin.Context) {
	user, shouldReturn := getUserFromContext(c)
	if shouldReturn {
		return
	}

	var dto DTOs.TimeBankConfigDTO
	// Validation automatique (datetime=2006-01-02)
	msgErrs := services.VerifyJSON(c, &dto)
	if len(msgErrs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": msgErrs})
		return
	}

	savedConfig, err := services.SaveTimeBankConfig(user.Id, dto)
	if err != nil {
		handleError(c, err, timeBankSubject)
		return
	}

	c.JSON(http.StatusOK, savedConfig)
}
