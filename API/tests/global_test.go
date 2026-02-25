package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"llio-api/database"
	"llio-api/handlers"
	"llio-api/models/DAOs"
	"llio-api/models/DTOs"
	"llio-api/models/enums"
	"llio-api/routes"
	"llio-api/services"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// variables globales utilisées pour led tests
var (
	router              *gin.Engine
	w                   *httptest.ResponseRecorder
	doNotDeleteUser     DAOs.User
	doNotDeleteUser2    DAOs.User
	pleaseDeleteUser    DAOs.User
	pleaseDeleteProject DAOs.Project
	doNotDeleteCategory DAOs.Category
	doNotDeleteProject  DAOs.Project
	doNotDeleteProject2 DAOs.Project
	doNotDeleteActivity DAOs.Activity
)

// Global JWT token for authentication in tests
var accessToken string

// Create and set JWT token for tests

func createAndSetAccessToken(role enums.UserRole, userId int) {
	// Create a JWT token for the test user
	token, err := services.CreateJWTToken(userId, doNotDeleteUser.Email, doNotDeleteUser.FirstName, doNotDeleteUser.LastName, time.Now().Add(time.Hour), role)
	if err != nil {
		log.Fatalf("Failed to create JWT token: %v", err)
	}
	accessToken = token
}

// agit comme un gros before each
func TestMain(m *testing.M) {
	r, recorder := setupTestRouter()
	router = r
	w = recorder

	//Ajoute des enregistrements de tests
	prepareTestData()
	// Permet de rouler tous les tests. Aucun test ne fonctionne sans  ça.
	exitCode := m.Run()

	// TODO vider la BD après les tests.
	drop_all_tables()
	// affiche le code de sortie des tests
	os.Exit(exitCode)
}

func prepareTestData() {
	testUser := DAOs.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Role:      enums.Administrator,
	}
	database.DB.Create(&testUser)
	doNotDeleteUser = testUser
	testUser2 := DAOs.User{
		FirstName: "Johnny",
		LastName:  "Joestar",
		Email:     "tusk@example.com",
		Id:        20, // Assurez-vous que l'ID est unique pour le test
	}
	database.DB.Create(&testUser2)
	pleaseDeleteUser = testUser2
	testUser3 := DAOs.User{
		Id:        3,
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Role:      enums.ProjectManager,
	}
	database.DB.Create(&testUser3)
	doNotDeleteUser2 = testUser3
	testProject := DAOs.Project{
		Id:             1, // Assurez-vous que l'ID est unique pour le test
		UniqueId:       "Interne-1234",
		Name:           "Sample project",
		ManagerId:      doNotDeleteUser.Id,
		Status:         enums.ProjectStatus(enums.InProgress),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		EndAt:          time.Now(),
		EstimatedHours: 10, // Ajout d'une valeur pour EstimatedHours
	}
	database.DB.Create(&testProject)
	doNotDeleteProject = testProject

	testProject2 := DAOs.Project{
		Id:             2, // Assurez-vous que l'ID est unique pour le test
		UniqueId:       "Externe-6789",
		ManagerId:      doNotDeleteUser.Id,
		Name:           "Sample project 2",
		Status:         enums.ProjectStatus(enums.InProgress),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		EndAt:          time.Now(),
		EstimatedHours: 100, // Ajout d'une valeur pour EstimatedHours
	}
	database.DB.Create(&testProject2)
	doNotDeleteProject2 = testProject2

	testProject3 := DAOs.Project{
		Id:             3, // Assurez-vous que l'ID est unique pour le test
		UniqueId:       "Interne-4444",
		Name:           "Sample yes",
		ManagerId:      doNotDeleteUser.Id,
		Status:         enums.ProjectStatus(enums.InProgress),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		EndAt:          time.Now(),
		EstimatedHours: 10, // Ajout d'une valeur pour EstimatedHours
	}
	database.DB.Create(&testProject3)
	pleaseDeleteProject = testProject3

	log.Println("doNotDeleteProject2:", doNotDeleteProject2)
	log.Println("doNotDeleteProject:", doNotDeleteProject)

	testCategory := DAOs.Category{
		Name:        "Test Category",
		Description: "Sample category",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ProjectId:   doNotDeleteProject.Id,
	}
	database.DB.Create(&testCategory)
	doNotDeleteCategory = testCategory

	testActivity := DAOs.Activity{
		Name:        "Test Activity",
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(time.Hour),
		Description: "test description",
		UserId:      doNotDeleteUser.Id,
		ProjectId:   doNotDeleteProject.Id,
		CategoryId:  doNotDeleteCategory.Id,
	}
	database.DB.Create(&testActivity)
	doNotDeleteActivity = testActivity
}

// Change de répertoir pour trouver le .env
func changeCurrentDiretory() {
	// Obtenir le répertoire courant actuel
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Erreur lors de la récupération du répertoire courant :", err)
		return
	}
	fmt.Println("Répertoire courant actuel :", currentDir)

	// Vérifier si le répertoire courant contient déjà le dossier "API"
	if !strings.HasSuffix(currentDir, "API") {
		// Remonter d'un niveau dans l'arborescence des répertoires pour trouver le .env
		err = os.Chdir("..")
		if err != nil {
			log.Fatalf("Erreur lors du changement de répertoire :%v", err)
			return
		}

		// Vérifier le nouveau répertoire courant
		updatedDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Erreur lors de la récupération du répertoire courant :", err)
			return
		}
		fmt.Println("Nouveau répertoire courant :", updatedDir)
	} else {
		fmt.Println("Déjà dans le répertoire API.")
	}
}

// setupTestRouter initialise un routeur de test et un enregistreur de réponse
func setupTestRouter() (*gin.Engine, *httptest.ResponseRecorder) {
	changeCurrentDiretory()
	os.Setenv("ENV", "TEST")
	database.Connect()

	router := gin.Default()
	// used for the api routes
	routes.RegisterRoutes(router)
	// used for the health check
	handlers.ApiStatus(router)

	return router, httptest.NewRecorder()
}

// sendRequest envoie une requête HTTP au routeur de test
// pour créer une requête http avec un role administrateur, on ajoute le role voulu à la fin : sendRequest(router, "POST", "/activity", activity, enums.Employee)
func sendRequest(router *gin.Engine, method, url string, body interface{}, userId *int, userRole ...enums.UserRole) *httptest.ResponseRecorder {
	var req *http.Request
	// If accessToken exists, we need to add it to the request cookies
	// This will be used in non-authenticated helper functions
	if method == "GET" || body == nil {
		// Pour GET ou body nil, ne pas inclure de corps
		req, _ = http.NewRequest(method, url, nil)
	} else {
		// Pour les autres méthodes avec body non-nil
		jsonValue, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(jsonValue))
	}

	var finalUserId int
	if userId != nil {
		finalUserId = *userId
	} else {
		finalUserId = doNotDeleteUser.Id
	}

	if (userRole != nil) && len(userRole) > 0 {
		createAndSetAccessToken(userRole[0], finalUserId)
	} else {
		createAndSetAccessToken(enums.Employee, finalUserId)
	}

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Path:     "/",
	}
	// Add cookies to the request if any exist
	req.AddCookie(cookie)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	return w
}

// assertResponse vérifie si le code de statut de la réponse et les erreurs correspondent à ceux attendus
func assertResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus int, expectedErrors []DTOs.FieldErrorDTO) {
	assert.Equal(t, expectedStatus, w.Code)

	if expectedErrors != nil {
		var responseBody struct {
			Errors []DTOs.FieldErrorDTO `json:"errors"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, expectedErrors, responseBody.Errors)
	}
}

func drop_all_tables() {
	// Supprimer toutes les tables de façon vraiment dégueulasse
	database.DB.Exec("SET FOREIGN_KEY_CHECKS = 0")
	database.DB.Exec("DROP TABLE IF EXISTS schema_migrations")
	database.DB.Exec("DROP TABLE IF EXISTS users")
	database.DB.Exec("DROP TABLE IF EXISTS projects")
	database.DB.Exec("DROP TABLE IF EXISTS categories")
	database.DB.Exec("DROP TABLE IF EXISTS activities")
	database.DB.Exec("DROP TABLE IF EXISTS co_managers")
	database.DB.Exec("SET FOREIGN_KEY_CHECKS = 1")
}
