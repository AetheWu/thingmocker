package main

type Metadata struct {
	Properties map[string]*MetaProperty `json:"properties"`         //属性
	Events     map[string]*MetaEvent    `json:"events"`             //事件
	Services   map[string]*MetaService  `json:"services,omitempty"` //服务
}

type MetaBaseElement struct {
	Id       string `json:"id"`       //产品下唯一属性标识
	Name     string `json:"name"`     //属性名称
	Required bool   `json:"required"` //是否标准功能必选属性
	Standard bool   `json:"standard"` //是否是标准属性
	FuncType string `json:"funcType,omitempty"`
}

type MetaProperty struct {
	MetaBaseElement
	AccessMode string   `json:"accessMode"` //属性权限: r(只读), w(只写), rw(读写)
	DType      DataType `json:"dataType"`
}

type MetaEvent struct {
	MetaBaseElement
	EventType  string                  `json:"eventType"` //事件类型
	Desc       string                  `json:"desc"`      //服务描述
	Method     string                  `json:"method,omitempty"`
	OutputData map[string]*StructField `json:"outputData"`
}

type MetaService struct {
	MetaEvent
	CallType  string                  `json:"callType"`  //async(异步), sync(同步)
	InputData map[string]*StructField `json:"inputData"` //输入参数
}

type DataType struct {
	Spec         `json:",inline"`
	Type         string                  `json:"type"`                   //value(数值), string(字符串), date(日期), bool(布尔), enum(枚举), array(数组), struct(自定义结构体)
	StructFields map[string]*StructField `json:"structFields,omitempty"` //结构体类型字段定义
}

type Spec struct {
	Unit        string         `json:"unit,omitempty"`        //单位
	UnitName    string         `json:"unitName,omitempty"`    //单位名称
	Min         string         `json:"min,omitempty"`         //参数最小值(数值类型特有)
	Max         string         `json:"max,omitempty"`         //参数最大值(数值类型特有)
	Step        string         `json:"step,omitempty"`        //步长
	Size        int            `json:"size,omitempty"`        //数组类型数组长度
	Length      int            `json:"length,omitempty"`      //字符串类型长度
	ElementType string         `json:"elementType,omitempty"` //数组元素类型
	Elements    map[int]string `json:"elements,omitempty"`    //枚举和布尔类型的元素定义
}

type StructField struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	DataType DataType `json:"dataType"`
}

type ThingJsonHeader struct {
	Id      string `json:"id"`
	Version string `json:"version"`
	Sys     struct {
		Ack int `json:"ack"`
	} `json:"sys"`
	Method string `json:"method"`
	Timestamp int64  `json:"timestamp"`
}

type ThingJsonPropParam struct {
	Value interface{} `json:"value"`
	Time  int64       `json:"time"`
}

type ThingJsonEventParam struct {
	Value map[string]interface{} `json:"value"`
	Time  int64                  `json:"time"`
}

type ThingJsonPropPost struct {
	ThingJsonHeader
	Params map[string]interface{}
}

type ThingJsonEventPost struct {
	ThingJsonHeader
	Params map[string]ThingJsonEventParam
}
