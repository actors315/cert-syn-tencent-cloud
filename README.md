# 腾讯云小工具

用 `letsencrypt` 来生成免费证书，并同步至腾讯云 cdn / ecdn。因为 `letsencrypt` 证书有效期只有三个月，所以开个定时器实时更新  

## 准备

需要有一个 MySQL 数据库，初始化 sql 见 `config/init.sql`，根据实际需要调整。

## 配置项

config 目录下,重建 config.simple.yaml 为 config.yaml

修改数据库对应配置

```
db:
  host: "localhost"
  port: 3306
  database: "qcloud-tools"
  user: "db_qcloud"
  password: "58117aec3b3252a97be0"

```

## 查看域名列表


## 填加域名

## Notice