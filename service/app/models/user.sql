-- 有缓存 goctl model mysql ddl -src user.sql -dir . -c
-- 无缓存 goctl model mysql ddl -src user.sql -dir . 
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '用户编号',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `open_id`     varchar(200) NOT NULL DEFAULT '' COMMENT '微信用户id',
    `mobile`      char(11) NOT NULL DEFAULT '' COMMENT '电话',
    `username`    varchar(50) NOT NULL DEFAULT '' COMMENT '用户名称',
    `password`    varchar(50) NOT NULL DEFAULT '' COMMENT '用户密码',
    `sex`         int NOT NULL DEFAULT '0' COMMENT '性别 0:男 1:女',
    `email`       varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
    `role`        int NOT NULL DEFAULT '0' COMMENT '身份 0-申请队员 1-岗前培训 2-见习队员 3-正式队员 4-督导老师 5-区域负责人 6-组委会 7-主任',
    `avatar`      varchar(255) NOT NULL DEFAULT '0' COMMENT '头像',
    `address`     varchar(255) NOT NULL DEFAULT '0' COMMENT '地址',
    `birthday`    varchar(20) NOT NULL DEFAULT '0000-00-00' COMMENT '生日 xxxx-xx-xx',
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_mobile` (`mobile`) 
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';

DROP TABLE IF EXISTS `users_auth`;
CREATE TABLE `users_auth`(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `user_id`     int unsigned NOT NULL DEFAULT '0',
    `auth_key`    varchar(64) NOT NULL DEFAULT '' COMMENT '平台唯一id',
    `auth_type`   varchar(12) NOT NULL DEFAULT '' COMMENT '平台类型',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_type_key` (`auth_type`,`auth_key`) USING BTREE, -- 复合索引，索引的存储类型为BTREE
    UNIQUE KEY `idx_userId_key` (`user_id`,`auth_type`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户授权表';
