package types

// CreateDatabaseRequest 包含创建数据库的请求参数
type CreateDatabaseRequest struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Version       string `json:"version,omitempty"`
	Namespace     string `json:"namespace,omitempty"`
	Token         string `json:"token,omitempty"`
	CPULimit      string `json:"cpu_limit,omitempty"`
	MemoryLimit   string `json:"memory_limit,omitempty"`
	CPURequest    string `json:"cpu_request,omitempty"`
	MemoryRequest string `json:"memory_request,omitempty"`
	Storage       string `json:"storage,omitempty"`
}

// GetDatabasesRequest 包含获取数据库列表的请求参数
type GetDatabasesRequest struct {
	Namespace string `json:"namespace,omitempty"`
	Token     string `json:"token,omitempty"`
	Type      string `json:"type,omitempty"`
}

// DeleteDatabaseRequest 包含删除数据库的请求参数
type DeleteDatabaseRequest struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	Token     string `json:"token,omitempty"`
}

// Response 定义了API响应的通用结构
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// DBClusterInfo 包含数据库集群的详细信息
type DBClusterInfo struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Version        string `json:"version"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	CPULimit       string `json:"cpu_limit,omitempty"`
	MemoryLimit    string `json:"memory_limit,omitempty"`
	CPURequest     string `json:"cpu_request,omitempty"`
	MemoryRequest  string `json:"memory_request,omitempty"`
	Storage        string `json:"storage,omitempty"`
	AccessMode     string `json:"access_mode,omitempty"`
	Replicas       int64  `json:"replicas,omitempty"`
	ServiceAccount string `json:"service_account,omitempty"`
}
