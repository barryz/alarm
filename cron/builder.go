package cron

import (
	"alarm/g"
	cm "common/model"
	"common/utils"
	"fmt"
	sm "sender/model"
)

func BuildCommonSMSContent(event *cm.Event) string {
	return fmt.Sprintf(
		"{%s} {%s} {主机:%s} {%s} {当前值:%s, 判定条件:%s%s} {指标:%s} {Tags: %s } {告警次数:%d/%d} {时间:%s}",
		event.AlarmLevel(),
		event.StatusString(),
		event.Endpoint,
		event.Note(),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		event.CurrentStep,
		event.MaxStep(),
		event.FormattedTime(),
	)
}

func BuildCommonMailContent(event *cm.Event) string {
	link := g.Link(event)
	return fmt.Sprintf(
		"%s\r\n%s\r\n主机:%s\r\n指标:%s\r\nTags:%s\r\n判定函数:%s: 当前值:%s\r\n  判定条件: %s%s\r\n告警内容:%s\r\n最大报警次数:%d, 当前报警次数:%d\r\n时间:%s\r\n%s\r\n",
		event.StatusString(),
		event.AlarmLevel(),
		event.Endpoint,
		event.Metric(),
		utils.SortedTags(event.PushedTags),
		event.Func(),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.Note(),
		event.MaxStep(),
		event.CurrentStep,
		event.FormattedTime(),
		link,
	)
}

func BuildCommonSlackContent(event *cm.Event) *sm.SlackContent {
	return &sm.SlackContent{EndPoint: event.Endpoint,
		Note:                     event.Note(),
		Status:                   event.StatusString(),
		Priority:                 event.AlarmLevel(),
		Metric:                   event.Metric(),
		Tags:                     utils.SortedTags(event.PushedTags),
		CurrentValue:             utils.ReadableFloat(event.LeftValue),
		Expression:               fmt.Sprintf("%s%s", event.Operator(), utils.ReadableFloat(event.RightValue())),
		AlarmCount:               fmt.Sprintf("%d/%d", event.CurrentStep, event.MaxStep()),
		TriggerTime:              event.FormattedTime()}
}

func GenerateSmsContent(event *cm.Event) string {
	return BuildCommonSMSContent(event)
}

func GenerateMailContent(event *cm.Event) string {
	return BuildCommonMailContent(event)
}

func GenerateSlackContent(event *cm.Event) *sm.SlackContent {
	return BuildCommonSlackContent(event)
}
