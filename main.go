package main

import (
	"container/list"
	"errors"
	chatwork2 "k8s-bot/chatwork"
	"k8s-bot/config"
	"k8s-bot/fun"
	"k8s-bot/k8s"
	"log"
	"time"
)

func main() {
	config.LoadConfig()
	queue := list.New()
	cw := chatwork2.NewChatworkUtil(queue)
	k8sClient := k8s.NewK8sRunner()
	tick := time.Tick(5000 * time.Millisecond)
	tick2 := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-tick:
			cw.FetchJob()
		case <-tick2:
			mess := cw.GetJob()
			if mess.Body != "" {
				job, err := k8s.JobParse(mess)

				if err != nil {
					log.Println(err)
					cw.ErrorReport(errors.New(fun.GetQuote()), mess)
					break
				}

				err = k8sClient.ProcessJobByManifest(job)

				if err != nil {
					log.Println(err)
					cw.ErrorReport(err, mess)
					break
				}

				time.Sleep(20 * time.Second)
				cw.ResultReport(job, mess)
			}
		default:
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
