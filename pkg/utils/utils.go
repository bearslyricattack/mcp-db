package utils

import (
	"os"
	"strings"
)

// GetEnvWithDefault 获取环境变量，如果不存在则返回默认值
func GetEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// SanitizeResourceName 确保资源名称符合Kubernetes命名规范
func SanitizeResourceName(name string) string {
	// 将大写转为小写
	name = strings.ToLower(name)

	// 将不符合规范的字符替换为"-"
	invalidChars := []string{" ", "_", ".", ","}
	for _, char := range invalidChars {
		name = strings.ReplaceAll(name, char, "-")
	}

	// 移除前后的"-"
	name = strings.Trim(name, "-")

	// 确保名称长度不超过63个字符
	if len(name) > 63 {
		name = name[:63]
	}

	return name
}

// IsValidDatabaseType 检查数据库类型是否有效
func IsValidDatabaseType(dbType string) bool {
	validTypes := []string{"postgresql", "mysql", "redis", "mongodb", "kafka", "milvus"}

	for _, validType := range validTypes {
		if dbType == validType {
			return true
		}
	}

	return false
}
