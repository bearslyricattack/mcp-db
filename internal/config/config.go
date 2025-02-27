package config

import (
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

// Config 包含服务器配置选项
type Config struct {
	// Kubernetes配置文件路径
	KubeconfigPath string
	// HTTP服务器端口
	Port string
	// 默认命名空间
	DefaultNamespace string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	config := &Config{
		KubeconfigPath:   os.Getenv("KUBECONFIG"),
		Port:             os.Getenv("PORT"),
		DefaultNamespace: os.Getenv("DEFAULT_NAMESPACE"),
	}

	// 设置默认值
	if config.Port == "" {
		config.Port = "8080"
	}

	if config.DefaultNamespace == "" {
		config.DefaultNamespace = "default"
	}

	// 如果没有指定kubeconfig，尝试使用默认位置
	if config.KubeconfigPath == "" {
		if home := homedir.HomeDir(); home != "" {
			config.KubeconfigPath = filepath.Join(home, ".kube", "config")
		}
	}

	return config
}
