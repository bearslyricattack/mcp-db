package k8s

import (
	"fmt"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"mcp-db/internal/config"
)

// Client 封装了与Kubernetes交互所需的客户端
type Client struct {
	// 标准Kubernetes客户端
	ClientSet *kubernetes.Clientset
	// 用于处理CRD的动态客户端
	DynamicClient dynamic.Interface
	// 配置
	Config *config.Config
}

// NewClient 创建一个新的Kubernetes客户端
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.KubeconfigPath == "" {
		return nil, fmt.Errorf("kubeconfig path not specified")
	}

	// 加载kubeconfig
	log.Printf("Loading kubeconfig from: %s", cfg.KubeconfigPath)
	restConfig, err := clientcmd.BuildConfigFromFlags("", cfg.KubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	// 创建clientset
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	// 创建用于CRD的动态客户端
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &Client{
		ClientSet:     clientset,
		DynamicClient: dynamicClient,
		Config:        cfg,
	}, nil
}
