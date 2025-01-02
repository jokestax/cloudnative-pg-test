package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getK8sClient(incluster, kubeconfig string) (*kubernetes.Clientset, error) {
	if incluster == "true" {
		fmt.Print("Using in-cluster config")
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get incluster config: %w", err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, fmt.Errorf("failed to get clientset: %w", err)
		}
		return clientset, nil
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, fmt.Errorf("failed to get clientset: %w", err)
		}
		return clientset, nil
	}
}

func main() {
	inCluster := os.Getenv("IN_CLUSTER")
	kubeconfig := os.Getenv("KUBECONFIG")
	clientset, err := getK8sClient(inCluster, kubeconfig)
	if err != nil {
		panic(err)
	}

	ctx := context.TODO()
	secrets, err := clientset.CoreV1().Secrets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	connStr := ""

	for _, secret := range secrets.Items {
		if secret.Name == "cluster-example-app" {
			for k := range secret.Data {
				fmt.Printf("%s : %s\n", k, string(secret.Data[k]))
				if k == "connpass" {
					connStr = string(secret.Data[k])
				}
			}
		}
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Connected to the database successfully!")
}
