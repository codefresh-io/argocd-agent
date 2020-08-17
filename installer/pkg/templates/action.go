package templates

import (
	"fmt"
	kubeobj "github.com/codefresh-io/argocd-listener/installer/pkg/obj/kubeobj"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type InstallOptions struct {
	Templates      map[string]string
	TemplateValues map[string]interface{}
	KubeClientSet  *kubernetes.Clientset
	Namespace      string
	KubeBuilder    interface {
		BuildClient() (*kubernetes.Clientset, error)
	}
}

type DeleteOptions struct {
	Templates      map[string]string
	TemplateValues map[string]interface{}
	KubeClientSet  *kubernetes.Clientset
	Namespace      string
	KubeBuilder    interface {
		BuildClient() (*kubernetes.Clientset, error)
	}
}

func Install(opt *InstallOptions) error {

	opt.TemplateValues["Namespace"] = opt.Namespace
	kubeObjects, err := KubeObjectsFromTemplates(opt.Templates, opt.TemplateValues)
	if err != nil {
		return err
	}

	for _, obj := range kubeObjects {
		var createErr error
		var kind, name string
		name, kind, createErr = kubeobj.CreateObject(opt.KubeClientSet, obj, opt.Namespace)

		if createErr == nil {
			fmt.Println(fmt.Sprintf("%s \"%s\" created", kind, name))
		} else if statusError, errIsStatusError := createErr.(*errors.StatusError); errIsStatusError {
			if statusError.ErrStatus.Reason == metav1.StatusReasonAlreadyExists {
				fmt.Println(fmt.Sprintf("%s \"%s\" already exists", kind, name))
			} else {
				fmt.Println(fmt.Sprintf("%s \"%s\" failed: %v ", kind, name, statusError))
				return statusError
			}
		} else {
			fmt.Println(fmt.Sprintf("%s \"%s\" failed: %v ", kind, name, createErr))
			return createErr
		}
	}

	return nil
}

func Delete(opt *DeleteOptions) error {

	kubeObjects, err := KubeObjectsFromTemplates(opt.Templates, opt.TemplateValues)
	if err != nil {
		return err
	}
	var kind, name string
	var deleteError error
	for _, obj := range kubeObjects {
		kind, name, deleteError = kubeobj.DeleteObject(opt.KubeClientSet, obj, opt.Namespace)
		if deleteError == nil {
			fmt.Println(fmt.Sprintf("%s \"%s\" deleted", kind, name))
		} else if statusError, errIsStatusError := deleteError.(*errors.StatusError); errIsStatusError {
			if statusError.ErrStatus.Reason == metav1.StatusReasonAlreadyExists {
				fmt.Println(fmt.Sprintf("%s \"%s\" already exist", kind, name))
			} else if statusError.ErrStatus.Reason == metav1.StatusReasonNotFound {
				fmt.Println(fmt.Sprintf("%s \"%s\" not found", kind, name))
			} else {
				fmt.Println(fmt.Sprintf("%s \"%s\" failed: %v ", kind, name, statusError))
				return statusError
			}
		} else {
			fmt.Println(fmt.Sprintf("%s \"%s\" failed: %v ", kind, name, deleteError))
			return deleteError
		}
	}
	return nil
}
