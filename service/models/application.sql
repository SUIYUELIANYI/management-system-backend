-- 有缓存 goctl model mysql ddl -src application.sql -dir . -c
-- 无缓存 goctl model mysql ddl -src application.sql -dir . 
DROP TABLE IF EXISTS `application_form`;
CREATE TABLE `application_form`(
    `id`                          int unsigned NOT NULL AUTO_INCREMENT COMMENT '申请表编号',
    `create_time`                 datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`                 datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time`                 datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`                   tinyint NOT NULL DEFAULT '0',
    `user_id`                     int NOT NULL DEFAULT '0' COMMENT '申请人id',
    `mobile`                      char(11) NOT NULL DEFAULT '' COMMENT '电话',
    `username`                    varchar(50) NOT NULL DEFAULT '' COMMENT '用户名称',
    `sex`                         int NOT NULL DEFAULT '0' COMMENT '性别 0-男 1-女',
    `email`                       varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
    `address`                     varchar(255) NOT NULL DEFAULT '0' COMMENT '地址',
    `birthday`                    varchar(20) NOT NULL DEFAULT '0000-00-00' COMMENT '生日 xxxx-xx-xx',
    `status`                      int NOT NULL DEFAULT '0' COMMENT '申请表状态 0-待审批 1-区域负责人通过 2-组织管理委员会通过 3-未通过',  -- 申请不通过移动到另一个，不通过超过3次则无法申请
    `regional_head_id`            int unsigned NOT NULL DEFAULT '0' COMMENT '区域负责人id', -- 仅记录最新一次修改的操作人
    `organizing_committee_id`     int unsigned NOT NULL DEFAULT '0' COMMENT '组委会成员id',
    `submission_time`             int unsigned NOT NULL DEFAULT '0' COMMENT '提交次数', -- 一共可以提交/修改3次
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_userId` (`user_id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='申请表';