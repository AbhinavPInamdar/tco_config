package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"tco-configurator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)


type KubeClientInterface interface {
	GetTeamBudget(name, namespace string) (*v1.TeamBudget, error)
	UpdateTeamBudgetStatus(name, namespace string, newUsage int64) error
}

type Client struct {
	dynClient dynamic.Interface
	gvr       schema.GroupVersionResource
}

func NewClient() (*Client, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := filepath.Join(homeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
		}
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

	return &Client{
		dynClient: dynClient,
		gvr:       gvr,
	}, nil

}

func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Could not determine home directory")
	}
	return home
}

func (c *Client) ListTeamBudgets() ([]v1.TeamBudget, error) {
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




func (c *Client) GetTeamBudget(name, namespace string) (*v1.TeamBudget, error) {
	obj, err := c.dynClient.Resource(c.gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list TeamBudget: %w", err)
	}
	var tb v1.TeamBudget
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, &tb)
	if err != nil {
		return nil, fmt.Errorf("conversion error: %w", err)
	}

	return &tb, nil
}


func (c *Client) UpdateTeamBudgetStatus(name, namespace string, newUsage int64) error {
	budget, err := c.GetTeamBudget(name, namespace)
	if err != nil {
		return fmt.Errorf("expected teambudget, instead got %v", err)
	}

	budget.Status.CurrentUsage = newUsage

	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(budget)
	if err != nil {
		return fmt.Errorf("conversion error: %w", err)
	}

	unstructuredObj := &unstructured.Unstructured{
		Object: unstructuredMap,
	}

	_, err = c.dynClient.
		Resource(c.gvr).
		Namespace(namespace).
		UpdateStatus(
			context.TODO(),
			unstructuredObj,
			metav1.UpdateOptions{},
		)

	if err != nil {
		return fmt.Errorf("Failed to update teambudget: %w", err)
	}

	return nil
}
