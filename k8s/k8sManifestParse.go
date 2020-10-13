package k8s

import (
	"fmt"
	"github.com/ghodss/yaml"
	"k8s.io/api/apps/v1beta1"
	"strings"
)

func (kr *k8sRunner) ManifestParse(job Job) (deployment *v1beta1.Deployment , err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in ManifestParse", r)
		}
	}()
	var key string

	if job.EnvName == "staging" {
		key  = "/staging/"
	} else {
		key = "/devteam/"
	}

	if  job.Source == "api" {
		key += "api-deployment.yaml"
	} else {
		key += "bff-deployment.yaml"
	}

	var content string
	content, err = kr.s3Client.GetManifestFromS3(key)

	if err != nil {
		return
	}

	if job.EnvName != "staging" {
		content = strings.ReplaceAll(content, "ENVNAME", job.EnvName)
		content = strings.ReplaceAll(content, "TAGNAME", job.ImageName)

		if job.Source == "bff" {
			if job.Ops != nil {
				if job.Ops["api"] == "" || strings.EqualFold(job.Ops["api"], "staging"){
					content = strings.ReplaceAll(content, "APIHOST", "homes-manager-kodate-api-svc")
				} else {
					content = strings.ReplaceAll(content, "APIHOST", fmt.Sprintf("%s-go-api-svc", job.Ops["api"]))
				}
			} else {
				content = strings.ReplaceAll(content, "APIHOST", "homes-manager-kodate-api-svc")
			}

			content = strings.ReplaceAll(content, "APIPORT", "80")
		}
	}

	err = yaml.Unmarshal([]byte(content), &deployment)

	return deployment, err
}
