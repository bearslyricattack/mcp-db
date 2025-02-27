# 数据库管理HTTP服务

这是一个简单的Go HTTP服务，用于在Kubernetes环境中管理数据库集群。服务支持创建、查询和删除多种类型的数据库，包括PostgreSQL、MySQL、Redis、MongoDB、Kafka和Milvus。

## 功能特性

- 创建数据库集群
- 查询数据库集群列表
- 删除数据库集群
- 支持多种数据库类型
- 自定义资源配置（CPU、内存、存储）

## 支持的数据库类型

- PostgreSQL
- MySQL
- Redis
- MongoDB
- Kafka
- Milvus

## API接口

### 创建数据库

```
POST /api/databases/create
```

请求体：

```json
{
  "name": "my-postgres",
  "type": "postgresql",
  "version": "14.8.0",
  "namespace": "default",
  "cpu_limit": "1000m",
  "memory_limit": "1024Mi",
  "cpu_request": "100m",
  "memory_request": "102Mi",
  "storage": "3Gi"
}
```

### 查询数据库列表

```
GET /api/databases?namespace=default&type=postgresql
```

或

```
POST /api/databases
```

请求体：

```json
{
  "namespace": "default",
  "type": "postgresql"
}
```

### 删除数据库

```
DELETE /api/databases/delete?name=my-postgres&namespace=default
```

或

```
POST /api/databases/delete
```

请求体：

```json
{
  "name": "my-postgres",
  "namespace": "default"
}
```

## 开发环境设置

### 先决条件

- Go 1.21+
- Kubernetes集群或Minikube
- KubeBlocks操作符已部署

### 构建和运行

1. 克隆仓库

```bash
git clone https://github.com/yourusername/database-manager.git
cd database-manager
```

1. 安装依赖

```bash
go mod tidy
```

1. 构建应用

```bash
go build -o database-manager cmd/server/main.go
```

1. 运行应用

```bash
export KUBECONFIG=/path/to/kubeconfig
export PORT=8080
export DEFAULT_NAMESPACE=default
./database-manager
```

## 在Kubernetes中运行

1. 构建Docker镜像

```bash
docker build -t database-manager:latest .
```

1. 将kubeconfig创建为secret

```bash
kubectl create secret generic database-manager-kubeconfig --from-file=config=/path/to/kubeconfig
```

1. 部署到Kubernetes

```bash
kubectl apply -f deploy/kubernetes/deployment.yaml
```

## 配置

服务支持以下环境变量：

- `KUBECONFIG`: Kubernetes配置文件路径
- `PORT`: HTTP服务器端口（默认：8080）
- `DEFAULT_NAMESPACE`: 默认命名空间（默认：default）

## 项目结构

```
database-manager/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── server.go
│   │   └── handlers.go
│   ├── config/
│   │   └── config.go
│   └── k8s/
│       ├── client.go
│       ├── database.go
│       └── rbac.go
├── pkg/
│   ├── types/
│   │   └── models.go
│   └── utils/
│       └── utils.go
├── deploy/
│   └── kubernetes/
│       └── deployment.yaml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```