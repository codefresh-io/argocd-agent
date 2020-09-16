package extract

import (
	"fmt"
	"github.com/codefresh-io/argocd-listener/agent/pkg/argo"
	codefresh2 "github.com/codefresh-io/argocd-listener/agent/pkg/codefresh"
	"github.com/codefresh-io/argocd-listener/agent/pkg/handler"
	"github.com/codefresh-io/argocd-listener/agent/pkg/transform"
	"github.com/codefresh-io/argocd-listener/agent/pkg/util"
	"github.com/golang/glog"
	"github.com/mitchellh/mapstructure"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	applicationCRD = schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "applications",
	}

	projectCRD = schema.GroupVersionResource{
		Group:    "argoproj.io",
		Version:  "v1alpha1",
		Resource: "appprojects",
	}
)

func buildConfig() *rest.Config {
	inCluster, _ := strconv.ParseBool(os.Getenv("IN_CLUSTER"))
	if inCluster {
		config, _ := rest.InClusterConfig()
		return config
	} else {
		kubeconfig := filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatal(err)
		}
		return config
	}
}
func updateEnv(obj interface{}) error {
	err, env := transform.PrepareEnvironment(obj.(*unstructured.Unstructured).Object)
	if err != nil {
		fmt.Println(fmt.Sprintf("Cant preapre env for codefresh because %v", err))
		return err
	}

	err = util.ProcessDataWithFilter("environment", env, func() error {
		_, err = codefresh2.GetInstance().SendEnvironment(*env)
		return err
	})

	return nil
}

func watchApplicationChanges() {
	clientset, err := dynamic.NewForConfig(buildConfig())
	if err != nil {
		glog.Errorln(err)
	}

	kubeInformerFactory := dynamicinformer.NewDynamicSharedInformerFactory(clientset, time.Minute*30)
	applicationInformer := kubeInformerFactory.ForResource(applicationCRD).Informer()

	api := codefresh2.GetInstance()

	applicationInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			var app argo.ArgoApplication
			err := mapstructure.Decode(obj.(*unstructured.Unstructured).Object, &app)

			err = updateEnv(obj)

			if err != nil {
				fmt.Println(fmt.Sprintf("Cant send env to codefresh because %v", err))
			}

			applications := argo.GetApplications()
			err = util.ProcessDataWithFilter("applications", applications, func() error {
				return api.SendResources("applications", transform.AdaptArgoApplications(applications))
			})

			applicationCreatedHandler := handler.GetApplicationCreatedHandlerInstance()
			err = applicationCreatedHandler.Handle(app)

			if err != nil {
				fmt.Print(err)
			}

			log.Println(applications)
		},
		DeleteFunc: func(obj interface{}) {
			var app argo.ArgoApplication
			err := mapstructure.Decode(obj.(*unstructured.Unstructured).Object, &app)
			if err != nil {
				fmt.Print(err)
			}

			applications := argo.GetApplications()
			err = util.ProcessDataWithFilter("applications", applications, func() error {
				return api.SendResources("applications", transform.AdaptArgoApplications(applications))
			})
			if err != nil {
				fmt.Print(err)
			}

			applicationRemovedHandler := handler.GetApplicationRemovedHandlerInstance()
			err = applicationRemovedHandler.Handle(app)

			if err != nil {
				fmt.Print(err)
			}

		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			err := updateEnv(newObj)
			if err != nil {
				fmt.Println(fmt.Sprintf("Cant send env to codefresh because %v", err))
			}
		},
	})

	projectInformer := kubeInformerFactory.ForResource(projectCRD).Informer()

	projectInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Printf("project added: %s \n", obj)
			projects := argo.GetProjects()

			err := util.ProcessDataWithFilter("projects", projects, func() error {
				return api.SendResources("projects", transform.AdaptArgoProjects(projects))
			})

			if err != nil {
				fmt.Print(err)
			}
			fmt.Println(projects)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Printf("project deleted: %s \n", obj)
			projects := argo.GetProjects()
			err := util.ProcessDataWithFilter("projects", projects, func() error {
				return api.SendResources("projects", transform.AdaptArgoProjects(projects))
			})
			if err != nil {
				fmt.Print(err)
			}
			fmt.Println(projects)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Printf("project updated: %s \n", newObj)
		},
	})

	stop := make(chan struct{})
	defer close(stop)
	kubeInformerFactory.Start(stop)

	for {
		time.Sleep(time.Second)
	}

}

func Watch() {
	watchApplicationChanges()

}
