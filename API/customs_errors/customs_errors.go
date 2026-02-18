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
	ErrUserForbidden        = errors.New("accès refusé")
	ErrUserIsManager        = errors.New("l'utilisateur est le chargé du projet")
	ErrUserRoleBelowManager = errors.New("l'utilisateur doit avoir le rôle de chargé de projet ou supérieur pour devenir co-chargé de projet")
	ErrProjectHasActivities = errors.New("Le projet a des activités associées, suppression impossible")
	ErrProjectNotFound      = errors.New("Projet non trouvé - Projet inexistant")
)
