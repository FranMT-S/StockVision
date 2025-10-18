package models

import "strings"

type Action string

const (
	downgradedBy    Action = "downgraded"
	initiatedBy     Action = "initiated"
	reiteratedBy    Action = "reiterated"
	targetLoweredBy Action = "target lowered"
	targetRaisedBy  Action = "target raised"
	targetSetBy     Action = "target set"
	upgradedBy      Action = "upgraded"
)

var actionsMap = map[string]Action{
	"downgraded":     downgradedBy,
	"initiated":      initiatedBy,
	"reiterated":     reiteratedBy,
	"target lowered": targetLoweredBy,
	"target raised":  targetRaisedBy,
	"target set":     targetSetBy,
	"upgraded":       upgradedBy,
}

func (a Action) IsValid() bool {
	return actionsMap[strings.ToLower(string(a))] != ""
}

func (a Action) Normalize() Action {
	var action string = strings.ToLower(string(a))
	if strings.Contains(action, "by") {
		action = strings.ReplaceAll(action, "by", "")
	}

	action = strings.TrimSpace(action)

	return actionsMap[action]
}
