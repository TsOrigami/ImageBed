# ImageV2 图床部署文档

## 环境要求
- Docker 20.10.0+
- Docker Compose v2.0.0+
- 至少 2GB 可用内存
- 至少 10GB 可用磁盘空间

## 快速部署

### 1. 解压部署包
```bash
tar -xzf imagev2-deploy.tar.gz
cd deploy
```

### 2. 导入镜像
```bash
docker load -i images.tar
```

### 3. 配置文件
检查并按需修改 `config.conf` 文件：
- MySQL 配置（默认无需修改）
- Redis 配置（默认无需修改）
- 服务器配置（根据实际部署环境修改 prefix）

### 4. 启动服务
```bash
docker-compose up -d
```

### 5. 验证服务
访问 http://localhost:8000 确认服务是否正常运行

## 目录结构
```
deploy/
├── images.tar          # Docker镜像包
├── config.conf         # 应用配置文件
└── README.md           # 部署说明文档
```

## 常用操作

### 查看服务状态
```bash
docker-compose ps
```

### 查看服务日志
```bash
docker-compose logs -f
```

### 停止服务
```bash
docker-compose down
```

### 重启服务
```bash
docker-compose restart
```

## 故障排查

1. 服务无法启动
   - 检查端口 8000, 3306, 6379 是否被占用
   - 查看容器日志 `docker-compose logs`

2. 无法连接数据库
   - 确认 MySQL 容器健康状态
   - 检查 config.conf 中的数据库配置

3. 上传图片失败
   - 检查 storage 目录权限
   - 确认磁盘空间是否充足

## 联系支持
如遇问题请提交 Issues 到项目仓库
