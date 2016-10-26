package cron

import (
	"fmt"
	"github.com/barryz/alarm/g"
	"github.com/barryz/common/model"
	"github.com/barryz/common/utils"
)

func BuildCommonSMSContent(event *model.Event) string {
	return fmt.Sprintf(
		"{%s} {%s} {主机:%s} {%s} {当前值:%s, 判定条件:%s%s} {指标:%s} {告警次数:%d/%d} {时间:%s}",
		event.AlarmLevel(),
		event.StatusString(),
		event.Endpoint,
		event.Note(),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.Metric(),
		event.CurrentStep,
		event.MaxStep(),
		event.FormattedTime(),
	)
}

func BuildCommonMailContent(event *model.Event) string {
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

func BuildCommonSlackContent(event *model.Event) (string, string, string) {
	return fmt.Sprintf(
		"主机: %s %s (%s)",
		event.Endpoint,
		event.Note(),
		event.StatusString(),
	), fmt.Sprintf(
		"%s",
		event.Status,
	), fmt.Sprintf(
		"{%s %s}{%s}{指标:%s}{当前值:%s,判定条件%s%s}{告警次数%d/%d}{异常开始时间:%s}",
		event.AlarmLevel(),
		event.Endpoint,
		event.Note(),
		event.Metric(),
		utils.ReadableFloat(event.LeftValue),
		event.Operator(),
		utils.ReadableFloat(event.RightValue()),
		event.CurrentStep,
		event.MaxStep(),
		event.FormattedTime(),
	)
}

func GenerateSmsContent(event *model.Event) string {
	return BuildCommonSMSContent(event)
}

func GenerateMailContent(event *model.Event) string {
	return BuildCommonMailContent(event)
}

func GenerateSlackContent(event *model.Event) (string, string, string) {
	return BuildCommonSlackContent(event)
}
