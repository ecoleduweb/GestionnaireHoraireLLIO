package controllers

import (
	"llio-api/models/DTOs"
	"llio-api/services"
	"log"
	"strconv"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

var activiteSTR = "activité"

func CreateActivity(c *gin.Context) {

	var activityDTO DTOs.ActivityDTO

	//Validation des données entrantes
	messageErrsJSON := services.VerifyJSON(c, &activityDTO)
	if len(messageErrsJSON) > 0 {
		log.Printf("Une ou plusieurs erreurs de format JSON sont survenues:%v", messageErrsJSON)
		c.JSON(http.StatusBadRequest, gin.H{"errors": messageErrsJSON})
		return
	}

	user, shouldReturn := getUserFromContext(c)
	if shouldReturn {
		return
	}

	messageErrs := services.VerifyActivityJSON(&activityDTO)
	if len(messageErrs) > 0 {
		log.Printf("Une ou plusieurs erreurs de verification du format de l'activité sont survenues:%v", messageErrs)
		c.JSON(http.StatusBadRequest, gin.H{"errors": messageErrs})
		return
	}

	activityAded, err := services.CreateActivity(&activityDTO, user.Id)
	if err != nil {
		handleError(c, err, activiteSTR)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"reponse":  "L'activité a bien été ajoutée à la base de données.",
		"activity": activityAded,
	})
}

func GetUsersActivities(c *gin.Context) {

	user, shouldReturn := getUserFromContext(c)
	if shouldReturn {
		return
	}
	activities, err := services.GetUsersActivities(user.Id)
	if err != nil {
		handleError(c, err, activiteSTR)
		return
	}

	// Retourne une liste vide plutôt que null pour être compatible avec les vérifiacations du frontend
	if activities == nil {
		// Retourner une liste vide au lieu de null
		c.JSON(http.StatusOK, gin.H{"activities": []DTOs.ActivityDTO{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"activities": activities})
}

func GetActivityById(c *gin.Context) {
	// Récupérer l'id de l'activité
	id := c.Param("id")
	activity, err := services.GetActivityById(id)
	if err != nil {
		handleError(c, err, activiteSTR)
		return
	}

	c.JSON(http.StatusOK, gin.H{"activity": activity})
}

func GetUsersOutlookCalendar(c *gin.Context) {
	userDto, shouldReturn := getUserFromContext(c)
	if shouldReturn {
		return
	}

	graphToken, err := services.GetUserGraphAccessToken(userDto.Id)
	if err != nil {
		handleError(c, err, activiteSTR)
		return
	}

	if graphToken == nil {
		log.Printf("La connexion Microsoft de l'utilisateur est invalide ou inexistante. : %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Votre compte n'est pas connecté à votre compte Microsoft. Veuillez vous reconnecter et réessayer."})
		return
	}

	isGraphTokenExpired, err := services.IsJWTExpired(*graphToken)
	if (err != nil) || isGraphTokenExpired {
		log.Printf("La connexion Microsoft de l'utilisateur est invalide ou inexistante. : %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Votre compte Microsoft s'est déconnecté. Veuillez vous reconnecter et réessayer."})
		return
	}

	events, err := services.GetCalendarEvents(*graphToken, time.Now())
	if err != nil {
		handleError(c, err, activiteSTR)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"events": events})
}

func UpdateActivity(c *gin.Context) {
	var updateActivityDTO DTOs.ActivityDTO

	//Validation des données entrantes
	messageErrsJSON := services.VerifyJSON(c, &updateActivityDTO)
	if len(messageErrsJSON) > 0 {
		log.Printf("Une ou plusieurs erreurs de format JSON sont survenues:%v", messageErrsJSON)
		c.JSON(http.StatusBadRequest, gin.H{"errors": messageErrsJSON})
		return
	}

	messageErrs := services.VerifyActivityJSON(&updateActivityDTO)
	if len(messageErrs) > 0 {
		log.Printf("Une ou plusieurs erreurs de verification du format de l'activité sont survenues:%v", messageErrs)
		c.JSON(http.StatusBadRequest, gin.H{"errors": messageErrs})
		return
	}

	id := strconv.Itoa(updateActivityDTO.Id)
	_, err := services.GetActivityById(id)
	if err != nil {
		handleError(c, err, activiteSTR)
		return
	}

	updatedActivityDTO, err := services.UpdateActivity(&updateActivityDTO)
	if err != nil {
		handleError(c, err, activiteSTR)
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated_activity": updatedActivityDTO})

}

func DeleteActivity(c *gin.Context) {

	id := c.Param("id")

	user, shouldReturn := getUserFromContext(c)

	if shouldReturn {
		return
	}

	activity, err := services.GetActivityById(id)
	if err != nil {
		handleError(c, err, activiteSTR)
		return
	}
	if activity == nil {
		log.Printf("L'activité à supprimer n'existe pas. : %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "L'activité à supprimer n'existe pas."})
		return
	}

	err = services.DeleteActivity(id, user)
	if err != nil {
		handleError(c, err, activiteSTR)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "La suppression de l'activité est un succès."})
}

func GetActivitiesFromRange(c *gin.Context) {
	from := c.Query("startDate")
	to := c.Query("endDate")

	user, shouldReturn := getUserFromContext(c)
	if shouldReturn {
		return
	}

	activities, err := services.GetActivitiesFromRange(from, to, user.Id)
	if err != nil {
		handleError(c, err, activiteSTR)
		return
	}

	if activities == nil {
		c.JSON(http.StatusOK, gin.H{"activities": []DTOs.ActivityDTO{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"activities": activities})
}

func GetDetailedActivitiesFromRange(c *gin.Context) {
	from := c.Query("startDate")
	to := c.Query("endDate")

	user, shouldReturn := getUserFromContext(c)

	if shouldReturn {
		return
	}

	detailedActivities, err := services.GetDetailedActivitiesFromRange(from, to, user.Id)
	if err != nil {
		handleError(c, err, activiteSTR)
		return
	}

	if detailedActivities == nil {
		c.JSON(http.StatusOK, gin.H{"activities": []DTOs.DetailedActivityDTO{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"activities": detailedActivities})
}

func getUserFromContext(c *gin.Context) (*DTOs.UserDTO, bool) {
	currentUser, exists := c.Get("current_user")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "Échec lors de la sauvegarde de l'utilisateur dans la session"})
		return nil, true
	}
	user, ok := currentUser.(*DTOs.UserDTO)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne du serveur"})
		return nil, true
	}
	return user, false
}
