package main

import "strings"

// sub topics
const (
	ThingDownlink = "fogcloud/+/+/thing/down/#"

	ThingPropertyPostReply = "fogcloud/+/+/thing/down/property/post_reply"
	ThingEventPostReply    = "fogcloud/+/+/thing/down/event/+/post_reply"
	ThingService           = "fogcloud/+/+/thing/down/service/+"

	ThingPropertySet = "fogcloud/+/+/thing/down/property/set"

	ThingTopoDownlink   = "fogcloud/+/+/thing/down/topo"
	ThingShadowDownlink = "fogcloud/+/+/thing/down/shadow"
	ThingNTPDownlink    = "fogcloud/+/+/thing/down/ntp"
)

// pub topics
const (
	ThingPropertyPost = "fogcloud/+/+/thing/up/property/post"
	ThingEventPost    = "fogcloud/+/+/thing/up/event/+/post"
	ThingServiceReply = "fogcloud/+/+/thing/up/service/+/reply"

	ThingTopoUplink   = "fogcloud/+/+/thing/up/topo"
	ThingShadowUplink = "fogcloud/+/+/thing/up/shadow"
	ThingNTPUplink    = "fogcloud/+/+/thing/up/ntp"
)

// sub topics index
const (
	IndexThingPropertyPostReply = iota
	IndexThingEventPostReply
	IndexThingService
	IndexThingPackDataPostReply
	IndexThingPropertySet
	IndexThingTopoDownlink
	IndexThingShadowDownlink
	IndexThingNTPDownlink
)

// pub topic index
const (
	IndexThingPropertyPost = iota
	IndexThingEventPost
	IndexThingServiceReply
	IndexThingPackDataPost
	IndexThingTopoUplink
	IndexThingShadowUplink
	IndexThingNTPUplink
)

var (
	SubTopics = []string{
		ThingDownlink,
	}

	PubTopics = []string{
		ThingPropertyPost,
		ThingEventPost,
		ThingServiceReply,
		ThingTopoUplink,
		ThingShadowUplink,
		ThingNTPUplink,
	}
)

func fillTopic(topic string, replaceStr ...string) string {
	s := topic
	for i := range replaceStr {
		s = strings.Replace(s, "+", replaceStr[i], 1)
	}
	return s
}

func fillTopics(rawTopics []string, productKey, deviceName string) []string {
	topics := make([]string, len(rawTopics))
	for i := range rawTopics {
		topics[i] = fillTopic(rawTopics[i], productKey, deviceName)
	}
	return topics
}
