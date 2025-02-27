package server

import (
	"fmt"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Server implementation
type Server struct {
	config        Config
	clientset     *kubernetes.Clientset
	dynamicClient dynamic.Interface
}

func NewServer(config Config) (*Server, error) {
	// Use provided kubeconfig or try default locations
	kubeconfigPath := config.KubeconfigPath
	if kubeconfigPath == "" {
		if envPath := os.Getenv("KUBECONFIG"); envPath != "" {
			kubeconfigPath = envPath
		} else if home := homedir.HomeDir(); home != "" {
			kubeconfigPath = filepath.Join(home, ".kube", "config")
		} else {
			return nil, fmt.Errorf("kubeconfig not found, please specify path")
		}
	}

	// Load kubeconfig
	log.Printf("Loading kubeconfig from: %s", kubeconfigPath)
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	// Create dynamic client for CRDs
	dynamicClient, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &Server{
		config:        config,
		clientset:     clientset,
		dynamicClient: dynamicClient,
	}, nil
}

func (s *Server) Start() error {
	// Register handlers
	http.HandleFunc("/api/databases", s.handleGetDatabases)
	http.HandleFunc("/api/databases/create", s.handleCreateDatabase)
	http.HandleFunc("/api/databases/delete", s.handleDeleteDatabase)

	// Start server
	log.Printf("Starting server on port %s", s.config.Port)
	return http.ListenAndServe(":"+s.config.Port, nil)
}
