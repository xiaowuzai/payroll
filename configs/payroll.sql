DROP DATABASE IF EXISTS payroll;

CREATE DATABASE IF NOT EXISTS payroll DEFAULT CHARSET utf8 COLLATE utf8_general_ci;
USE payroll;

DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
    `id` varchar(36) NOT NULL,
    `name` varchar(30) UNIQUE NOT NULL,
    `description` varchar(255) NOT NULL,
    `roles` varchar(300) UNIQUE NOT NULL,
    `deleted` tinyint,
    `created` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp,
    `updated` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    PRIMARY  KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
INSERT INTO role VALUES('4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','系统管理员', '系统管理员',
                        'salary-management,salary-preview,salary-approve,salary-query,salary-carry-forward,salary-note-history,salary-statistics,cost-statistics,project-cost-statistic,system-config,user-management,role-management,organization,staff-management',
                        0, NOW(), NOW());

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `id` varchar(36) NOT NULL,
    `username` varchar(45) UNIQUE NOT NULL,
    `account_name` varchar(20) UNIQUE NOT NULL,
    `email` varchar(45) UNIQUE,
    `role_id` varchar(36) NOT NULL,
    `password` blob NOT NULL,
    `salt` blob NOT NULL,
    `status` tinyint NOT NULL DEFAULT 0,
    `deleted` tinyint,
    `created` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp,
    `updated` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    PRIMARY  KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `employee`;
CREATE TABLE `employee` (
    `id` varchar(36) NOT NULL,
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
    `id` varchar(36) NOT NULL,
    `parent_id` int unsigned NOT NULL,
    `name` varchar(255) NOT NULL,
    `path` varchar(255),
    `type` tinyint,
    `salary_type` tinyint,
    `employee_type` tinyint,
    `deleted` tinyint,
    `created` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp,
    `updated` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
    PRIMARY KEY(`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `organization` VALUES ('16f5f7e7-e021-406c-9a8f-aeb2f136d518','','机关总部','.16f5f7e7-e021-406c-9a8f-aeb2f136d518', 0, 0, 0, 0, NOW(),NOW());
