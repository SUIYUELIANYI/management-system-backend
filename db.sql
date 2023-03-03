DROP DATABASE IF EXISTS `managementsystem`;
CREATE DATABASE `managementsystem`;
USE `managementsystem`;

-- 用户
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
    `realname`    varchar(255) NOT NULL DEFAULT '' COMMENT '真实姓名',
    `sex`         tinyint(1) NOT NULL DEFAULT '0' COMMENT '性别 0:男 1:女',
    `email`       varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
    `status`      int NOT NULL DEFAULT '0' COMMENT '用户身份',
    `avatar`      varchar(255) NOT NULL DEFAULT '0' COMMENT '头像',
    `address`     varchar(255) NOT NULL DEFAULT '0' COMMENT '地址',
    `birthday`    varchar(255) NOT NULL DEFAULT '0' COMMENT '生日',
    `info`        varchar(255) NOT NULL DEFAULT '0' COMMENT '信息',
    PRIMARY KEY(`id`),
    UNIQUE KEY (`mobile`) 
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO users(mobile,realname,username,password) VALUES("18258456918","方瑜诚","千人","123456");