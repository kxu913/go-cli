package operator

import (
	"context"
	"fmt"
	"serverless-generator/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/retry"
)

func DeployService(ns string, metaData *model.MetaData, replicas int, container *model.Container) *unstructured.Unstructured {
	client := DynamicClient()
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	object := map[string]any{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata":   model.MetaDataToMap(metaData),
		"spec":       model.ToDeploymentSpec(metaData, replicas, container),
	}
	deployment := &unstructured.Unstructured{
		Object: object,
	}

	yaml, err := client.Resource(deploymentRes).Namespace(ns).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	return yaml
}

func DestroyDeployment(ns string, metaData *model.MetaData) {
	client := DynamicClient()
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	if err := client.Resource(deploymentRes).Namespace(ns).Delete(context.TODO(), metaData.Name+"-"+metaData.Version, deleteOptions); err != nil {
		panic(err)
	}
}

func UpdateDeployment(ns string, metaData *model.MetaData, image string) {
	client := DynamicClient()
	deploymentRes := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := client.Resource(deploymentRes).Namespace(ns).Get(context.TODO(), metaData.Name+"-"+metaData.Version, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("failed to get latest version of Deployment: %v", getErr))
		}
		containers, found, err := unstructured.NestedSlice(result.Object, "spec", "template", "spec", "containers")
		if err != nil || !found || containers == nil {
			panic(fmt.Errorf("deployment containers not found or error in spec: %v", err))
		}
		container := containers[0].(map[string]any)
		if err := unstructured.SetNestedField(container, image, "image"); err != nil {
			panic(err)
		}

		if err := unstructured.SetNestedField(result.Object, containers, "spec", "template", "spec", "containers"); err != nil {
			panic(err)
		}
		_, updateErr := client.Resource(deploymentRes).Namespace(ns).Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("update failed: %v", retryErr))
	}
}
