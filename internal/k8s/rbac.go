package k8s

import (
	"context"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateServiceAccount 创建数据库集群用的服务账号
func (c *Client) CreateServiceAccount(ctx context.Context, name, namespace string) error {
	sa := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"sealos-db-provider-cr":        name,
				"app.kubernetes.io/instance":   name,
				"app.kubernetes.io/managed-by": "kbcli",
			},
		},
	}

	_, err := c.ClientSet.CoreV1().ServiceAccounts(namespace).Create(ctx, sa, metav1.CreateOptions{})
	return err
}

// CreateRole 创建数据库集群用的角色
func (c *Client) CreateRole(ctx context.Context, name, namespace string) error {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"sealos-db-provider-cr":        name,
				"app.kubernetes.io/instance":   name,
				"app.kubernetes.io/managed-by": "kbcli",
			},
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
		},
	}

	_, err := c.ClientSet.RbacV1().Roles(namespace).Create(ctx, role, metav1.CreateOptions{})
	return err
}

// CreateRoleBinding 创建数据库集群用的角色绑定
func (c *Client) CreateRoleBinding(ctx context.Context, name, namespace string) error {
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"sealos-db-provider-cr":        name,
				"app.kubernetes.io/instance":   name,
				"app.kubernetes.io/managed-by": name,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     name,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "ServiceAccount",
				Name: name,
			},
		},
	}

	_, err := c.ClientSet.RbacV1().RoleBindings(namespace).Create(ctx, roleBinding, metav1.CreateOptions{})
	return err
}
