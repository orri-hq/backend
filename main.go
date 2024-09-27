package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Get the Kubernetes config, either from the in-cluster config or from the kubeconfig file (e.g., when running locally)
func getKubernetesConfig() (*rest.Config, error) {
	// Check if we are running inside a Kubernetes cluster
	if _, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token"); err == nil {
		// Use the in-cluster config
		return rest.InClusterConfig()
	}

	// Check if the KUBECONFIG environment variable is set
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		// If KUBECONFIG is not set, default to the local kubeconfig file
		home := homedir.HomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// Use the local kubeconfig file
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

// REST handler to fetch the OpenAPI schema from Kubernetes
func fetchOpenAPISchema(w http.ResponseWriter, r *http.Request) {
	// Get the Kubernetes configuration
	config, err := getKubernetesConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting Kubernetes config: %v", err), http.StatusInternalServerError)
		return
	}

	// Create the Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating Kubernetes client: %v", err), http.StatusInternalServerError)
		return
	}

	// Fetch the OpenAPI schema from the Kubernetes API
	openapi, err := clientset.RESTClient().Get().AbsPath("/openapi/v2").DoRaw(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching OpenAPI schema: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the OpenAPI schema as a response
	w.Header().Set("Content-Type", "application/json")
	w.Write(openapi)
}

func main() {
	// Create a new router
	r := mux.NewRouter()

	// Define the endpoint to fetch OpenAPI schema
	r.HandleFunc("/openapi", fetchOpenAPISchema).Methods("GET")

	// Start the HTTP server
	fmt.Println("Server starting at :8080")
	http.ListenAndServe("0.0.0.0:8080", r)
}
