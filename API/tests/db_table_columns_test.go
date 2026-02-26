package tests

import (
	"llio-api/database"
	"testing"
)

func TestUserIdColumnCount(t *testing.T) {
	// Requête pour compter les tables ayant une colonne nommée 'userid'
	var count int64

	query := `
	SELECT COUNT(*) 
	FROM information_schema.columns 
	WHERE column_name = 'user_id' 
 	AND table_schema = DATABASE()
	`

	err := database.DB.Raw(query).Scan(&count).Error
	if err != nil {
		t.Fatalf("Erreur lors de la vérification des colonnes user_id: %v", err)
	}

	// Vérifier qu'il y a exactement 2 tables avec une colonne user_id
	if count != 2 {
		t.Errorf("ERREUR CRITIQUE: Il y a %d tables avec une colonne 'user_id', mais il devrait y en avoir exactement 2.", count)
		t.Error("COMMENTAIRE IMPORTANT: Si ce test échoue, il faut ajouter un check dans la fonction DeleteUserById")
		t.Error("pour prévenir la suppression involontaire d'un utilisateur. Vérifiez que toutes les")
		t.Error("tables avec des références à 'user_id' sont correctement gérées lors de la suppression.")
		t.FailNow()
	}

	t.Logf("✓ Vérification réussie: %d tables contiennent une colonne 'user_id'", count)
}
