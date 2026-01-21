# API LLIO

> Une API développée en Go avec le framework Gin-Gonic.

## Structure du projet

```plaintext
API/
├── controllers/   # Logique des routes
├── database/      # Connexion et configuration de la base de données
|__ auth/          # Configuration de l'authentification
├── main.go        # Point d'entrée de l'application
├── middleware/    # Middlewares
├── models/        # Modèles de données
│   ├── DAOs/      # Objets d'accès aux données
│   └── DTOs/      # Objets de transfert de données
├── repositories/  # Gestion des accès aux données
├── routes/        # Définition des routes
├── services/      # Logique métier
├── tests/         # Tests unitaires et d'intégration
└── useful/        # Fonctions utilitaires
```

# Prérequis

1. Installation de la derniere version stable de go sur golang.org
   (go version go1.23.5 windows/amd64)

2. Une instance de MySQL Client (mariaDB)

# Variables d'environnement

Ce projet utilise les variables d'environement.

1. Se créer un fichier .env a la racine du dossier API

2. Copier le contenu du fichier `.env.template` et le coller dans le `.env`

3. Ajuster les valeurs selon votre réalité

# Dépendances

Installation de toutes les dépendances avec la commande suivante

```bash
go mod tidy
```

# Migrations et mise en place de la base de données

1. Créer une base de données vide nommée llio et une autre nommée llio_test sur l'instance de MariaDb

2. À la racine du projet API, exécuter la commande

```bash
go run main.go migrate up
```

Cette commande permet d'exécuter les nouvelles migrations en se basant sur la table schema_migrations.

## Pour créer une nouvelle migration

1. Exécuter la commande

```bash
go run main.go migrate create un_nom_significatif_pour_ta_migration
```

Cela va ajouter deux fichiers dans le dossier database/migrations.

2. Complète les fichiers de migration. L'ia peut t'aider à le faire, mais assures-toi de bien valider ce qui est produit! Le ficher up applique le changement, le fichier down revient en arrière.

Test ensuite ta migration. Si elle échoue, la table schema_migration aura ton numéro de migration et la valeur 1 à isDirty qui montre que ta migration a échouée.

## Pour annuler une migration

La commande

```bash
go run main.go migrate down
```

annule la dernière migration qui a été exécutée.

## Si jamais ma migration a tout cassé?

Le champ isDirty dans schema_migration sera à 1. Le changer pour 0 dans la base de données et re rouler la migration!

# Test

## Installation de la librairie de test pour du golang : testify

Installation de la librairie de base:
go get -u github.com/stretchr/testify

Pour l'utilisation graphique des tests avec l'extension Testing, il faut déactiver le cache des tests antérieurs (mesure de précaution)
Étapes :

1. Ouvrez les Paramètres de VS Code (via Ctrl+, ou Cmd+, sur macOS).

2. Recherchez go.testFlags.

3. Ajoutez -count=1 à la liste des flags de test:
   `"go.testFlags": ["-count=1"]`

Pour excécution des tests en ligne de commande:

```bash
go test -v -count=1 ./...
```

# Démarrage du serveur

Pour démarrer le serveur exécutez la ligne suivante :

```bash
go run main.go
```

# Mise en place du compte administateur

Par defaut, le compte administrateur est un simple utilisteur.
Il faut donc modifier son role manuellement dans la base de donnée

```bash
UPDATE users SET role=2 WHERE id=1;
```
