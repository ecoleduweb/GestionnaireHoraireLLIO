package controllers

import (
	"llio-api/models/DTOs"
	"llio-api/models/enums"
	"llio-api/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var projectSTR = "projet"

func CreatedProject(c *gin.Context) {
	var projetDTO DTOs.ProjectDTO

	msgErrsJson := services.VerifyJSON(c, &projetDTO)
	if len(msgErrsJson) > 0 {
		log.Printf("Une ou plusieurs erreurs de format JSON sont survenues:%v", msgErrsJson)
		c.JSON(http.StatusBadRequest, gin.H{"errors": msgErrsJson})
		return
	}

	msgErrsMore := services.VerifyProjectJSON(&projetDTO)
	if len(msgErrsMore) > 0 {
		log.Printf("Une ou plusieurs erreurs de verification du format du projet sont survenues:%v", msgErrsMore)
		c.JSON(http.StatusBadRequest, gin.H{"errors": msgErrsMore})
		return
	}

	_, err := services.GetUserById(projetDTO.ManagerId)
	if err != nil {
		handleError(c, err, userSTR)
		return
	}

	projectAdded, err := services.CreateProject(&projetDTO)
	if err != nil {
		handleError(c, err, projectSTR)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": "Le projet a bien été ajouté à la base de données",
		"project":  projectAdded,
	})
}

func GetProjectById(c *gin.Context) {
	id := c.Param("id")

	project, err := services.GetProjectById(id)
	if err != nil {
		handleError(c, err, projectSTR)
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

func GetDetailedProjects(c *gin.Context) {
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

	var projects []map[string]any
	var err error

	switch user.Role {
	case enums.Administrator:
		projects, err = services.GetDetailedProjects(user.Id)
	case enums.ProjectManager:
		projects, err = services.GetDetailedProjectsByManagerId(user.Id)
	case enums.Employee:
		projects, err = services.GetDetailedProjectsByUserId(user.Id)
	}

	if err != nil {
		handleError(c, err, projectSTR)
		return
	}

	if projects == nil {
		// Retourner une liste vide au lieu de null
		c.JSON(http.StatusOK, gin.H{"projects": []map[string]any{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func GetProjects(c *gin.Context) {
    currentUser, exists := c.Get("current_user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
        return
    }

    user, ok := currentUser.(*DTOs.UserDTO)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne du serveur"})
        return
    }

    projects, err := services.GetProjects(user.Id)
    if err != nil {
        handleError(c, err, projectSTR)
        return
    }

    if projects == nil {
        c.JSON(http.StatusOK, gin.H{"projects": []DTOs.ProjectDTO{}})
        return
    }

    c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func UpdateProject(c *gin.Context) {
	var projectToUpdate DTOs.ProjectDTO

	msgErrsJson := services.VerifyJSON(c, &projectToUpdate)
	if len(msgErrsJson) > 0 {
		log.Printf("Une ou plusieurs erreurs de format JSON sont survenues:%v", msgErrsJson)
		c.JSON(http.StatusBadRequest, gin.H{"errors": msgErrsJson})
		return
	}

	msgErrsMore := services.VerifyProjectJSON(&projectToUpdate)
	if len(msgErrsMore) > 0 {
		log.Printf("Une ou plusieurs erreurs de verification du format du projet sont survenues:%v", msgErrsMore)
		c.JSON(http.StatusBadRequest, gin.H{"errors": msgErrsMore})
		return
	}

	id := strconv.Itoa(projectToUpdate.Id)
	_, err := services.GetProjectById(id)
	if err != nil {
		handleError(c, err, projectSTR)
		return
	}

	_, err = services.GetUserById(projectToUpdate.ManagerId)
	if err != nil {
		handleError(c, err, userSTR)
		return
	}

	updatedProjectDTO, err := services.UpdateProject(&projectToUpdate)
	if err != nil {
		handleError(c, err, projectSTR)
		return
	}

	c.JSON(http.StatusOK, gin.H{"updatedProject": updatedProjectDTO})
}

func DeleteProject(c *gin.Context) {
	projectIdToDelete := c.Param("id")

	// Convert string to int
	projectIDInt, err := strconv.Atoi(projectIdToDelete)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page non trouvé - Erreur de conversion du ID en int"})
		return
	}

	err = services.DeleteProjectById(projectIDInt)
	if err != nil {
		handleError(c, err, projectSTR)
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

func GetDetailedProjectsByUser(c *gin.Context) {
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

	projects, err := services.GetDetailedProjectsByUserId(user.Id)
	if err != nil {
		handleError(c, err, projectSTR)
		return
	}

	if projects == nil {
		// Retourner une liste vide au lieu de null
		c.JSON(http.StatusOK, gin.H{"projects": []map[string]any{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}
