# Application WEB LLIO

Une application WEB de gestion de calendrier développée avec SvelteKit et TypeScript, intégrant FullCalendar..

---

## Technologies utilisées
- **Framework** : SvelteKit
- **Langage** : TypeScript
- **Bibliothèque calendrier** : FullCalendar
- **Bundler** : Vite
- **CSS** : Tailwindcss

## Structure du projet

```plaintext
LLIO2025/
├── .svelte-kit/      # Fichiers générés par SvelteKit
├── .vscode/          # Configuration VS Code
├── node_modules/     # Dépendances du projet
├── src/
│   ├── Components/   # Composants réutilisables
│   │   └── Calendar/ # Composants liés au calendrier
│   ├── forms/        # Formulaires et validations
│   │   └── activity/     # Formulaires liés aux tâches
│   ├── lib/          # Bibliothèques et utilitaires partagés
│   ├── Models/       # Types et interfaces TypeScript
│   ├── routes/       # Pages et routage de l'application
│   │   ├── login/    # Gestion de l'authentification
│   ├── services/     # Services pour la logique métier et API
│   ├── style/        # Fichiers CSS et styles globaux
│   ├── ts/           # Configuration et types TypeScript
│   ├── utils/        # Fonctions utilitaires (dates, formatage, etc.)
└── static/           # Ressources statiques (images, fonts, etc.)
```

## Prérequis
- Node.js (version 16 ou supérieur)
- npm (gestionnaire de paquets Node.js)

## Installation
1. Cloner le projet
```bash
git clone [URL_DU_REPO]
cd LLIO2025
```
2. Installer les dépendances
```bash
npm install
```
3. Installer les dépendances FullCalendar
```bash
npm install @fullcalendar/core @fullcalendar/daygrid @fullcalendar/timegrid @fullcalendar/interaction
```
4. Créer un .env à la racine du projet avec cette structure
````
VITE_API_BASE_URL=http://localhost:8080
````

## Développement
Lancer le serveur de développement :
```bash
npm run dev
```
L'application sera disponible à l'adresse afficher dans le terminal
