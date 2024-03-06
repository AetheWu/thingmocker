package thingmocker

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

var (
	onceMetadata sync.Once
	metadata     Metadata

	rawThingModelExample = `
{"events": {"post": {"id": "post", "desc": "属性上报", "name": "post", "method": "thing.event.property.post", "required": true, "standard": false, "eventType": "", "outputData": {"LightLux": {"id": "LightLux", "name": "光照值", "dataType": {"max": "65535", "min": "0", "step": "1", "type": "int", "unit": "LUX", "unitName": "照度"}}, "date_test": {"id": "date_test", "name": "日期", "dataType": {"type": "date"}}, "enum_test": {"id": "enum_test", "name": "meiju", "dataType": {"type": "enum", "elements": {"0": "红", "1": "橙", "2": "黄"}}}, "array_test": {"id": "array_test", "name": "数组", "dataType": {"size": 10, "type": "array", "elementType": "int"}}, "float_test": {"id": "float_test", "name": "浮点", "dataType": {"max": "99.9", "min": "0", "step": "0.5", "type": "float", "unit": "W/㎡", "unitName": "太阳总辐射"}}, "struct_test": {"id": "struct_test", "name": "结构体", "dataType": {"type": "struct", "structFields": {"field1": {"id": "field1", "name": "字段1", "dataType": {"max": "99", "min": "1", "step": "2", "type": "int"}}, "field2": {"id": "field2", "name": "字段2", "dataType": {"max": "99", "min": "1", "step": "0.8", "type": "float"}}}}}}}}, "services": {"get": {"id": "get", "desc": "属性获取", "name": "get", "method": "thing.service.property.get", "callType": "async", "required": true, "standard": false, "eventType": "", "inputData": {"LightLux": {"id": "LightLux", "name": "光照值", "dataType": {"max": "65535", "min": "0", "step": "1", "type": "int", "unit": "LUX", "unitName": "照度"}}, "date_test": {"id": "date_test", "name": "日期", "dataType": {"type": "date"}}, "enum_test": {"id": "enum_test", "name": "meiju", "dataType": {"type": "enum", "elements": {"0": "红", "1": "橙", "2": "黄"}}}, "array_test": {"id": "array_test", "name": "数组", "dataType": {"size": 10, "type": "array", "elementType": "int"}}, "float_test": {"id": "float_test", "name": "浮点", "dataType": {"max": "99.9", "min": "0", "step": "0.5", "type": "float", "unit": "W/㎡", "unitName": "太阳总辐射"}}, "struct_test": {"id": "struct_test", "name": "结构体", "dataType": {"type": "struct", "structFields": {"field1": {"id": "field1", "name": "字段1", "dataType": {"max": "99", "min": "1", "step": "2", "type": "int"}}, "field2": {"id": "field2", "name": "字段2", "dataType": {"max": "99", "min": "1", "step": "0.8", "type": "float"}}}}}}, "outputData": {"LightLux": {"id": "LightLux", "name": "光照值", "dataType": {"max": "65535", "min": "0", "step": "1", "type": "int", "unit": "LUX", "unitName": "照度"}}, "date_test": {"id": "date_test", "name": "日期", "dataType": {"type": "date"}}, "enum_test": {"id": "enum_test", "name": "meiju", "dataType": {"type": "enum", "elements": {"0": "红", "1": "橙", "2": "黄"}}}, "array_test": {"id": "array_test", "name": "数组", "dataType": {"size": 10, "type": "array", "elementType": "int"}}, "float_test": {"id": "float_test", "name": "浮点", "dataType": {"max": "99.9", "min": "0", "step": "0.5", "type": "float", "unit": "W/㎡", "unitName": "太阳总辐射"}}, "struct_test": {"id": "struct_test", "name": "结构体", "dataType": {"type": "struct", "structFields": {"field1": {"id": "field1", "name": "字段1", "dataType": {"max": "99", "min": "1", "step": "2", "type": "int"}}, "field2": {"id": "field2", "name": "字段2", "dataType": {"max": "99", "min": "1", "step": "0.8", "type": "float"}}}}}}}, "set": {"id": "set", "desc": "属性设置", "name": "set", "method": "thing.service.property.set", "callType": "async", "required": true, "standard": false, "eventType": "", "inputData": {"date_test": {"id": "date_test", "name": "日期", "dataType": {"type": "date"}}, "enum_test": {"id": "enum_test", "name": "meiju", "dataType": {"type": "enum", "elements": {"0": "红", "1": "橙", "2": "黄"}}}, "array_test": {"id": "array_test", "name": "数组", "dataType": {"size": 10, "type": "array", "elementType": "int"}}, "float_test": {"id": "float_test", "name": "浮点", "dataType": {"max": "99.9", "min": "0", "step": "0.5", "type": "float", "unit": "W/㎡", "unitName": "太阳总辐射"}}, "struct_test": {"id": "struct_test", "name": "结构体", "dataType": {"type": "struct", "structFields": {"field1": {"id": "field1", "name": "字段1", "dataType": {"max": "99", "min": "1", "step": "2", "type": "int"}}, "field2": {"id": "field2", "name": "字段2", "dataType": {"max": "99", "min": "1", "step": "0.8", "type": "float"}}}}}}, "outputData": {}}, "TimeReset": {"id": "TimeReset", "desc": "", "name": "设备校时服务", "method": "thing.service.TimeReset", "callType": "async", "required": false, "standard": false, "eventType": "", "inputData": {"TimeReset": {"id": "TimeReset", "name": "TimeReset", "dataType": {"type": "text", "length": 255}}}, "outputData": {"output1": {"id": "output1", "name": "输出2", "dataType": {"max": "999999999", "min": "1", "step": "1", "type": "int"}}}}}, "properties": {"LightLux": {"id": "LightLux", "name": "光照值", "dataType": {"max": "65535", "min": "0", "step": "1", "type": "int", "unit": "LUX", "unitName": "照度"}, "required": false, "standard": false, "accessMode": "r"}, "date_test": {"id": "date_test", "name": "日期", "dataType": {"type": "date"}, "required": false, "standard": false, "accessMode": "rw"}, "enum_test": {"id": "enum_test", "name": "meiju", "dataType": {"type": "enum", "elements": {"0": "红", "1": "橙", "2": "黄"}}, "required": false, "standard": false, "accessMode": "rw"}, "array_test": {"id": "array_test", "name": "数组", "dataType": {"size": 10, "type": "array", "elementType": "int"}, "required": false, "standard": false, "accessMode": "rw"}, "float_test": {"id": "float_test", "name": "浮点", "dataType": {"max": "99.9", "min": "0", "step": "0.5", "type": "float", "unit": "W/㎡", "unitName": "太阳总辐射"}, "required": false, "standard": false, "accessMode": "rw"}, "struct_test": {"id": "struct_test", "name": "结构体", "dataType": {"type": "struct", "structFields": {"field1": {"id": "field1", "name": "字段1", "dataType": {"max": "99", "min": "1", "step": "2", "type": "int"}}, "field2": {"id": "field2", "name": "字段2", "dataType": {"max": "99", "min": "1", "step": "0.8", "type": "float"}}}}, "required": false, "standard": false, "accessMode": "rw"}}}
`
)

func getExampleThingModel() *Metadata {
	onceMetadata.Do(func() {
		err := json.Unmarshal([]byte(rawThingModelExample), &metadata)
		if err != nil {
			panic(fmt.Sprintf("getExampleThingModel: %s", err))
		}
	})
	return &metadata
}

func generateExampleProperties(id uint32, timestamp int64) []byte {
	msg := ThingJsonPropPost{
		ThingJsonHeader: ThingJsonHeader{
			Id:        id,
			Version:   "1.0",
			Timestamp: timestamp,
		},
		Params: map[string]interface{}{
			"ACOnlineNum":             id,
			"ACTotalNum":              id,
			"BatteryClusterOnlineNum": timestamp,
			"BatteryClusterTotalNum":  id,
			"BatteryStackOnlineNum":   id,
			"BatteryStackTotalNum":    id,
			"CabinetOnlineNum":        id,
			"CabinetTotalNum":         id,
			"ChargesNum":              id,
		},
	}

	rawData, _ := json.Marshal(msg)
	return rawData
}

func generateExampleEvents(id uint32, timestamp int64) []byte {
	msg := ThingJsonEventPost{
		ThingJsonHeader: ThingJsonHeader{
			Id:      id,
			Version: "1.0",
		},
		Params: map[string]ThingJsonEventParam{
			"post": {
				Value: map[string]interface{}{
					"array_test": []int{1, 2, 3},
					"date_test":  "2021-08-10",
				},
				Time: timestamp,
			},
		},
	}
	rawData, _ := json.Marshal(msg)
	return rawData
}

func NewMockerMsg(topic, payload, format string) MockerMsg {
	return MockerMsg{
		Topic:   topic,
		Payload: payload,
		Format:  format,
	}
}

type MockerMsg struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
	Format  string `json:"format"` // json,hex
}

func (m MockerMsg) GetPayload() []byte {
	if m.Format == "hex" {
		bytes, _ := hex.DecodeString(m.Payload)
		return bytes
	}

	return []byte(m.Payload)
}

func (m MockerMsg) GetTopic(pk, dn string) string {
	s := strings.Replace(m.Topic, "+", pk, 1)
	return strings.Replace(s, "+", dn, 1)
}
