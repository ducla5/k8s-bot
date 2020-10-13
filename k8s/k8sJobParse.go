package k8s

import (
	"errors"
	"fmt"
	chatwork "github.com/ducla5/go-chatwork"
	"strings"
)

// JobParse parse message to job
// sample: deploy bff image lftv-develop on staging
// sample: deploy api image lftv-develop on staging
func JobParse(message chatwork.Message, ) (job Job, err error)  {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in JobParse", r)
		}
	}()

	body := strings.Split(message.Body,"\n")

	body = strings.Split(body[1]," ")

	if len(body) != 6 && len(body) != 8 {
		err = errors.New("message invalid")
		return
	}

	if !strings.EqualFold(body[0],"deploy") || !strings.EqualFold(body[4],"on") {
		err = errors.New("message invalid")
		return
	}

	job.Action = strings.ToLower(body[0])

	if !strings.EqualFold(body[1],"bff") && !strings.EqualFold(body[1],"api") {
		err = errors.New("please specify: bff or api")
		return
	}

	job.Source = strings.ToLower(body[1])

	if !strings.EqualFold(body[2],"image") && !strings.EqualFold(body[2],"branch") {
		err = errors.New("please specify: image or branch name")
		return
	}

	if strings.EqualFold(body[2],"image") {
		job.ImageName = strings.ToLower(body[3])
	} else {
		job.BranchName = strings.ToLower(body[3])

	}

	job.EnvName = strings.ToLower(body[5])

	if strings.EqualFold(body[1],"bff") && len(body) > 6 && strings.EqualFold(body[6], "api") {
		job.Ops = map[string]string{"api":strings.ToLower(body[7])}
	}

	job.Pic = message.Account.Name
	job.Date = message.SendTime

	return
}
