### Documentation : Automatisation des Backups MariaDB (DB `llio`)

Cette procédure détaille la mise en place d'une sauvegarde automatisée, compressée et rotative (rétention 365 jours) pour la base de données `llio`.

#### 1. Script de Sauvegarde

**Fichier :** `/usr/local/bin/mariadb-backup.sh`
Important : Mettre le fichier dans le path   /usr/local/bin/mariadb-backup.sh

Ce script utilise `set -o pipefail` pour garantir que l'échec du dump stoppe l'exécution, même si la compression réussit. 
Il inclut également un horodatage dans les logs pour faciliter le débogage via cron.
```sh
#!/bin/bash
set -e
set -o pipefail

# ===== CONFIGURATION =====
DB_NAME="llio"
BACKUP_DIR="/var/backups/mariadb"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
FILE_NAME="${BACKUP_DIR}/${DB_NAME}_${TIMESTAMP}.sql.gz"
RETENTION_DAYS=365
LOG_TAG="[mariadb-backup]"

# Fonction de log simple
log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') $LOG_TAG $1"
}

# ===== PRÉPARATION =====
# Création du dossier avec permissions restreintes si inexistant
if [ ! -d "$BACKUP_DIR" ]; then
    mkdir -p "$BACKUP_DIR"
    chmod 700 "$BACKUP_DIR"
fi

log "Démarrage du backup pour $DB_NAME"

# ===== BACKUP =====
# --single-transaction : Consistance sans verrouiller les tables (InnoDB)
# --routines --triggers --events : Sauvegarde complète de la logique
mariadb-dump \
  --defaults-extra-file=/root/.my.cnf \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  --databases "$DB_NAME" \
| gzip > "$FILE_NAME"

log "Succès : $FILE_NAME créé."

# ===== NETTOYAGE =====
# Suppression des fichiers vieux de plus de X jours
log "Nettoyage des archives de plus de $RETENTION_DAYS jours..."
find "$BACKUP_DIR" -type f -name "${DB_NAME}_*.sql.gz" -mtime +"$RETENTION_DAYS" -print -delete | while read -r deleted_file; do
    log "Supprimé : $deleted_file"
done

log "Terminé."

```


# Installation et Permissions

- Créer/Éditer le fichier :
```sh
sudo nano /root/.my.cnf
```

- Ajouter les identifiants :
```sh
[client]
user=root
password=VOTRE_MOT_DE_PASSE_ICI
```

- Ouvrir l'éditeur cron :
```sh
sudo crontab -e
```

- Ajouter la ligne suivante :
```bash
0 0 * * * /usr/local/bin/mariadb-backup.sh >> /var/log/mariadb-backup.log 2>&1
```
Note : 0 0 : indique 00h, ex 15 2 = 2h15min du matin
- Lancer les commandes une par une.
```bash
# Rendre le script exécutable
sudo chmod +x /usr/local/bin/mariadb-backup.sh

# Préparer le fichier de log avec les bons droits
sudo touch /var/log/mariadb-backup.log
sudo chmod 640 /var/log/mariadb-backup.log
# (Optionnel : changer le propriétaire si un user spécifique lance le cron, ex: root)
sudo chown root:root /var/log/mariadb-backup.log

sudo chmod 600 /root/.my.cnf
sudo chmod -R 755 /var/backups/mariadb


```


