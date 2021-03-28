CREATE DATABASE todo;

USE todo;

DROP TABLE IF EXISTS `tasks`;
CREATE TABLE `tasks` (
                         `id` int(11) NOT NULL AUTO_INCREMENT,
                         `name` text NOT NULL,
                         `created_at` timestamp NOT NULL,
                         `updated_at` timestamp NOT NULL,
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;