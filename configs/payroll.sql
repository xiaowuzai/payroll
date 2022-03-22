
CREATE DATABASE IF NOT EXISTS payroll DEFAULT CHARSET utf8 COLLATE utf8_general_ci;
USE payroll;

DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `description` varchar(255) NOT NULL,
    `roles` varchar(255) NOT NULL,
    `createdAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp,
    `updatedAt` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    PRIMARY  KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(45) NOT NULL,
    `telephone` varchar(11) NOT NULL,
    `email` varchar(45),
    `role` int NOT NULL,
    `password` varchar(255) NOT NULL,
    `salt` varchar(255) NOT NULL,
    `status` tinyint NOT NULL DEFAULT 0,
    `created` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp,
    `updated` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    PRIMARY  KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `employee`;
CREATE TABLE `employee` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `number` int unsigned NOT NULL,
    `id_card` varchar(18) NOT NULL,
    `telephone` varchar(11) NOT NULL,
    `offer_time` timestamp NOT NULL,
    `retire_time` timestamp,
    `duty` varchar(24),
    `post` varchar(24),
    `level` varchar(24),
    `base_salary` int,
    `identity` tinyint(3),
    `deleted` tinyint,
    `created` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp,
    `updated` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE  IF EXISTS `organization`;
CREATE TABLE `organization` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `parent_id` int unsigned NOT NULL,
    `name` varchar(255) NOT NULL,
    `type` tinyint,
    `path` varchar(255),
    `created` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp,
    `updated` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `organization` VALUES (1,0,'机关总部',0,'.1',NOW(),NOW());

DROP TABLE  IF EXISTS `organization_salary`;
CREATE TABLE `organization_salary` (
    `id` int unsigned NOT NULL AUTO_INCREMENT,
    `salary_type` tinyint,
    `employee_type` tinyint,
    `created` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp,
    `updated` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;