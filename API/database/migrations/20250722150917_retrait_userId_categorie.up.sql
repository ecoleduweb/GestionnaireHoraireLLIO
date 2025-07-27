-- Migration Up
ALTER TABLE `categories` 
DROP FOREIGN KEY `fk_categories_users`,
DROP COLUMN `user_id`;

