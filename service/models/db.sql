-- 设置了软删除字段，但业务上还没有实现
-- 修改全局sql_mode：set global sql_mode='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';
-- 设置时区：set global time_zone='Asia/Shanghai';
DROP DATABASE IF EXISTS `management_system`;
CREATE DATABASE `management_system`;
USE `management_system`;