-- 有缓存 goctl model mysql ddl -src role.sql -dir . -c
-- 无缓存 goctl model mysql ddl -src role.sql -dir . 
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles`(
    `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '身份编号',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `user_id`     INT UNSIGNED NOT NULL DEFAULT '0',
    `user_role`   INT NOT NULL DEFAULT '0' COMMENT '身份 0-申请队员 1-岗前培训 2-见习队员 3-正式队员 4-督导老师 5-区域负责人 6-组委会 7-主任',
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_user_id` (`user_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='身份表';

DROP TABLE IF EXISTS `roles_change`;
CREATE TABLE `roles_change`(
    `id`          INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `user_id`     INT UNSIGNED NOT NULL DEFAULT '0' COMMENT '身份变动用户',
    `operator_id` INT UNSIGNED NOT NULL DEFAULT '0' COMMENT '操作人',
    `new_role`    INT UNSIGNED NOT NULL DEFAULT '0' COMMENT '新身份',
    `old_role`    INT UNSIGNED NOT NULL DEFAULT '0' COMMENT '旧身份',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='身份变动表';