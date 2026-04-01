#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/db.config"

if [ -f "$CONFIG_FILE" ]; then
    source "$CONFIG_FILE"
else
    echo "ERREUR: Config introuvable."
    exit 1
fi

if [ -z "$1" ]; then
    echo " Usage: ./restore.sh <CHEMIN_DU_FICHIER_SQL_GZ>"
    echo "Exemple: ./restore.sh $BACKUP_DIR/llio_20240101.sql.gz"
    exit 1
fi

BACKUP_FILE="$1"

if [ ! -f "$BACKUP_FILE" ]; then
    echo " Erreur : Le fichier $BACKUP_FILE n'existe pas."
    exit 1
fi

echo " ATTENTION : La base de données '$DB_NAME' sera écrasée !"
read -p "Confirmer la restauration ? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Annulé."
    exit 1
fi

export MYSQL_PWD="$DB_PASS"

echo " Restauration en cours..."
zcat "$BACKUP_FILE" | mariadb -u "$DB_USER" "$DB_NAME"

unset MYSQL_PWD

echo " Base de données restaurée avec succès."