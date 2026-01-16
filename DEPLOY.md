# Life RPG 服务器部署指南

## 部署架构

```
┌─────────────────────────────────────────────────┐
│                   服务器                          │
│  ┌─────────────────┐    ┌─────────────────────┐ │
│  │  Frontend (80)  │───▶│  Backend (8080)     │ │
│  │    Nginx        │    │    Go + Gin         │ │
│  └─────────────────┘    └─────────────────────┘ │
│                                   │             │
│                                   ▼             │
│                         ┌─────────────────────┐ │
│                         │  MySQL (mysql01)    │ │
│                         └─────────────────────┘ │
└─────────────────────────────────────────────────┘
```

## 部署步骤

### 1. 拉取代码

```bash
cd /opt
git clone https://github.com/YOUR_USERNAME/my-life-rpg.git
cd my-life-rpg
```

### 2. 使用 Docker Compose 一键部署

```bash
# 构建并启动所有服务
docker-compose up -d --build

# 查看运行状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 3. 单独部署（可选）

如果不使用 docker-compose，可以手动部署：

**后端：**
```bash
cd backend
docker build -t life-rpg-backend:latest .
docker run -d \
  --name life-rpg-backend \
  -p 8080:8080 \
  -e DB_HOST=mysql01 \
  -e DB_PORT=3306 \
  -e DB_USER=root \
  -e DB_PASSWORD=123456 \
  -e DB_NAME=life_rpg \
  --network life-rpg-network \
  life-rpg-backend:latest
```

**前端：**
```bash
cd frontend
docker build -t life-rpg-frontend:latest .
docker run -d \
  --name life-rpg-frontend \
  -p 80:80 \
  --network life-rpg-network \
  life-rpg-frontend:latest
```

## 关于 IP 地址配置

| 场景 | nginx.conf 中的 proxy_pass |
|------|---------------------------|
| Docker Compose (推荐) | `http://life-rpg-backend:8080` |
| 同一宿主机 | `http://host.docker.internal:8080` |
| 不同机器 | `http://内网IP:8080` |

**当前配置使用容器名 `life-rpg-backend`**，这是 Docker 内部 DNS 解析，只在同一个 Docker 网络内有效。

## 更新部署

```bash
cd /opt/my-life-rpg

# 拉取最新代码
git pull

# 重新构建并启动
docker-compose up -d --build
```

## 常用命令

```bash
# 停止所有服务
docker-compose down

# 重启服务
docker-compose restart

# 查看后端日志
docker logs -f life-rpg-backend

# 进入前端容器
docker exec -it life-rpg-frontend sh
```

## 访问地址

- 前端: `http://服务器IP:80`
- 后端 API: `http://服务器IP:8080/api`
