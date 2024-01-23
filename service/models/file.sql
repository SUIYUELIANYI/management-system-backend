-- 有缓存 goctl model mysql ddl -src file.sql -dir . -c
-- 无缓存 goctl model mysql ddl -src file.sql -dir .
DROP TABLE IF EXISTS `file`;
CREATE TABLE `file`(
    `id`                         int unsigned NOT NULL AUTO_INCREMENT COMMENT '文件编号',
    `create_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`                  tinyint NOT NULL DEFAULT '0',
    `file_name`                  varchar(100) NOT NULL DEFAULT '' COMMENT '文件名',
    `folder_id`                  int unsigned NOT NULL DEFAULT '0' COMMENT '文件夹名',
    `url`                        varchar(500) NOT NULL DEFAULT '' COMMENT '文件地址',
    PRIMARY KEY(`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文件表';

CREATE TABLE `folder`(
    `id`                         int unsigned NOT NULL AUTO_INCREMENT COMMENT '文件夹编号',
    `create_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`                  tinyint NOT NULL DEFAULT '0',
    `folder_name`                varchar(100) NOT NULL DEFAULT '' COMMENT '文件名',
    `role`                       int NOT NULL DEFAULT '0' COMMENT '文件查看权限',
    PRIMARY KEY(`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文件夹表';

DROP TABLE IF EXISTS `viewing_record`;
CREATE TABLE `viewing_record`(
    `id`                         int unsigned NOT NULL AUTO_INCREMENT COMMENT '观看记录编号',
    `create_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time`                datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`                  tinyint NOT NULL DEFAULT '0',
    `user_id`                    int unsigned NOT NULL DEFAULT '0',
    `file_id`                    int unsigned NOT NULL DEFAULT '0',
    `duration`                   int unsigned NOT NULL DEFAULT '0' COMMENT '观看时长（单位s）',
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_userId_fileId` (`user_id`,`file_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='观看记录表';