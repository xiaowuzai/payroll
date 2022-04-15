CREATE DATABASE IF NOT EXISTS payroll DEFAULT CHARSET utf8 COLLATE utf8_general_ci;
USE payroll;

INSERT INTO role VALUES('4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','系统管理员', '系统管理员',0, NOW(), NOW());

INSERT INTO menu VALUES
    ('3f946719-679d-4d0a-b265-a33e16cdd7be','salary-management', '工资管理', NOW(), NOW()),
    ('4b1c567e-c856-4053-afb0-f7c3094f174f','salary-preview', '当月工资预览', NOW(), NOW()),
    ('4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','salary-approve', '当月工资审批', NOW(), NOW()),
    ('559a1bf5-03dc-4dd3-8b79-ab816b8c30ac','salary-query', '当月工资统计', NOW(), NOW()),
    ('57344b90-e8ab-4ef6-8163-67c105abf6ef','salary-carry-forward-record', '结转记录', NOW(), NOW()),
    ('65ca014a-bf3a-42a5-8e8a-00fd52b22699','salary-note-history', '备忘记录', NOW(), NOW()),
    ('67b18bdb-dadf-4a0d-b6fb-95cfe360e21a','salary-statistics', '工资统计', NOW(), NOW()),
    ('6d2b6450-6cd3-4956-ab55-c18c400bb036','cost-statistics', '成本统计', NOW(), NOW()),
    ('7497609c-8fba-4c33-8bdd-bbc443865b73','project-cost-statistic', '项目费用统计', NOW(), NOW()),
    ('7fad0844-4563-408c-a269-e3a2d401ae88','system-config', '系统设置', NOW(), NOW()),
    ('8668919e-1112-49cc-8ae0-74808cb9e9cc','user-management', '用户管理', NOW(), NOW()),
    ('8d7fe0fb-4b98-450b-8e48-aac8a5809370','role-management', '角色管理', NOW(), NOW()),
    ('76b18bdb-dadf-4a0d-b6fb-95cfe360e21a','organization', '组织机构', NOW(), NOW()),
    ('6d7fe0fb-4b98-450b-8e48-aac8a5809370','staff-management', '员工管理', NOW(), NOW());

INSERT INTO role_menu VALUES
    ('de504569-42f5-4a90-9fc4-2481f05cdd89','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','3f946719-679d-4d0a-b265-a33e16cdd7be',NOW(), NOW()),
    ('73853828-d922-4c6a-b64c-6acc03df6831','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','4b1c567e-c856-4053-afb0-f7c3094f174f',NOW(), NOW()),
    ('2a9ccd2b-192a-4865-a3e0-b50656441305','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81',NOW(), NOW()),
    ('4e821e0b-297c-4578-b1ee-379a84091192','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','559a1bf5-03dc-4dd3-8b79-ab816b8c30ac',NOW(), NOW()),
    ('e7e255c3-e608-49f8-9db7-d66b72e1382c','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','57344b90-e8ab-4ef6-8163-67c105abf6ef',NOW(), NOW()),
    ('e4af79ee-2ae8-49c1-8a1f-e8e4c7b5e2e3','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','65ca014a-bf3a-42a5-8e8a-00fd52b22699',NOW(), NOW()),
    ('e3488d80-b05c-442d-8feb-94c2eadebc46','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','67b18bdb-dadf-4a0d-b6fb-95cfe360e21a',NOW(), NOW()),
    ('d86b6b91-60bf-42b6-a846-b1831baa0344','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','6d2b6450-6cd3-4956-ab55-c18c400bb036',NOW(), NOW()),
    ('7ccee999-37ca-49b9-881b-a6c95809a6a1','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','7497609c-8fba-4c33-8bdd-bbc443865b73',NOW(), NOW()),
    ('1b73b515-16e0-4841-8397-5442552ec708','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','7fad0844-4563-408c-a269-e3a2d401ae88',NOW(), NOW()),
    ('26566951-6b8d-4bed-a3d2-d12f29c90fa2','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','8668919e-1112-49cc-8ae0-74808cb9e9cc',NOW(), NOW()),
    ('16f5f7e7-e021-406c-9a8f-aeb2f136d518','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','8d7fe0fb-4b98-450b-8e48-aac8a5809370',NOW(), NOW()),
    ('40c7b944-8c0a-421a-b3f4-b52c996b5524','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','76b18bdb-dadf-4a0d-b6fb-95cfe360e21a',NOW(), NOW()),
    ('40c7b944-8c0a-421a-b3f4-b52c996b1234','4fba1999-a7f9-4b34-b82a-b3f34cfe4d81','6d7fe0fb-4b98-450b-8e48-aac8a5809370',NOW(), NOW());

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

INSERT INTO `organization` VALUES ('16f5f7e7-e021-406c-9a8f-aeb2f136d518','','机关总部','.16f5f7e7-e021-406c-9a8f-aeb2f136d518', 0, 0, 0, 0, NOW(),NOW());
