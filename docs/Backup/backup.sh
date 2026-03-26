#!/bin/bash
set -e
set -o pipefail

# 1. Récupération du chemin du dossier où se trouve ce script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/db.config"

# 2. Chargement de la configuration
if [ -f "$CONFIG_FILE" ]; then
    source "$CONFIG_FILE"
else
    echo "ERREUR : Fichier de configuration introuvable : $CONFIG_FILE"
    echo "Veuillez copier db.config.example vers db.config et le configurer."
    exit 1
fi

TIMESTAMP=$(date +%Y%m%d_%H%M%S)
FILE_NAME="${BACKUP_DIR}/${DB_NAME}_${TIMESTAMP}.sql.gz"

# 3. Vérification/Création du dossier de destination
if [ ! -d "$BACKUP_DIR" ]; then
    mkdir -p "$BACKUP_DIR"
    chmod 700 "$BACKUP_DIR"
fi

echo "[$(date +'%Y-%m-%d %H:%M:%S')] Démarrage backup : $DB_NAME"

# 4. Dump (Mot de passe passé via variable temporaire sécurisée)
export MYSQL_PWD="$DB_PASS"

mariadb-dump \
  -u "$DB_USER" \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  --databases "$DB_NAME" \
| gzip > "$FILE_NAME"

unset MYSQL_PWD

echo "[$(date +'%Y-%m-%d %H:%M:%S')] Succès : $FILE_NAME"

# 5. Rotation (Suppression des vieux backups)
find "$BACKUP_DIR" -type f -name "${DB_NAME}_*.sql.gz" -mtime +"$RETENTION_DAYS" -delete

echo "[$(date +'%Y-%m-%d %H:%M:%S')] Terminé."