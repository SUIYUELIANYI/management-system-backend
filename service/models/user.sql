-- 有缓存 goctl model mysql ddl -src user.sql -dir . -c
-- 无缓存 goctl model mysql ddl -src user.sql -dir . 
DROP TABLE IF EXISTS `user`; -- 用户表
CREATE TABLE `user`(
    `id`          int unsigned NOT NULL AUTO_INCREMENT COMMENT '用户编号',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`   tinyint NOT NULL DEFAULT '0',
    `mobile`      char(11) NOT NULL DEFAULT '' COMMENT '电话', -- 微信登录和系统注册都需要电话
    `username`    varchar(50) NOT NULL DEFAULT '' COMMENT '用户名',
    `password`    varchar(50) NOT NULL DEFAULT '' COMMENT '用户密码',
    `sex`         int NOT NULL DEFAULT '0' COMMENT '性别 0-男 1-女', -- 默认为男
    `email`       varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
    `avatar`      varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
    `address`     varchar(255) NOT NULL DEFAULT '' COMMENT '地址',
    `birthday`    varchar(20) NOT NULL DEFAULT '' COMMENT '生日 xxxx-xx-xx',
    `role`        int NOT NULL DEFAULT '0' COMMENT '身份代号 -1-待处理人员 0-非在册队员 1-申请队员 2-岗前培训 3-见习队员 4-正式队员 5-督导老师 6-树洞之友 40-普通队员 41-核心队员 42-区域负责人 43-组委会成员 44-组委会主任',
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_mobile` (`mobile`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';

DROP TABLE IF EXISTS `user_auth`; -- 用户授权表（分为平台内部和微信小程序，平台内部的唯一id就是手机号，微信小程序的唯一id就是open_id）
CREATE TABLE `user_auth`(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`   tinyint NOT NULL DEFAULT '0',
    `user_id`     int unsigned NOT NULL DEFAULT '0',
    `auth_key`    varchar(64) NOT NULL DEFAULT '' COMMENT '平台唯一id', -- 电话/openid
    `auth_type`   varchar(12) NOT NULL DEFAULT '' COMMENT '平台类型', -- 平台内部/微信小程序
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_type_key` (`auth_type`,`auth_key`) USING BTREE, -- 复合索引，索引的存储类型为BTREE
    UNIQUE KEY `idx_userId_key` (`user_id`,`auth_type`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户授权表';


DROP TABLE IF EXISTS `pending_personnel`; -- 待处理人员表（申请表三次不通过将该用户移入此表）
CREATE TABLE `pending_personnel`(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`   tinyint NOT NULL DEFAULT '0',
    `user_id`     int unsigned NOT NULL DEFAULT '0' COMMENT '待处理人员id',
    `reason`      varchar(255) NOT NULL DEFAULT '' COMMENT '原因',
    `operate_id`  int unsigned NOT NULL DEFAULT '0' COMMENT '操作人id',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_userId` (`user_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='待处理人员表';

