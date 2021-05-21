# 腾讯云小工具

用 `letsencrypt` 来生成免费证书，并同步至腾讯云 cdn / ecdn。因为 `letsencrypt` 证书有效期只有三个月，所以开个定时器实时更新  

## 部署

见 [deployment.md](deployment.md)

## 目录结构
```
- config #配置文件
- src # 代码目录
    - certificate/ # 证书处理核心代码
        - issue.go # 通过 shell 脚本的方式签发 `letsencrypt` 证书
        - sync.go # 同步证书至 cdn / ecdn
        - task.go # 定时器
    - cmd/certificate-monitor/main.go # 程序入口
    - config/ # 配置文件解析
    - db/ # db 操作相关代码
    - tools/utils.go # 工具类
    - web/ # http 页面处理相关代码
- web # 静态文件和模板目录
```

## 访问

### 查看域名列表

/list

### 填加域名

/add

## Notice

项目不包含鉴权服务，谨慎部署至公网，注意信息安全

## CloudBase

[![](https://main.qcloudimg.com/raw/67f5a389f1ac6f3b4d04c7256438e44f.svg)](https://console.cloud.tencent.com/tcb/env/index?action=CreateAndDeployCloudBaseProject&appUrl=https%3A%2F%2Fgithub.com%2Factors315%2Fqcloud-tools&branch=master)