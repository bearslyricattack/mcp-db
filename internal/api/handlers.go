package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mcp-db/pkg/types"
	"net/http"
)

// CreateDatabase 处理创建数据库的请求
func (s *Server) CreateDatabase(w http.ResponseWriter, r *http.Request) {
	var req types.CreateDatabaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// 验证请求
	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Database name is required")
		return
	}

	// 设置默认值
	if req.Type == "" {
		req.Type = "postgresql"
	}
	if req.Namespace == "" {
		req.Namespace = s.config.DefaultNamespace
	}
	if req.CPULimit == "" {
		req.CPULimit = "1000m"
	}
	if req.MemoryLimit == "" {
		req.MemoryLimit = "1024Mi"
	}
	if req.CPURequest == "" {
		req.CPURequest = "100m"
	}
	if req.MemoryRequest == "" {
		req.MemoryRequest = "102Mi"
	}
	if req.Storage == "" {
		req.Storage = "3Gi"
	}

	// 创建数据库集群
	ctx := context.Background()
	if err := s.k8sClient.CreateDatabaseCluster(ctx, &req); err != nil {
		log.Printf("Failed to create database cluster: %v", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create database cluster: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, types.Response{
		Success: true,
		Message: fmt.Sprintf("Successfully created database cluster '%s'", req.Name),
	})
}

// GetDatabases 处理获取数据库列表的请求
func (s *Server) GetDatabases(w http.ResponseWriter, r *http.Request) {
	var req types.GetDatabasesRequest

	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request format")
			return
		}
	} else {
		// 提取查询参数
		req.Namespace = r.URL.Query().Get("namespace")
		req.Type = r.URL.Query().Get("type")
		req.Token = r.URL.Query().Get("token")
	}

	// 设置默认命名空间
	if req.Namespace == "" {
		req.Namespace = s.config.DefaultNamespace
	}

	// 获取数据库集群
	ctx := context.Background()
	clusters, err := s.k8sClient.GetDatabaseClusters(ctx, req.Namespace, req.Type)
	if err != nil {
		log.Printf("Failed to get database clusters: %v", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to get database clusters: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, types.Response{
		Success: true,
		Message: fmt.Sprintf("Found %d database clusters in namespace '%s'", len(clusters), req.Namespace),
		Data:    clusters,
	})
}

// DeleteDatabase 处理删除数据库的请求
func (s *Server) DeleteDatabase(w http.ResponseWriter, r *http.Request) {
	var req types.DeleteDatabaseRequest

	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request format")
			return
		}
	} else {
		// 提取DELETE方法的查询参数
		req.Name = r.URL.Query().Get("name")
		req.Namespace = r.URL.Query().Get("namespace")
		req.Token = r.URL.Query().Get("token")
	}

	// 验证请求
	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Database name is required")
		return
	}

	// 设置默认命名空间
	if req.Namespace == "" {
		req.Namespace = s.config.DefaultNamespace
	}

	// 删除数据库集群
	ctx := context.Background()
	if err := s.k8sClient.DeleteDatabaseCluster(ctx, req.Name, req.Namespace); err != nil {
		log.Printf("Failed to delete database cluster: %v", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete database cluster: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, types.Response{
		Success: true,
		Message: fmt.Sprintf("Successfully deleted database cluster '%s'", req.Name),
	})
}

// respondWithJSON 帮助函数，用于发送JSON响应
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// respondWithError 帮助函数，用于发送错误响应
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, types.Response{
		Success: false,
		Message: message,
	})
}
