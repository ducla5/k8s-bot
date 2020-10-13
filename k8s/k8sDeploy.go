package k8s

import (
	"context"
	"flag"
	"fmt"
	"k8s-bot/config"
	"k8s-bot/s3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type k8sRunner struct {
	clientSet *kubernetes.Clientset
	s3Client  s3.Util
}

type Runner interface {
	ProcessJob(job Job) error
	ProcessJobByManifest(job Job) error
}

func NewK8sRunner() Runner {
	configK8s, err := rest.InClusterConfig()

	if err != nil {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		configK8s, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	clientSet, err := kubernetes.NewForConfig(configK8s)

	if err != nil {
		panic(err.Error())
	}

	s3Client := s3.NewS3Client()
	return &k8sRunner{clientSet: clientSet, s3Client: s3Client}
}

func (kr *k8sRunner) ProcessJob(job Job) (err error) {
	var deploymentName string
	var imageName1 string
	var imageName2 string

	if strings.EqualFold(job.Source, "api") {
		if strings.EqualFold(job.EnvName, "staging") {
			deploymentName = "homes-manager-kodate-api"
		} else {
			deploymentName = fmt.Sprintf("%s-new-manager-go", job.EnvName)
		}

		imageName1 = fmt.Sprintf("%s:%s", config.C.ECR.Api, job.ImageName)

	} else {
		if strings.EqualFold(job.EnvName, "staging") {
			deploymentName = "homes-manager-kodate"
		} else {
			deploymentName = fmt.Sprintf("%s-new-manager-bff", job.EnvName)
		}

		imageName1 = fmt.Sprintf("%s:%s", config.C.ECR.Bff, job.ImageName)
		imageName2 = fmt.Sprintf("%s:%s", config.C.ECR.BffNginx, job.ImageName)
	}

	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {

		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getError := kr.clientSet.AppsV1().Deployments(config.C.K8s.NameSpace).Get(context.TODO(), deploymentName, metav1.GetOptions{})

		if getError != nil {
			return getError
		}

		if result.Spec.Template.Spec.Containers[0].Image == imageName1 {
			selector := result.ObjectMeta.Labels["app"]
			result, err := kr.clientSet.CoreV1().Pods(config.C.K8s.NameSpace).List(context.TODO(), metav1.ListOptions{LabelSelector: fmt.Sprintf("app=%s", selector)})

			if err != nil {
				return err
			}

			for _, v := range result.Items {
				err = kr.clientSet.CoreV1().Pods(config.C.K8s.NameSpace).Delete(context.TODO(), v.Name, metav1.DeleteOptions{})
			}

			return err
		}

		// Change image name
		result.Spec.Template.Spec.Containers[0].Image = imageName1
		if strings.EqualFold(job.Source, "bff") {
			result.Spec.Template.Spec.Containers[1].Image = imageName2
		}

		_, updateErr := kr.clientSet.AppsV1().Deployments(config.C.K8s.NameSpace).Update(context.TODO(), result, metav1.UpdateOptions{})

		return updateErr
	})
	return
}

func (kr *k8sRunner) ProcessJobByManifest(job Job) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in ProcessJobByManifest", r)
		}
	}()
	dep, err := kr.ManifestParse(job)

	if err != nil {
		log.Println(err)
	}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}

	err = kr.clientSet.AppsV1beta1().Deployments(config.C.K8s.NameSpace).Delete(context.TODO(), dep.Name, deleteOptions)

	if err != nil && !errors.IsNotFound(err) {
		return
	}

	maxTry := 5

	for i := 0; i < maxTry; i++ {
		_, err = kr.clientSet.AppsV1beta1().Deployments(config.C.K8s.NameSpace).Create(context.TODO(), dep, metav1.CreateOptions{})

		if err == nil || !errors.IsAlreadyExists(err) {
			break
		}

		time.Sleep(10 * time.Second)
	}

	return
}
