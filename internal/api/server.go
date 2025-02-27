package api

import (
	"github.com/gorilla/mux"
	"mcp-db/internal/config"
	"mcp-db/internal/k8s"
	"net/http"
)

// Server 封装了HTTP服务器及其依赖
type Server struct {
	router    *mux.Router
	k8sClient *k8s.Client
	config    *config.Config
}

// NewServer 创建一个新的HTTP服务器
func NewServer(cfg *config.Config, client *k8s.Client) *Server {
	server := &Server{
		router:    mux.NewRouter(),
		k8sClient: client,
		config:    cfg,
	}

	// 注册路由
	server.setupRoutes()

	return server
}

// setupRoutes 设置HTTP路由
func (s *Server) setupRoutes() {
	// API版本前缀
	api := s.router.PathPrefix("/api").Subrouter()

	// 数据库集群相关端点
	api.HandleFunc("/databases", s.GetDatabases).Methods(http.MethodGet, http.MethodPost)
	api.HandleFunc("/databases/create", s.CreateDatabase).Methods(http.MethodPost)
	api.HandleFunc("/databases/delete", s.DeleteDatabase).Methods(http.MethodDelete, http.MethodPost)

	// 健康检查
	s.router.HandleFunc("/health", s.HealthCheck).Methods(http.MethodGet)
}

// Start 启动HTTP服务器
func (s *Server) Start() error {
	return http.ListenAndServe(":"+s.config.Port, s.router)
}

// HealthCheck 提供健康检查接口
func (s *Server) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]bool{"healthy": true})
}
