package k8s

import (
	"context"
	"fmt"
	"log"
	"mcp-db/pkg/types"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// DatabaseClusterGVR 定义了数据库集群的GroupVersionResource
var DatabaseClusterGVR = schema.GroupVersionResource{
	Group:    "apps.kubeblocks.io",
	Version:  "v1alpha1",
	Resource: "clusters",
}

// DatabaseConfigs 存储了不同数据库类型的配置信息
var DatabaseConfigs = map[string]struct {
	Definition string
	Version    string
	Component  string
}{
	"postgresql": {
		Definition: "postgresql",
		Version:    "postgresql-%s",
		Component:  "postgresql",
	},
	"mysql": {
		Definition: "apecloud-mysql",
		Version:    "ac-mysql-%s",
		Component:  "mysql",
	},
	"redis": {
		Definition: "redis",
		Version:    "redis-%s",
		Component:  "redis",
	},
	"mongodb": {
		Definition: "mongodb",
		Version:    "mongodb-%s",
		Component:  "mongodb",
	},
	"kafka": {
		Definition: "kafka",
		Version:    "kafka-%s",
		Component:  "kafka",
	},
	"milvus": {
		Definition: "milvus",
		Version:    "milvus-%s",
		Component:  "milvus",
	},
}

// DefaultVersions 存储了不同数据库类型的默认版本
var DefaultVersions = map[string]string{
	"postgresql": "14.8.0",
	"mysql":      "8.0.30-1",
	"redis":      "7.0.6",
	"mongodb":    "6.0",
	"kafka":      "3.3.2",
	"milvus":     "2.4.5",
}

// CreateDatabaseCluster 创建一个新的数据库集群
func (c *Client) CreateDatabaseCluster(ctx context.Context, req *types.CreateDatabaseRequest) error {
	// 检查数据库类型是否支持
	dbConfig, ok := DatabaseConfigs[req.Type]
	if !ok {
		return fmt.Errorf("unsupported database type: %s", req.Type)
	}

	// 使用用户提供的版本或默认版本
	version := req.Version
	if version == "" {
		if defaultVer, ok := DefaultVersions[req.Type]; ok {
			version = defaultVer
		} else {
			return fmt.Errorf("version not provided and no default available")
		}
	}

	// 格式化版本字符串
	formattedVersion := fmt.Sprintf(dbConfig.Version, version)

	// 创建RBAC资源
	if err := c.CreateServiceAccount(ctx, req.Name, req.Namespace); err != nil {
		return fmt.Errorf("failed to create ServiceAccount: %w", err)
	}

	if err := c.CreateRole(ctx, req.Name, req.Namespace); err != nil {
		return fmt.Errorf("failed to create Role: %w", err)
	}

	if err := c.CreateRoleBinding(ctx, req.Name, req.Namespace); err != nil {
		return fmt.Errorf("failed to create RoleBinding: %w", err)
	}

	// 等待ServiceAccount的token生成
	time.Sleep(1 * time.Second)

	// 创建数据库集群的unstructured对象
	cluster := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps.kubeblocks.io/v1alpha1",
			"kind":       "Cluster",
			"metadata": map[string]interface{}{
				"name":      req.Name,
				"namespace": req.Namespace,
				"finalizers": []string{
					"cluster.kubeblocks.io/finalizer",
				},
				"labels": map[string]interface{}{
					"clusterdefinition.kubeblocks.io/name": dbConfig.Definition,
					"clusterversion.kubeblocks.io/name":    formattedVersion,
					"sealos-db-provider-cr":                req.Name,
				},
			},
			"spec": map[string]interface{}{
				"affinity": map[string]interface{}{
					"nodeLabels":      map[string]interface{}{},
					"podAntiAffinity": "Preferred",
					"tenancy":         "SharedNode",
					"topologyKeys": []string{
						"kubernetes.io/hostname",
					},
				},
				"clusterDefinitionRef": dbConfig.Definition,
				"clusterVersionRef":    formattedVersion,
				"componentSpecs": []map[string]interface{}{
					{
						"componentDefRef": dbConfig.Component,
						"monitor":         true,
						"name":            dbConfig.Component,
						"replicas":        1,
						"resources": map[string]interface{}{
							"limits": map[string]interface{}{
								"cpu":    req.CPULimit,
								"memory": req.MemoryLimit,
							},
							"requests": map[string]interface{}{
								"cpu":    req.CPURequest,
								"memory": req.MemoryRequest,
							},
						},
						"serviceAccountName": req.Name,
						"switchPolicy": map[string]interface{}{
							"type": "Noop",
						},
						"volumeClaimTemplates": []map[string]interface{}{
							{
								"name": "data",
								"spec": map[string]interface{}{
									"accessModes": []string{
										"ReadWriteOnce",
									},
									"resources": map[string]interface{}{
										"requests": map[string]interface{}{
											"storage": req.Storage,
										},
									},
								},
							},
						},
					},
				},
				"terminationPolicy": "Delete",
				"tolerations":       []interface{}{},
			},
		},
	}

	// 创建集群
	_, err := c.DynamicClient.
		Resource(DatabaseClusterGVR).
		Namespace(req.Namespace).
		Create(ctx, cluster, metav1.CreateOptions{})

	return err
}

// GetDatabaseClusters 获取指定命名空间中的数据库集群列表
func (c *Client) GetDatabaseClusters(ctx context.Context, namespace, dbType string) ([]types.DBClusterInfo, error) {
	// 获取集群列表
	clusters, err := c.DynamicClient.
		Resource(DatabaseClusterGVR).
		Namespace(namespace).
		List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	result := make([]types.DBClusterInfo, 0)

	for _, cluster := range clusters.Items {
		// 提取元数据
		metadata, found, err := unstructured.NestedMap(cluster.Object, "metadata")
		if err != nil || !found {
			log.Printf("Failed to get metadata for cluster: %v", err)
			continue
		}

		name, _ := metadata["name"].(string)
		creationTimestamp, _ := metadata["creationTimestamp"].(string)

		// 获取标签
		labels, found, _ := unstructured.NestedMap(metadata, "labels")
		if !found {
			labels = map[string]interface{}{}
		}

		// 提取定义类型
		definitionType, _ := labels["clusterdefinition.kubeblocks.io/name"].(string)
		versionString, _ := labels["clusterversion.kubeblocks.io/name"].(string)

		// 按类型过滤
		if dbType != "" && definitionType != dbType {
			continue
		}

		// 提取状态
		status := "Unknown"
		statusObj, found, _ := unstructured.NestedMap(cluster.Object, "status")
		if found {
			if phase, ok := statusObj["phase"].(string); ok {
				status = phase
			}
		}

		// 提取规格详情
		spec, found, _ := unstructured.NestedMap(cluster.Object, "spec")
		if !found {
			spec = map[string]interface{}{}
		}

		// 提取组件规格
		componentSpecsUntyped, found, _ := unstructured.NestedSlice(spec, "componentSpecs")

		// 默认值
		cpuLimit := ""
		memLimit := ""
		cpuRequest := ""
		memRequest := ""
		storage := ""
		accessMode := ""
		var replicas int64 = 0
		serviceAccount := ""

		if found && len(componentSpecsUntyped) > 0 {
			// 获取第一个组件
			mainComponent, ok := componentSpecsUntyped[0].(map[string]interface{})
			if ok {
				// 获取资源
				resources, found, _ := unstructured.NestedMap(mainComponent, "resources")
				if found {
					limits, limitsFound, _ := unstructured.NestedMap(resources, "limits")
					if limitsFound {
						if cpu, ok := limits["cpu"].(string); ok {
							cpuLimit = cpu
						}
						if mem, ok := limits["memory"].(string); ok {
							memLimit = mem
						}
					}

					requests, reqFound, _ := unstructured.NestedMap(resources, "requests")
					if reqFound {
						if cpu, ok := requests["cpu"].(string); ok {
							cpuRequest = cpu
						}
						if mem, ok := requests["memory"].(string); ok {
							memRequest = mem
						}
					}
				}

				// 获取副本数
				if rep, ok := mainComponent["replicas"].(int64); ok {
					replicas = rep
				}

				// 获取服务账号
				if sa, ok := mainComponent["serviceAccountName"].(string); ok {
					serviceAccount = sa
				}

				// 获取卷声明模板
				volumeTemplates, found, _ := unstructured.NestedSlice(mainComponent, "volumeClaimTemplates")
				if found && len(volumeTemplates) > 0 {
					// 寻找数据卷
					for _, volUntyped := range volumeTemplates {
						vol, ok := volUntyped.(map[string]interface{})
						if !ok {
							continue
						}

						volName, _ := vol["name"].(string)
						if volName == "data" {
							spec, specFound, _ := unstructured.NestedMap(vol, "spec")
							if specFound {
								resourcesMap, resFound, _ := unstructured.NestedMap(spec, "resources")
								if resFound {
									requestsMap, reqFound, _ := unstructured.NestedMap(resourcesMap, "requests")
									if reqFound {
										if st, ok := requestsMap["storage"].(string); ok {
											storage = st
										}
									}
								}

								accessModes, modesFound, _ := unstructured.NestedStringSlice(spec, "accessModes")
								if modesFound && len(accessModes) > 0 {
									accessMode = accessModes[0]
								}
							}

							break
						}
					}
				}
			}
		}

		// 创建集群信息
		clusterInfo := types.DBClusterInfo{
			Name:           name,
			Type:           definitionType,
			Version:        versionString,
			Status:         status,
			CreatedAt:      creationTimestamp,
			CPULimit:       cpuLimit,
			MemoryLimit:    memLimit,
			CPURequest:     cpuRequest,
			MemoryRequest:  memRequest,
			Storage:        storage,
			AccessMode:     accessMode,
			Replicas:       replicas,
			ServiceAccount: serviceAccount,
		}

		result = append(result, clusterInfo)
	}

	return result, nil
}

// DeleteDatabaseCluster 删除指定的数据库集群
func (c *Client) DeleteDatabaseCluster(ctx context.Context, name, namespace string) error {
	return c.DynamicClient.
		Resource(DatabaseClusterGVR).
		Namespace(namespace).
		Delete(ctx, name, metav1.DeleteOptions{})
}
