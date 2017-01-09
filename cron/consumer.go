package cron

import (
	"encoding/json"
	"log"

	"alarm/api"

	"alarm/g"

	"alarm/redis"

	"common/model"
)

func consume(event *model.Event, isHigh bool) {
	// 获取当前事件的表达式id或者策略id
	actionId := event.ActionId()
	if actionId <= 0 {
		return
	}

	// 根据actionid请求portal api 获取action
	action := api.GetAction(actionId)
	if action == nil {
		return
	}

	// 如果有回调方法， 先处理回调
	if action.Callback == 1 {
		HandleCallback(event, action)
		return
	}

	// 告警事件分优先级处理
	if isHigh {
		consumeHighEvents(event, action)
	} else {
		consumeLowEvents(event, action)
	}
}

// 高优先级的不做报警合并
func consumeHighEvents(event *model.Event, action *api.Action) {
	// action.Uic 在这里是用户组
	if action.Uic == "" {
		return
	}

	var (
		phones, mails, slacks []string
	)
	// 通过action查询是否发送到channel, 否则发送到user
	if len(action.SlackChannel) != 0 {
		slacks = []string{action.SlackChannel}
	} else {
		//  这里是通过portal 的api 查询到用户组对应的联系人信息
		phones, mails, slacks = api.ParseTeams(action.Uic)
	}

	smsContent := GenerateSmsContent(event)
	mailContent := GenerateMailContent(event)
	slackContent := GenerateSlackContent(event)

	// 这里定义了 如果优先级 < 3 时需要发短信
	if event.Priority() < 3 {
		redis.WriteSms(phones, smsContent)
	}

	redis.WriteMail(mails, smsContent, mailContent)
	redis.WriteSlack(slacks, slackContent)
}

// 低优先级的做报警合并
func consumeLowEvents(event *model.Event, action *api.Action) {
	if action.Uic == "" {
		return
	}

	var slacks []string

	if len(action.SlackChannel) != 0 {
		slacks := []string{action.SlackChannel}
	} else {
		_, _, slacks := api.ParseTeams(action.Uic)
	}

	slackContent := GenerateSlackContent(event)

	// 低优先级且优先级 < 3的情况下, 解析用户的告警短信，并写入redis， 后续做短信合并
	if event.Priority() < 3 {
		ParseUserSms(event, action)
	}

	//  解析用户的告警邮件， 并写入redis， 后续做邮件合并
	ParseUserMail(event, action)

	// 低优先级的告警事件也需要发到slack里
	redis.WriteSlack(slacks, slackContent)
}

func ParseUserSms(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)

	content := GenerateSmsContent(event)
	metric := event.Metric()
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserSmsQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := SmsDto{
			Priority: priority,
			Metric:   metric,
			Content:  content,
			Phone:    user.Phone,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Println("json marshal SmsDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Println("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}

func ParseUserMail(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)

	metric := event.Metric()
	subject := GenerateSmsContent(event)
	content := GenerateMailContent(event)
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserMailQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := MailDto{
			Priority: priority,
			Metric:   metric,
			Subject:  subject,
			Content:  content,
			Email:    user.Email,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Println("json marshal MailDto fail:", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Println("LPUSH redis", queue, "fail:", err, "dto:", string(bs))
		}
	}
}
