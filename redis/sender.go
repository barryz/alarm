package redis

import (
	"alarm/g"
	"encoding/json"
	"log"
	"strings"

	"sender/model"
)

func LPUSH(queue, message string) {
	rc := g.RedisConnPool.Get()
	defer rc.Close()
	_, err := rc.Do("LPUSH", queue, message)
	if err != nil {
		log.Println("LPUSH redis", queue, "fail:", err, "message:", message)
	}
}

func WriteSmsModel(sms *model.Sms) {
	if sms == nil {
		return
	}

	bs, err := json.Marshal(sms)
	if err != nil {
		log.Println(err)
		return
	}

	LPUSH(g.Config().Queue.Sms, string(bs))
}

func WriteMailModel(mail *model.Mail) {
	if mail == nil {
		return
	}

	bs, err := json.Marshal(mail)
	if err != nil {
		log.Println(err)
		return
	}

	LPUSH(g.Config().Queue.Mail, string(bs))
}

func WriteSlackModel(slack *model.Slack) {
	if slack == nil {
		return
	}

	bs, err := json.Marshal(slack)
	if err != nil {
		log.Println(err)
		return
	}

	LPUSH(g.Config().Queue.Slack, string(bs))
}

func WriteSms(tos []string, content string) {
	if len(tos) == 0 {
		return
	}

	sms := &model.Sms{Tos: strings.Join(tos, ","), Content: content}
	WriteSmsModel(sms)
}

func WriteMail(tos []string, subject, content string) {
	if len(tos) == 0 {
		return
	}

	mail := &model.Mail{Tos: strings.Join(tos, ","), Subject: subject, Content: content}
	WriteMailModel(mail)
}

func WriteSlack(tos []string, content *model.SlackContent) {
	if len(tos) == 0 {
		return
	}

	slack := &model.Slack{Tos: strings.Join(tos, ","), Content: content}
	WriteSlackModel(slack)
}
