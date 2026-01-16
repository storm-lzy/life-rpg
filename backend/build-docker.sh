#!/bin/bash
# Life RPG 后端 Docker 镜像构建脚本

set -e

IMAGE_NAME="life-rpg-backend"
IMAGE_TAG="latest"
TAR_FILE="${IMAGE_NAME}.tar"

echo "🏗️  开始构建 Docker 镜像..."
docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

echo "📦  导出镜像为 ${TAR_FILE}..."
docker save -o ${TAR_FILE} ${IMAGE_NAME}:${IMAGE_TAG}

echo "✅  构建完成！"
echo "镜像文件: $(pwd)/${TAR_FILE}"
echo "文件大小: $(du -h ${TAR_FILE} | cut -f1)"
echo ""
echo "在目标服务器上执行以下命令导入镜像:"
echo "  docker load -i ${TAR_FILE}"
echo "  docker run -d -p 8080:8080 --name life-rpg-api ${IMAGE_NAME}:${IMAGE_TAG}"
