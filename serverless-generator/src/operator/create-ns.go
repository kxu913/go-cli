package operator

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNS(ns string) *corev1.Namespace {
	typeMeta := &metav1.TypeMeta{
		Kind:       "Namespace",
		APIVersion: "v1",
	}
	objectMeta := &metav1.ObjectMeta{
		Name: ns,
	}
	namespace := &corev1.Namespace{
		TypeMeta:   *typeMeta,
		ObjectMeta: *objectMeta,
	}
	clientSet := ClientSet()
	createdNS, err := clientSet.CoreV1().Namespaces().Create(context.TODO(), namespace, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	return createdNS
}

func DestroyNS(ns string) {

	clientSet := ClientSet()

	err := clientSet.CoreV1().Namespaces().Delete(context.TODO(), ns, metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}

}
