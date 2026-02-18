ALTER TABLE `users`
ADD COLUMN `time_bank_start_date` DATE NULL COMMENT 'Date de début de calcul de la banque d''heures',
ADD COLUMN `time_bank_hours_per_week` FLOAT NULL COMMENT 'Nombre d''heures attendues par semaine',
ADD COLUMN `time_bank_balance_offset` FLOAT DEFAULT 0 COMMENT 'Ajustement manuel du solde';