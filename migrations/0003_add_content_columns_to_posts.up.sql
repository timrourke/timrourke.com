ALTER TABLE `posts`
ADD COLUMN `title` VARCHAR(250) AFTER `updated_at`,
ADD COLUMN `excerpt` TEXT AFTER `title`,
ADD COLUMN `content` LONGTEXT AFTER `excerpt`,
ADD COLUMN `permalink` VARCHAR(250) AFTER `content`,
ADD INDEX `permalink` (`permalink`);
