// Code generated by go generate; DO NOT EDIT.

package kubeobj

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"

	appsv1 "k8s.io/api/apps/v1"
	v1api "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	apiextensionv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"

	storagev1 "k8s.io/api/storage/v1"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
)

// CreateObject - creates kubernetes object from *runtime.Object. Returns object name, kind and creation error
func CreateObject(clientset *kubernetes.Clientset, apiextensionsClientSet *apixv1beta1client.ApiextensionsV1beta1Client, obj runtime.Object, namespace string) (string, string, error) {

	var name, kind string
	var err error
	switch objT := obj.(type) {

	case *appsv1.DaemonSet:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.AppsV1().DaemonSets(namespace).Create(objT)

	case *appsv1.Deployment:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.AppsV1().Deployments(namespace).Create(objT)

	case *batchv1.Job:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.BatchV1().Jobs(namespace).Create(objT)

	case *batchv1beta1.CronJob:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.BatchV1beta1().CronJobs(namespace).Create(objT)

	case *rbacv1.ClusterRole:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1().ClusterRoles().Create(objT)

	case *rbacv1.ClusterRoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1().ClusterRoleBindings().Create(objT)

	case *rbacv1.Role:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1().Roles(namespace).Create(objT)

	case *rbacv1.RoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1().RoleBindings(namespace).Create(objT)

	case *rbacv1beta1.ClusterRole:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1beta1().ClusterRoles().Create(objT)

	case *rbacv1beta1.ClusterRoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1beta1().ClusterRoleBindings().Create(objT)

	case *rbacv1beta1.Role:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1beta1().Roles(namespace).Create(objT)

	case *rbacv1beta1.RoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1beta1().RoleBindings(namespace).Create(objT)

	case *storagev1.StorageClass:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.StorageV1().StorageClasses().Create(objT)

	case *v1.ConfigMap:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().ConfigMaps(namespace).Create(objT)

	case *v1.PersistentVolume:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().PersistentVolumes().Create(objT)

	case *v1.PersistentVolumeClaim:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().PersistentVolumeClaims(namespace).Create(objT)

	case *v1.Pod:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().Pods(namespace).Create(objT)

	case *v1.Secret:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().Secrets(namespace).Create(objT)

	case *v1.Service:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().Services(namespace).Create(objT)

	case *v1.ServiceAccount:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().ServiceAccounts(namespace).Create(objT)

	case *v1beta1.DaemonSet:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.ExtensionsV1beta1().DaemonSets(namespace).Create(objT)

	case *v1beta1.Deployment:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.ExtensionsV1beta1().Deployments(namespace).Create(objT)

	case *netv1.NetworkPolicy:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.NetworkingV1().NetworkPolicies(namespace).Create(objT)

	case *v1api.StatefulSet:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.AppsV1().StatefulSets(namespace).Create(objT)

	case *apiextensionv1beta1.CustomResourceDefinition:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = apiextensionsClientSet.CustomResourceDefinitions().Create(objT)

	default:
		return "", "", fmt.Errorf("Unknown object type %T\n ", objT)
	}
	return name, kind, err
}

// CheckObject - checks kubernetes object from *runtime.Object. Returns object name, kind and creation error
func CheckObject(clientset *kubernetes.Clientset, obj runtime.Object, namespace string) (string, string, error) {

	var name, kind string
	var err error
	switch objT := obj.(type) {

	case *appsv1.DaemonSet:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.AppsV1().DaemonSets(namespace).Get(name, metav1.GetOptions{})

	case *appsv1.Deployment:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})

	case *batchv1.Job:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.BatchV1().Jobs(namespace).Get(name, metav1.GetOptions{})

	case *batchv1beta1.CronJob:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.BatchV1beta1().CronJobs(namespace).Get(name, metav1.GetOptions{})

	case *rbacv1.ClusterRole:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1().ClusterRoles().Get(name, metav1.GetOptions{})

	case *rbacv1.ClusterRoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1().ClusterRoleBindings().Get(name, metav1.GetOptions{})

	case *rbacv1.Role:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1().Roles(namespace).Get(name, metav1.GetOptions{})

	case *rbacv1.RoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1().RoleBindings(namespace).Get(name, metav1.GetOptions{})

	case *rbacv1beta1.ClusterRole:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1beta1().ClusterRoles().Get(name, metav1.GetOptions{})

	case *rbacv1beta1.ClusterRoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1beta1().ClusterRoleBindings().Get(name, metav1.GetOptions{})

	case *rbacv1beta1.Role:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1beta1().Roles(namespace).Get(name, metav1.GetOptions{})

	case *rbacv1beta1.RoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1beta1().RoleBindings(namespace).Get(name, metav1.GetOptions{})

	case *storagev1.StorageClass:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.StorageV1().StorageClasses().Get(name, metav1.GetOptions{})

	case *v1.ConfigMap:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().ConfigMaps(namespace).Get(name, metav1.GetOptions{})

	case *v1.PersistentVolume:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().PersistentVolumes().Get(name, metav1.GetOptions{})

	case *v1.PersistentVolumeClaim:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().PersistentVolumeClaims(namespace).Get(name, metav1.GetOptions{})

	case *v1.Pod:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})

	case *v1.Secret:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})

	case *v1.Service:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})

	case *v1.ServiceAccount:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().ServiceAccounts(namespace).Get(name, metav1.GetOptions{})

	case *v1beta1.DaemonSet:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.ExtensionsV1beta1().DaemonSets(namespace).Get(name, metav1.GetOptions{})

	case *v1beta1.Deployment:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.ExtensionsV1beta1().Deployments(namespace).Get(name, metav1.GetOptions{})

	case *netv1.NetworkPolicy:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.NetworkingV1().NetworkPolicies(namespace).Get(name, metav1.GetOptions{})

	default:
		return "", "", fmt.Errorf("Unknown object type %T\n ", objT)
	}
	return name, kind, err
}

// DeleteObject - checks kubernetes object from *runtime.Object. Returns object name, kind and creation error
func DeleteObject(clientset *kubernetes.Clientset, apiextensionsClientSet *apixv1beta1client.ApiextensionsV1beta1Client, obj runtime.Object, namespace string) (string, string, error) {
	var propagationPolicy metav1.DeletionPropagation = "Background"
	var name, kind string
	var err error
	switch objT := obj.(type) {

	case *apiextensionv1beta1.CustomResourceDefinition:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = apiextensionsClientSet.CustomResourceDefinitions().Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *v1api.StatefulSet:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.AppsV1().StatefulSets(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *appsv1.DaemonSet:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.AppsV1().DaemonSets(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *appsv1.Deployment:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.AppsV1().Deployments(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *batchv1.Job:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.BatchV1().Jobs(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *batchv1beta1.CronJob:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.BatchV1beta1().CronJobs(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *rbacv1.ClusterRole:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.RbacV1().ClusterRoles().Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *rbacv1.ClusterRoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.RbacV1().ClusterRoleBindings().Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *rbacv1.Role:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.RbacV1().Roles(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *rbacv1.RoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.RbacV1().RoleBindings(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *rbacv1beta1.ClusterRole:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.RbacV1beta1().ClusterRoles().Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *rbacv1beta1.ClusterRoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.RbacV1beta1().ClusterRoleBindings().Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *rbacv1beta1.Role:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.RbacV1beta1().Roles(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *rbacv1beta1.RoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.RbacV1beta1().RoleBindings(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *storagev1.StorageClass:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.StorageV1().StorageClasses().Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *v1.ConfigMap:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.CoreV1().ConfigMaps(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *v1.PersistentVolume:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.CoreV1().PersistentVolumes().Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *v1.PersistentVolumeClaim:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *v1.Pod:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.CoreV1().Pods(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *v1.Secret:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.CoreV1().Secrets(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *v1.Service:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.CoreV1().Services(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *v1.ServiceAccount:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.CoreV1().ServiceAccounts(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *v1beta1.DaemonSet:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.ExtensionsV1beta1().DaemonSets(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *v1beta1.Deployment:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.ExtensionsV1beta1().Deployments(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	case *netv1.NetworkPolicy:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		err = clientset.NetworkingV1().NetworkPolicies(namespace).Delete(name, &metav1.DeleteOptions{
			PropagationPolicy: &propagationPolicy,
		})

	default:
		return "", "", fmt.Errorf("Unknown object type %T\n ", objT)
	}
	return name, kind, err
}

// ReplaceObject - replaces kubernetes object from *runtime.Object. Returns object name, kind and creation error
func ReplaceObject(clientset *kubernetes.Clientset, obj runtime.Object, namespace string) (string, string, error) {
	var name, kind string
	var err error
	switch objT := obj.(type) {

	case *appsv1.DaemonSet:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.AppsV1().DaemonSets(namespace).Update(objT)

	case *appsv1.Deployment:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.AppsV1().Deployments(namespace).Update(objT)

	case *batchv1.Job:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.BatchV1().Jobs(namespace).Update(objT)

	case *batchv1beta1.CronJob:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.BatchV1beta1().CronJobs(namespace).Update(objT)

	case *rbacv1.ClusterRole:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1().ClusterRoles().Update(objT)

	case *rbacv1.ClusterRoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1().ClusterRoleBindings().Update(objT)

	case *rbacv1.Role:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1().Roles(namespace).Update(objT)

	case *rbacv1.RoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1().RoleBindings(namespace).Update(objT)

	case *rbacv1beta1.ClusterRole:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1beta1().ClusterRoles().Update(objT)

	case *rbacv1beta1.ClusterRoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1beta1().ClusterRoleBindings().Update(objT)

	case *rbacv1beta1.Role:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1beta1().Roles(namespace).Update(objT)

	case *rbacv1beta1.RoleBinding:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.RbacV1beta1().RoleBindings(namespace).Update(objT)

	case *storagev1.StorageClass:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.StorageV1().StorageClasses().Update(objT)

	case *v1.ConfigMap:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().ConfigMaps(namespace).Update(objT)

	case *v1.PersistentVolume:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().PersistentVolumes().Update(objT)

	case *v1.PersistentVolumeClaim:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().PersistentVolumeClaims(namespace).Update(objT)

	case *v1.Pod:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().Pods(namespace).Update(objT)

	case *v1.Secret:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().Secrets(namespace).Update(objT)

	case *v1.Service:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().Services(namespace).Update(objT)

	case *v1.ServiceAccount:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.CoreV1().ServiceAccounts(namespace).Update(objT)

	case *v1beta1.DaemonSet:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.ExtensionsV1beta1().DaemonSets(namespace).Update(objT)

	case *v1beta1.Deployment:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.ExtensionsV1beta1().Deployments(namespace).Update(objT)

	case *netv1.NetworkPolicy:
		name = objT.ObjectMeta.Name
		kind = objT.TypeMeta.Kind
		_, err = clientset.NetworkingV1().NetworkPolicies(namespace).Update(objT)

	default:
		return "", "", fmt.Errorf("Unknown object type %T\n ", objT)
	}
	return name, kind, err
}

func DeletePod(clientset *kubernetes.Clientset, namespace string, labelSelector string) error {
	var propagationPolicy metav1.DeletionPropagation = "Background"

	listOptions := metav1.ListOptions{LabelSelector: labelSelector}
	return clientset.CoreV1().Pods(namespace).DeleteCollection(&metav1.DeleteOptions{PropagationPolicy: &propagationPolicy}, listOptions)
}

func GetDeployments(clientset *kubernetes.Clientset, namespace string, labelSelector string) (*appsv1.DeploymentList, error) {
	listOptions := metav1.ListOptions{LabelSelector: labelSelector}
	return clientset.AppsV1().Deployments(namespace).List(listOptions)
}

func UpdateDeployment(clientset *kubernetes.Clientset, deployment *v1api.Deployment, namespace string) (*v1api.Deployment, error) {
	return clientset.AppsV1().Deployments(namespace).Update(deployment)
}
