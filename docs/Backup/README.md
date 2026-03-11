# Système de Backup Automatisé

## Installation (À faire une seule fois sur le serveur)

1. **Configurer les identifiants :**
   Dans ce dossier `backup/`, créez le fichier de configuration : db.config basé sur le fichier db.config.example.

   Remplir : DB_USER, DB_PASS, DB_NAME

- exemple du model 
```bash
DB_USER="root"
DB_PASS="remplacer_par_password"
DB_NAME="remplacer_par_nom_bd"
BACKUP_DIR="/var/backups/mariadb"
RETENTION_DAYS=365
```
2. **Sécuriser et activer le script:**
```bash
chmod 600 db.config       # Lecture seule pour le propriétaire
chmod +x backup.sh        # Rendre exécutable
```
3. **Configurer le Cron**
ance cette commande dans le dossier **backup/**
 `echo "0 0 * * * $(pwd)/backup.sh >> /var/log/mariadb-backup.log 2>&1"

Copie le résultat qui s'affiche (c'est ta ligne crontab) et colle-le dans :
` sudo crontab -e
