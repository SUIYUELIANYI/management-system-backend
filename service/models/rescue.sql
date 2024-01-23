-- 有缓存 goctl model mysql ddl -src rescue.sql -dir . -c
-- 无缓存 goctl model mysql ddl -src rescue.sql -dir .
-- 救援评价对救援对象和救援信息都有，每个老师都可以对救援信息评价，但对救援对象的评价所有老师只能写一次。
-- 本来救援对象表想把weibo_address当唯一索引的，但是考虑到救援结束后救援对象仍有可能会求助，所以不设置唯一索引。后来我又想了想，如果再发生求助，直接修改救援状态，然后继续就是了。
DROP TABLE IF EXISTS `rescue_target`;
CREATE TABLE `rescue_target`(
    `id`                         int unsigned NOT NULL AUTO_INCREMENT COMMENT '救援对象编号',
    `create_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`                  tinyint NOT NULL DEFAULT '0',
    `weibo_address`              varchar(200) DEFAULT '' NOT NULL COMMENT '微博地址',
    `status`                     INT NOT NULL DEFAULT '0' COMMENT '救援状态 0-待救援 1-救援中 2-已救援',
    `start_time`                 datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '救援起始时间',
    `end_time`                   datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '救援结束时间',
    `rescue_teacher1_id`         int unsigned NOT NULL DEFAULT '0' COMMENT '救助老师1id', -- 救援老师可以为3-见习队员 4-正式队员（40-普通队员 41-核心队员 42-区域负责人 43-组委会成员 44-组委会主任）
    `rescue_teacher2_id`         int unsigned NOT NULL DEFAULT '0' COMMENT '救助老师2id',
    `rescue_teacher3_id`         int unsigned NOT NULL DEFAULT '0' COMMENT '救助老师3id',
    `description`                varchar(10000) NOT NULL DEFAULT '' COMMENT '救援过程描述', -- 由救援老师填写，可接续
    `evaluation`                 varchar(5000) NOT NULL DEFAULT '' COMMENT '最终评价',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='救援对象表';

DROP TABLE IF EXISTS `rescue_info`;
CREATE TABLE `rescue_info`(
    `id`                         int unsigned NOT NULL AUTO_INCREMENT COMMENT '救援信息编号',
    `create_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`                  tinyint NOT NULL DEFAULT '0',
    `rescue_target_id`           int unsigned NOT NULL COMMENT '救援对象编号',
    `release_time`               varchar(100) DEFAULT '' NOT NULL COMMENT '发布时间',
    `weibo_account`              varchar(200) DEFAULT '' NOT NULL COMMENT '微博账号',
    `weibo_address`              varchar(200) DEFAULT '' NOT NULL COMMENT '微博地址',
    `nickname`                   varchar(200) DEFAULT '' NOT NULL COMMENT '昵称',
    `risk_level`                 varchar(200) DEFAULT '' COMMENT '危险级别',
    `area`                       varchar(500) DEFAULT '' COMMENT '所在城市',
    `sex`                        int NOT NULL DEFAULT '0' COMMENT '性别 0-男 1-女',
    `birthday`                   varchar(200)  DEFAULT '' COMMENT '生日',
    `brief_introduction`         varchar(4000) DEFAULT '' COMMENT '简介',
    `text`                       text  COMMENT '信息原文',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='救援信息表';

DROP TABLE IF EXISTS `rescue_process`;
CREATE TABLE `rescue_process`(
    `id`                         int unsigned NOT NULL AUTO_INCREMENT COMMENT '救援过程编号',
    `create_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`                  tinyint NOT NULL DEFAULT '0',
    `rescue_teacher_id`          int unsigned NOT NULL DEFAULT '0' COMMENT '救援老师编号',
    `rescue_info_id`             int unsigned NOT NULL DEFAULT '0' COMMENT '救援信息编号',
    `start_time`                 varchar(100) NOT NULL DEFAULT '' COMMENT '救援起始时间', -- 0000-00-00 00:00
    `end_time`                   varchar(100) NOT NULL DEFAULT '' COMMENT '救援结束时间',
    `duration`                   varchar(100) NOT NULL DEFAULT '0h0m' COMMENT '救援时长',
    `evaluation`                 varchar(5000) NOT NULL DEFAULT '' COMMENT '救援评价',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='救援过程表';


DROP TABLE IF EXISTS `signature`;
CREATE TABLE `signature`(
    `id`                         int unsigned NOT NULL AUTO_INCREMENT COMMENT '签字编号',
    `create_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`                  tinyint NOT NULL DEFAULT '0',
    `rescue_teacher_id`          int unsigned NOT NULL DEFAULT '0' COMMENT '救援老师编号',
    `rescue_target_id`           int unsigned NOT NULL DEFAULT '0' COMMENT '救援对象编号',
    `image`                      text COMMENT '签字图片编码',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_teacherId_targetId` (`rescue_teacher_id`,`rescue_target_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='签字表';