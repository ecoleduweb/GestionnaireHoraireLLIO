package controllers

import (
	"llio-api/models/DTOs"
	"llio-api/models/enums"
	"llio-api/services"
	"llio-api/useful"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func GetAuthCallback(c *gin.Context) {
	useful.LoadEnv()
	frontendURL := os.Getenv("FRONTEND_URL")
	jwtDurationHours, err := strconv.Atoi(os.Getenv("JWT_DURATION"))
	if err != nil {
		log.Printf("Erreur lors de la conversion de JWT_DURATION: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Durée JWT invalide."})
		return
	}

	useful.SetupAuthProvider(c)
	userAzure, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		log.Printf("Erreur lors de l'authentification: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var dbUser DTOs.UserDTO
	dbUser.FirstName = userAzure.FirstName
	dbUser.LastName = userAzure.LastName
	dbUser.Email = userAzure.Email
	dbUser.Role = enums.Employee

	userInDb, err := services.FirstOrCreateUser(&dbUser)
	if err != nil {
		log.Printf("Impossible d'interagir avec l'utilisateur dans la base de données: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible d'ajouter l'utilisateur à la base de données."})
		return
	}

	tokenDuration := time.Now().Add(time.Duration(jwtDurationHours) * time.Hour)
	accessToken, err := services.CreateJWTToken(userInDb.Id, userAzure.Email, userAzure.FirstName, userAzure.LastName, tokenDuration, userInDb.Role)
	if err != nil {
		log.Printf("Erreur lors de l'authentification: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	isRunningSecure := useful.IsRunningSecure()

	// On sauvegarde le refresh token de l'utilisateur afin d'obtenir
	// les évènements Outlook à l'aide de celui-ci plus tard
	services.UpdateUserGraphAccessToken(userInDb.Id, &userAzure.AccessToken)

	// cookie pour l'authentification de l'utilisateur
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   isRunningSecure,
		SameSite: http.SameSiteStrictMode,
	})
	http.Redirect(c.Writer, c.Request, frontendURL+"/calendar", http.StatusFound)
}

func Auth(c *gin.Context) {
	useful.SetupAuthProvider(c)
	if gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
		log.Printf("Utilisateur déja authentifié")
		c.JSON(http.StatusOK, gin.H{"access_token": gothUser.AccessToken})
		return
	} else {
		log.Printf("Début de l'authentification")
		gothic.BeginAuthHandler(c.Writer, c.Request)
	}
}

func Logout(c *gin.Context) {
	useful.LoadEnv()
	if err := gothic.Logout(c.Writer, c.Request); err != nil {
		log.Printf("Erreur lors de la déconnexion: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   useful.IsRunningSecure(),
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Déconnexion réussie"})
}
