-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS `users` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `users`;

-- 创建用户表
CREATE TABLE IF NOT EXISTS `users`
(
    `id`          bigint(20)   NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username`    varchar(50)  NOT NULL COMMENT '用户名',
    `password`    varchar(255) NOT NULL COMMENT '密码（加密后）',
    `nike_name`   varchar(50) DEFAULT NULL COMMENT '昵称',
    `create_time` datetime    DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete`      tinyint(1)  DEFAULT 0 COMMENT '是否删除 0-未删除 1-已删除',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户表';

CREATE TABLE IF NOT EXISTS `permission`
(
    `id`     bigint(20)   NOT NULL AUTO_INCREMENT COMMENT 'id',
    `permit` varchar(100) NOT NULL COMMENT '权限标识',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='权限表';

CREATE TABLE IF NOT EXISTS `user_permission`
(
    `user_id`       bigint(20) NOT NULL COMMENT '用户id',
    `permission_id` bigint(20) NOT NULL COMMENT '权限id'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户-权限关联表';
