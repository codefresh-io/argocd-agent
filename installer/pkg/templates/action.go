package templates

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/installer/pkg/fs"
	"github.com/codefresh-io/argocd-listener/installer/pkg/logger"
	kubeobj "github.com/codefresh-io/argocd-listener/installer/pkg/obj/kubeobj"
	apixv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"reflect"
)

type InstallOptions struct {
	Templates        map[string]string
	TemplateValues   map[string]interface{}
	KubeClientSet    *kubernetes.Clientset
	KubeCrdClientSet *apixv1beta1client.ApiextensionsV1beta1Client
	Namespace        string
	KubeManifestPath string
	KubeBuilder      interface {
		BuildClient() (*kubernetes.Clientset, error)
	}
}

type DeleteOptions struct {
	Templates        map[string]string
	TemplateValues   map[string]interface{}
	KubeClientSet    *kubernetes.Clientset
	KubeCrdClientSet *apixv1beta1client.ApiextensionsV1beta1Client
	Namespace        string
	KubeBuilder      interface {
		BuildClient() (*kubernetes.Clientset, error)
	}
}

func Install(opt *InstallOptions) (error, string, string) {
	opt.TemplateValues["Namespace"] = opt.Namespace
	kubeObjects, parsedTemplates, err := KubeObjectsFromTemplates(opt.Templates, opt.TemplateValues)
	if err != nil {
		return err, "", ""
	}

	if opt.KubeManifestPath != "" {
		manifest := GenerateSingleManifest(parsedTemplates)
		err = fs.WriteFile(opt.KubeManifestPath, manifest)
		if err != nil {
			return err, "", ""
		}
	}

	kubeObjectKeys := reflect.ValueOf(kubeObjects).MapKeys()

	for _, key := range kubeObjectKeys {
		kind, name, createErr := kubeobj.CreateObject(opt.KubeClientSet, opt.KubeCrdClientSet, kubeObjects[key.String()], opt.Namespace)

		if createErr == nil {
			// skip, everything ok
		} else if statusError, errIsStatusError := createErr.(*errors.StatusError); errIsStatusError {
			if statusError.ErrStatus.Reason == metav1.StatusReasonAlreadyExists {
				logger.Warning(fmt.Sprintf("%s \"%s\" already exists", kind, name))
			} else {
				logger.Error(fmt.Sprintf("%s \"%s\" failed: %v ", kind, name, statusError))
				return statusError, kind, name
			}
		} else {
			logger.Error(fmt.Sprintf("%s \"%s\" failed: %v ", kind, name, createErr))
			return createErr, kind, name
		}
	}

	return nil, "", ""
}

func Delete(opt *DeleteOptions) (error, string, string) {

	kubeObjects, _, err := KubeObjectsFromTemplates(opt.Templates, opt.TemplateValues)
	if err != nil {
		return err, "", ""
	}
	var kind, name string
	var deleteError error
	for _, obj := range kubeObjects {
		kind, name, deleteError = kubeobj.DeleteObject(opt.KubeClientSet, opt.KubeCrdClientSet, obj, opt.Namespace)
		if deleteError == nil {
			fmt.Println(fmt.Sprintf("%s \"%s\" deleted", kind, name))
		} else if statusError, errIsStatusError := deleteError.(*errors.StatusError); errIsStatusError {
			if statusError.ErrStatus.Reason == metav1.StatusReasonNotFound {
				logger.Warning(fmt.Sprintf("Resource %s \"%s\" not found", kind, name))
			} else {
				return statusError, kind, name
			}
		} else {
			return deleteError, kind, name
		}
	}
	return nil, "", ""
}
