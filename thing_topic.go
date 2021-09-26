package main

import "strings"

// sub topics
const (
	ThingPropertyPostReply = "$fogcloud/+/+/thing/event/property/post_reply"
	ThingEventPostReply    = "$fogcloud/+/+/thing/event/+/post_reply"
	ThingService           = "$fogcloud/+/+/thing/service/+"
	ThingPackDataPostReply = "$fogcloud/+/+/thing/pack/data/post_reply"

	ThingPropertySet = "$fogcloud/+/+/thing/event/property/set"

	ThingTopoDownlink   = "$fogcloud/+/+/thing/topo/downlink"
	ThingShadowDownlink = "$fogcloud/+/+/thing/shadow/downlink"
	ThingNTPDownlink    = "$fogcloud/+/+/thing/ntp/downlink"
)

// pub topics
const (
	ThingPropertyPost = "$fogcloud/+/+/thing/event/property/post"
	ThingEventPost    = "$fogcloud/+/+/thing/event/+/post"
	ThingServiceReply = "$fogcloud/+/+/thing/service/+/reply"
	ThingPackDataPost = "$fogcloud/+/+/thing/pack/data/post"

	ThingTopoUplink   = "$fogcloud/+/+/thing/topo/uplink"
	ThingShadowUplink = "$fogcloud/+/+/thing/shadow/uplink"
	ThingNTPUplink    = "$fogcloud/+/+/thing/ntp/uplink"
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
		ThingPropertyPostReply,
		ThingEventPostReply,
		ThingService,
		ThingPackDataPostReply,
		ThingPropertySet,
		ThingTopoDownlink,
		ThingShadowDownlink,
		ThingNTPDownlink,
	}

	PubTopics = []string{
		ThingPropertyPost,
		ThingEventPost,
		ThingServiceReply,
		ThingPackDataPost,
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
