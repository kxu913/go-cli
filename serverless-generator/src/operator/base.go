package operator

import (
	"log"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

func DynamicClient() *dynamic.DynamicClient {
	config, err := clientcmd.BuildConfigFromFlags("", ConfigFile)
	if err != nil {
		panic(err)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return client
}

func ClientSet() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", ConfigFile)
	if err != nil {
		panic(err)
	}
	config.APIPath = "api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs

	clientSet, err := kubernetes.NewForConfig(config)
	return clientSet
}

func IstioClient() *versionedclient.Clientset {
	restConfig, err := clientcmd.BuildConfigFromFlags("", ConfigFile)
	if err != nil {
		log.Fatalf("Failed to create k8s rest client: %s", err)
	}
	ic, err := versionedclient.NewForConfig(restConfig)
	return ic
}
