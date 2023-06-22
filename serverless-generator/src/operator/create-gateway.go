package operator

import (
	"context"
	"serverless-generator/model"

	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateVirtualService(ns string, metadata *model.MetaData, port int32) *v1alpha3.VirtualService {
	client := IstioClient()

	typeMeta := &metav1.TypeMeta{
		Kind: "VirtualService",
	}
	gatewayName := metadata.Name + "-vs"
	objectMeta := &metav1.ObjectMeta{
		Name: gatewayName,
		Labels: map[string]string{
			"name": metadata.Name,
		},
	}
	uri := &networkingv1alpha3.StringMatch{
		MatchType: &networkingv1alpha3.StringMatch_Prefix{
			Prefix: metadata.Prefix,
		},
	}
	match := &networkingv1alpha3.HTTPMatchRequest{
		Uri: uri,
	}
	destination := &networkingv1alpha3.HTTPRouteDestination{
		Destination: &networkingv1alpha3.Destination{
			Host: metadata.Name + "-svc",
			Port: &networkingv1alpha3.PortSelector{
				Number: uint32(port),
			},
		},
	}
	route := &networkingv1alpha3.HTTPRoute{
		Match: []*networkingv1alpha3.HTTPMatchRequest{match},
		Route: []*networkingv1alpha3.HTTPRouteDestination{destination},
	}

	service := &v1alpha3.VirtualService{
		TypeMeta:   *typeMeta,
		ObjectMeta: *objectMeta,
		Spec: networkingv1alpha3.VirtualService{
			Hosts:    []string{"kxu.kevin.com"},
			Gateways: []string{ns + "-gw"},
			Http:     []*networkingv1alpha3.HTTPRoute{route},
		},
	}
	vs, err := client.NetworkingV1alpha3().VirtualServices(ns).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	return vs
}

func CreateGateway(ns string) *v1alpha3.Gateway {
	client := IstioClient()
	typeMeta := &metav1.TypeMeta{
		Kind: "Gateway",
	}
	gwName := ns + "-gw"
	objectMeta := &metav1.ObjectMeta{
		Name: gwName,
		Labels: map[string]string{
			"name": gwName,
		},
	}

	server := &networkingv1alpha3.Server{
		Port: &networkingv1alpha3.Port{
			Number:   80,
			Name:     "http",
			Protocol: "HTTP",
		},
		Hosts: []string{
			"kxu.kevin.com",
		},
	}

	gateway := &v1alpha3.Gateway{
		TypeMeta:   *typeMeta,
		ObjectMeta: *objectMeta,
		Spec: networkingv1alpha3.Gateway{
			Selector: map[string]string{
				"istio": "ingressgateway",
			},
			Servers: []*networkingv1alpha3.Server{server},
		},
	}
	result, err := client.NetworkingV1alpha3().Gateways(ns).Create(context.TODO(), gateway, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	return result
}
