CREATE TABLE IF NOT EXISTS `posts` (
	`id` INT NOT NULL AUTO_INCREMENT,
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, 
	`updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`user_id` INT NOT NULL,
	INDEX `user_id` (`user_id`),
	FOREIGN KEY (`user_id`)
		REFERENCES users(`id`),
	PRIMARY KEY (`id`)
) ENGINE=InnoDB;
