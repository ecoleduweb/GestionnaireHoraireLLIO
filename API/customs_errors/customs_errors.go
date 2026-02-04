package customs_errors

import (
	"errors"
)

// Types d'erreurs personnalisées pour la base de données
var (
	ErrDuplicateEntry       = errors.New("duplication")
	ErrNotFound             = errors.New("ressource introuvable")
	ErrDatabase             = errors.New("erreur de la BD")
	ErrUserHasActivities    = errors.New("l'utilisateur a des activités associées, suppression impossible")
	ErrUserHasProjects      = errors.New("l'utilisateur a des projets associées, suppression impossible")
	ErrProjectHasActivities = errors.New("Le projet a des activités associées, suppression impossible")
)
