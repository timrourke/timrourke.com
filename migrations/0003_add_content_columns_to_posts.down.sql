ALTER TABLE `posts`
DROP COLUMN `title`,
DROP COLUMN `excerpt`,
DROP COLUMN `content`,
DROP COLUMN `permalink`,
DROP INDEX `permalink`;
