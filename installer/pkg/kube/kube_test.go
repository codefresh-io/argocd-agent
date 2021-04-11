package kube

import (
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var _ = func() bool {
	testing.Init()
	return true
}()

func TestGetLoadBalancerHost(t *testing.T) {
	kube, _ := New(&Options{
		ContextName:      "",
		Namespace:        "",
		PathToKubeConfig: "",
		InCluster:        false,
		FailFast:         true,
	})

	ingresses := make([]v1.LoadBalancerIngress, 0)
	ingresses = append(ingresses, v1.LoadBalancerIngress{
		IP:       "ip",
		Hostname: "host",
	})

	host, err := kube.GetLoadBalancerHost(v1.Service{
		TypeMeta:   v12.TypeMeta{},
		ObjectMeta: v12.ObjectMeta{},
		Spec:       v1.ServiceSpec{},
		Status: v1.ServiceStatus{
			LoadBalancer: v1.LoadBalancerStatus{
				Ingress: ingresses,
			},
		},
	})

	if err != nil {
		t.Errorf("Should be executed without error, got error %v", err.Error())
	}

	if host != "https://host" {
		t.Errorf("Host should be \"https://host\" but actually got %s", host)
	}
}

func TestGetLoadBalancerHostByIp(t *testing.T) {
	kube, _ := New(&Options{
		ContextName:      "",
		Namespace:        "",
		PathToKubeConfig: "",
		InCluster:        false,
		FailFast:         true,
	})

	ingresses := make([]v1.LoadBalancerIngress, 0)
	ingresses = append(ingresses, v1.LoadBalancerIngress{
		IP:       "ip",
		Hostname: "",
	})

	host, err := kube.GetLoadBalancerHost(v1.Service{
		TypeMeta:   v12.TypeMeta{},
		ObjectMeta: v12.ObjectMeta{},
		Spec:       v1.ServiceSpec{},
		Status: v1.ServiceStatus{
			LoadBalancer: v1.LoadBalancerStatus{
				Ingress: ingresses,
			},
		},
	})

	if err != nil {
		t.Errorf("Should be executed without error, got error %v", err.Error())
	}

	if host != "https://ip" {
		t.Errorf("Host should be \"https://ip\" but actually got %s", host)
	}
}

func TestGetLoadBalancerHostWithoutHostAndIp(t *testing.T) {
	kube, _ := New(&Options{
		ContextName:      "",
		Namespace:        "",
		PathToKubeConfig: "",
		InCluster:        false,
		FailFast:         true,
	})

	ingresses := make([]v1.LoadBalancerIngress, 0)
	ingresses = append(ingresses, v1.LoadBalancerIngress{
		IP:       "",
		Hostname: "",
	})

	_, err := kube.GetLoadBalancerHost(v1.Service{
		TypeMeta:   v12.TypeMeta{},
		ObjectMeta: v12.ObjectMeta{},
		Spec:       v1.ServiceSpec{},
		Status: v1.ServiceStatus{
			LoadBalancer: v1.LoadBalancerStatus{
				Ingress: ingresses,
			},
		},
	})

	if err == nil || err.Error() != "Failed to retrieve Load Balancer Hostname or IP" {
		t.Errorf("Should be executed with error")
	}

}

func TestGetLoadBalancerHostWithoutIngress(t *testing.T) {
	kube, _ := New(&Options{
		ContextName:      "",
		Namespace:        "",
		PathToKubeConfig: "",
		InCluster:        false,
		FailFast:         true,
	})

	ingresses := make([]v1.LoadBalancerIngress, 0)

	_, err := kube.GetLoadBalancerHost(v1.Service{
		TypeMeta:   v12.TypeMeta{},
		ObjectMeta: v12.ObjectMeta{},
		Spec:       v1.ServiceSpec{},
		Status: v1.ServiceStatus{
			LoadBalancer: v1.LoadBalancerStatus{
				Ingress: ingresses,
			},
		},
	})

	if err == nil {
		t.Errorf("Should be executed with error")
	}
}
