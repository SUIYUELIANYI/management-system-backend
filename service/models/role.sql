-- 有缓存 goctl model mysql ddl -src role.sql -dir . -c
-- 无缓存 goctl model mysql ddl -src role.sql -dir . 
DROP TABLE IF EXISTS `role_change`;
CREATE TABLE `role_change`(
    `id`            int unsigned NOT NULL AUTO_INCREMENT,
    `create_time`   datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time`   datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `delete_time`   datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `del_state`     tinyint NOT NULL DEFAULT '0',
    `user_id`       int NOT NULL DEFAULT '0' COMMENT '身份变动用户id',
    `operator_id`   int NOT NULL DEFAULT '0' COMMENT '操作人id',
    `new_role`      int NOT NULL DEFAULT '0' COMMENT '新身份',
    `old_role`      int NOT NULL DEFAULT '0' COMMENT '旧身份',
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='身份变动表';