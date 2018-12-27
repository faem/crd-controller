package main

import (
	clientset "crd-controller/pkg/client/clientset/versioned"
	informers "crd-controller/pkg/client/informers/externalversions"
	cntrlr "crd-controller/pkg/controller"
	"fmt"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"time"
)

var (
	masterURL  string
	kubeconfig string
)

func main() {
	kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube/config")

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		fmt.Println(err.Error())
	}

	kubeClient := kubernetes.NewForConfigOrDie(cfg)
	exampleClient := clientset.NewForConfigOrDie(cfg)

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	exampleInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*30)

	controller := cntrlr.NewController(kubeClient, exampleClient,
		kubeInformerFactory.Apps().V1().Deployments(),
		exampleInformerFactory.Crdcntrlr().V1alpha1().Foos())

	stopCh := make(chan struct{})
	kubeInformerFactory.Start(stopCh)
	exampleInformerFactory.Start(stopCh)

	if err = controller.Run(stopCh); err != nil {
		fmt.Println("Error running controller: ", err.Error())
	}

}
