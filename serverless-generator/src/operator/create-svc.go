package operator

import (
	"context"
	"serverless-generator/model"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateServiceAccount(ns string, metadata *model.MetaData) *corev1.ServiceAccount {
	clientSet := ClientSet()
	typeMeta := &metav1.TypeMeta{
		Kind:       "ServiceAccount",
		APIVersion: "v1",
	}
	acctName := metadata.Name + "-acct"
	objectMeta := &metav1.ObjectMeta{
		Name: acctName,
		Labels: map[string]string{
			"account": acctName,
		},
	}
	serviceAccount := &corev1.ServiceAccount{
		TypeMeta:   *typeMeta,
		ObjectMeta: *objectMeta,
	}
	acct, err := clientSet.CoreV1().ServiceAccounts(ns).Create(context.TODO(), serviceAccount, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	return acct

}

func CreateService(ns string, metadata *model.MetaData, port int32) *corev1.Service {
	clientSet := ClientSet()
	typeMeta := &metav1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	}
	serviceName := metadata.Name + "-svc"
	objectMeta := &metav1.ObjectMeta{
		Name: serviceName,
		Labels: map[string]string{
			"app":     serviceName,
			"service": serviceName,
		},
	}
	ports := []corev1.ServicePort{}
	ports = append(ports, corev1.ServicePort{
		Name: "http",
		Port: port,
	})
	spec := &corev1.ServiceSpec{
		Ports: ports,
		Selector: map[string]string{
			"app": metadata.Name, // be sure that same with deployment label `app: minicloud-api`
		},
	}
	service := &corev1.Service{
		TypeMeta:   *typeMeta,
		ObjectMeta: *objectMeta,
		Spec:       *spec,
	}
	acct, err := clientSet.CoreV1().Services(ns).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	return acct
}
