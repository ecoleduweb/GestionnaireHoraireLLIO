ALTER TABLE `users`
DROP COLUMN IF EXISTS `time_bank_start_date`,
DROP COLUMN IF EXISTS `time_bank_hours_per_week`,
DROP COLUMN IF EXISTS `time_bank_balance_offset`;