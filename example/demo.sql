CREATE TABLE `demo` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
    `name` varchar(256) NOT NULL COMMENT 'name',
    `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态',
    `content` text NOT NULL COMMENT '配置id',
    `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4