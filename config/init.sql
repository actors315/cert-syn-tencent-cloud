CREATE DATABASE `qcloud-tools` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;

GRANT ALL PRIVILEGES ON `qcloud-tools`.* TO 'db_qcloud'@'%' IDENTIFIED BY '58117aec3b3252a97be0';

CREATE TABLE IF NOT EXISTS `qcloud-tools`.`issue_info` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `secret_id` varchar(64) NOT NULL DEFAULT '' COMMENT '腾讯云密钥ID',
  `secret_key` varchar(64) NOT NULL DEFAULT '' COMMENT '腾讯云密钥key',
  `dns_api` varchar(8) NOT NULL DEFAULT '' COMMENT 'DNS API',
  `app_id` varchar(16) NOT NULL DEFAULT '' COMMENT 'DNS产商ID标识',
  `app_id_value` varchar(64) NOT NULL DEFAULT '' COMMENT 'DNS产商ID',
  `app_key` varchar(32) NOT NULL DEFAULT '' COMMENT 'DNS产商KEY标识',
  `app_key_value` varchar(64) NOT NULL DEFAULT '' COMMENT 'DNS产商KEY',
  `type` varchar(8) NOT NULL DEFAULT 'cdn' COMMENT 'cdn类型',
  `cdn_domain` varchar(128) NOT NULL DEFAULT '' COMMENT '加速域名|需求配置一致',
  `main_domain` varchar(128) NOT NULL DEFAULT '' COMMENT '主域名|泛解析时该值不支持为子域名',
  `extra_domain` text NOT NULL COMMENT '额外域名|需以 -d domain 的方式，多个以空格分隔',
  `last_issue_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最后签发时间',
  `last_check_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最后一次执行脚本时间',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uniq_domain` (`cdn_domain`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `issue_history` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `issue_domain` varchar(128) NOT NULL DEFAULT '' COMMENT '证书根域名',
  `public_key` longtext NOT NULL COMMENT '公钥内容',
  `private_key` longtext NOT NULL COMMENT '私钥内容',
  `created_at` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

