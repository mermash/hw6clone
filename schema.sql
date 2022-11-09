SET NAMES utf8mb4;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

CREATE DATABASE IF NOT EXISTS`redditclone`;

USE `redditclone`;

DROP TABLE IF EXISTS `redditclone`.`category`;
CREATE TABLE `category` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `redditclone`.`category` (`name`) VALUES 
('music'),
('funny'),
('videos'),
('programming'),
('news'),
('fashion');

DROP TABLE IF EXISTS `redditclone`.`post`;
CREATE TABLE `redditclone`.`post` (
  `id` varchar(36) NOT NULL,
  `title` varchar(255) NOT NULL,
  `type` ENUM('text', 'link') DEFAULT NULL,
  `description` text NOT NULL,
  `score` int(11) DEFAULT NULL,
  `user_id` varchar(36) NOT NULL,
  `category_id` int(11) NOT NULL, 
  `created` varchar(255) DEFAULT NULL,
   UNIQUE KEY `id` (`id`),
   KEY `user_id` (`user_id`),
   CONSTRAINT `posts_user_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `redditclone`.`comment`;
CREATE TABLE `redditclone`.`comment` (
  `id` varchar(36) NOT NULL,
  `post_id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `body` text NOT NULL,
  `created` varchar(255) DEFAULT NULL,
   UNIQUE KEY `id` (`id`),
   KEY `user_id` (`user_id`),
   CONSTRAINT `user_comments_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


DROP TABLE IF EXISTS `redditclone`.`user`;
CREATE TABLE `redditclone`.`user` (
  `id` varchar(36) NOT NULL,
  `login` varchar(255) NOT NULL,
  `password` varchar(60) NOT NULL,
  `created` varchar(255) DEFAULT NULL,
   UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `redditclone`.`user` (`id`, `login`, `password`, `created`) VALUES 
('34420d9d-91c0-4c6f-96fa-e4346eb9361c', 'test', 'test', '2022-11-02 15:24:00');

DROP TABLE IF EXISTS `redditclone`.`vote`;
CREATE TABLE `redditclone`.`vote` (
    `post_id` varchar(36) NOT NULL,
    `user_id` varchar(36) NOT NULL,
    `vote` int(11) NOT NULL,
    UNIQUE KEY `post_id_user_id` (`post_id`, `user_id`),
    KEY `post_id` (`post_id`),
    KEY `user_id` (`user_id`),
    CONSTRAINT `posts_votes_ibfk_1` FOREIGN KEY (`post_id`) REFERENCES `post`(`id`),
    CONSTRAINT `users_votes_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


DROP TABLE IF EXISTS `redditclone`.`sessions`;
CREATE TABLE `redditclone`.`sessions` (
    `id` varchar(36) NOT NULL,
    `user_id` varchar(36) NOT NULL,
    UNIQUE KEY `id` (`id`),
    KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;