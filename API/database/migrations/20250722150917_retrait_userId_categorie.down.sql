-- Migration Down
ALTER TABLE `categories` 
ADD COLUMN `user_id` bigint(20) unsigned NOT NULL,
ADD CONSTRAINT `fk_categories_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
