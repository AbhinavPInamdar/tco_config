package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"tco-configurator/api/v1"
    "k8s.io/apimachinery/pkg/runtime"
)

type Client struct {
	dynClient dynamic.Interface
	gvr schema.GroupVersionResource

}

func NewClient() (*Client ,error) {
	kubeconfig := filepath.Join(homeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	gvr := schema.GroupVersionResource{
		Group:    "tco.io",
		Version:  "v1",
		Resource: "teambudgets",
	}

	return &Client {
		dynClient:dynClient,
		gvr:gvr,
	}, nil

}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Could not determine home directory")
	}
	return home
}

func (c *Client) ListTeamBudgets() ([]v1.TeamBudget, error  ) {
	
	tbList, err := c.dynClient.Resource(c.gvr).Namespace("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list TeamBudgets: %w", err)
	}
	var results []v1.TeamBudget

	for _, item := range tbList.Items {
		var tb v1.TeamBudget
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &tb)
    	if err != nil {
			return nil, fmt.Errorf("conversion error: %w", err)
    	}
		results = append(results, tb)

	}
	

	return results, nil
}




func GetTeamBudget(name, namespace string) {

	
}



func WatchTeamBudgets() {}
