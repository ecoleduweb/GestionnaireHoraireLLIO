package controllers

import (
	"fmt"
	"llio-api/customs_errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleError(ctx *gin.Context, err error, subject string) {
	switch err {
	case customs_errors.ErrDuplicateEntry:
		errorMsg := fmt.Sprintf("Un(e) %s avec ces données existe déjà", subject)
		log.Printf("ERREUR - Conflit de données: %s - %v", errorMsg, err)
		ctx.JSON(http.StatusConflict, gin.H{"error": errorMsg})

	case customs_errors.ErrNotFound:
		errorMsg := fmt.Sprintf("Le(La) %s n'a pas été trouvé(e)", subject)
		log.Printf("ERREUR - Ressource non trouvée: %s - %v", errorMsg, err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": errorMsg})

	case customs_errors.ErrDatabase:
		errorMsg := fmt.Sprintf("Erreur interne lors du traitement du(de la) %s", subject)
		log.Printf("ERREUR CRITIQUE - Erreur base de données: %s - %v", errorMsg, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})

	case customs_errors.ErrUserHasActivities:
		errorMsg := fmt.Sprintf("L'utilisateur a une ou des activités associées, suppression impossible")
		log.Printf("ERREUR - Suppression impossible: %s - %v", errorMsg, err)
		ctx.JSON(http.StatusForbidden, gin.H{"error": errorMsg})

	case customs_errors.ErrProjectHasActivities:
		errorMsg := fmt.Sprintf("Le projet a une ou des activités associées, suppression impossible")
		log.Printf("ERREUR - Suppression impossible: %s - %v", errorMsg, err)
		ctx.JSON(http.StatusForbidden, gin.H{"error": errorMsg})

	case customs_errors.ErrUserHasProjects:
		errorMsg := fmt.Sprintf("L'utilisateur a un ou des projets associées, suppression impossible")
		log.Printf("ERREUR - Suppression impossible: %s - %v", errorMsg, err)
		ctx.JSON(http.StatusForbidden, gin.H{"error": errorMsg})

	default:
		errorMsg := fmt.Sprintf("Erreur inconnue lors du traitement du(de la) %s", subject)
		log.Printf("ERREUR INCONNUE: %s - %v", errorMsg, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
	}
}
