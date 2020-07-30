package templates

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/src/pkg/obj/kubeobj"
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

func Install(opt *InstallOptions) error {

	kubeObjects, err := KubeObjectsFromTemplates(opt.Templates, opt.TemplateValues)
	if err != nil {
		return err
	}

	for _, obj := range kubeObjects {
		var createErr error
		var kind, name string
		name, kind, createErr = kubeobj.CreateObject(opt.KubeClientSet, obj, opt.Namespace)

		if createErr == nil {
			fmt.Println("%s \"%s\" created", kind, name)
		} else if statusError, errIsStatusError := createErr.(*errors.StatusError); errIsStatusError {
			if statusError.ErrStatus.Reason == metav1.StatusReasonAlreadyExists {
				fmt.Println("%s \"%s\" already exists", kind, name)
			} else {
				fmt.Println(fmt.Sprintf("%s \"%s\" failed: %v ", kind, name, statusError))
				return statusError
			}
		} else {
			fmt.Println("%s \"%s\" failed: %v ", kind, name, createErr)
			return createErr
		}
	}

	return nil
}
