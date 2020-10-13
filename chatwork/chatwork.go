package chatwork

import (
	"container/list"
	"fmt"
	chatwork "github.com/ducla5/go-chatwork"
	"k8s-bot/config"
	"k8s-bot/k8s"
	"strings"
)

type util struct {
	queue *list.List
	cw *chatwork.Client
	me chatwork.Me
}

type Util interface {
	FetchJob()
	GetJob() chatwork.Message
	ErrorReport(error, chatwork.Message)
	ResultReport(k8s.Job, chatwork.Message)
}

func NewChatworkUtil(queue *list.List) Util {
	client := &util{queue: queue, cw: chatwork.NewClient(config.C.ChatWork.APIKey)}
	client.me = client.cw.Me()
	return client
}

func (u *util) FetchJob() {
	messages := u.cw.RoomMessages(config.C.ChatWork.RoomID)
	for _, v := range messages {
		if strings.Contains(v.Body, fmt.Sprintf("[To:%d]", u.me.AccountId)){
			u.queue.PushBack(v)
		}
	}

}

func (u *util) GetJob() chatwork.Message {
	if u.queue.Len() == 0 {
		return chatwork.Message{}
	}

	job := u.queue.Front()
	u.queue.Remove(job)

	return job.Value.(chatwork.Message)
}

func (u *util) ErrorReport(err error, message chatwork.Message)  {
	messageBody := fmt.Sprintf("[rp aid=%d to=%s-%s]%s \n%s", message.Account.AccountId, config.C.ChatWork.RoomID, message.MessageId, message.Account.Name, err.Error())
	u.cw.PostRoomMessage(config.C.ChatWork.RoomID, messageBody)
}

func (u *util) ResultReport(job k8s.Job, message chatwork.Message)  {
	content := fmt.Sprintf("%s env %s has been updated", job.Source, job.EnvName)
	messageBody := fmt.Sprintf("[rp aid=%d to=%s-%s]%s \n%s", message.Account.AccountId, config.C.ChatWork.RoomID, message.MessageId, message.Account.Name, content)

	u.cw.PostRoomMessage(config.C.ChatWork.RoomID, messageBody)
}