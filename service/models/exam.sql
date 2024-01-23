-- 有缓存 goctl model mysql ddl -src exam.sql -dir . -c
-- 无缓存 goctl model mysql ddl -src exam.sql -dir .
DROP TABLE IF EXISTS `subjective_exam`; -- 主观题考试成绩表
CREATE TABLE `subjective_exam`(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`   tinyint NOT NULL DEFAULT '0',
    `user_id`     int unsigned NOT NULL DEFAULT '0',
    `result`      tinyint(1) NOT NULL DEFAULT '0' COMMENT '结果 0:不合格 1:合格',
    `time`        varchar(100) DEFAULT '' NOT NULL COMMENT '考试日期',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_userId_time` (`user_id`,`time`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='主观题考试成绩表';

DROP TABLE IF EXISTS `objective_exam`; -- 客观题考试成绩表
CREATE TABLE `objective_exam`(
    `id`          int unsigned NOT NULL AUTO_INCREMENT,
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`   tinyint NOT NULL DEFAULT '0',
    `user_id`     int unsigned NOT NULL DEFAULT '0',
    `result`      tinyint(1) NOT NULL DEFAULT '0' COMMENT '结果 0:不合格 1:合格',
    `time`        varchar(100) DEFAULT '' NOT NULL COMMENT '考试日期',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_userId_time` (`user_id`,`time`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='客观题考试成绩表';
