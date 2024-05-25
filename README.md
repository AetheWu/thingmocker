## 使用教程

### 1. 配置说明
thingmocker需要读取本地配置文件`config.yaml`，配置文件格式参考`configs/config.yaml`；

| 配置项 | 默认值 | 说明 | 
| --- | --- |  --- |
| MQTT_HOST | "localhost" | MQTT服务地址 |
| MQTT_PORT | 1883 | MQTT服务端口 |
| MQTT_TLS | false |  是否使用mqtt tls |
| IF_ADDR | "" | 客户端连接时使用的网卡地址，不填则为默认 |
| MESSAGE_RATE | 1000 | 模拟设备上报数据频率，单位为：messages/秒 |
| MESSAGE_DURATION | 3600*24*7 | 模拟器运行时间，单位为：秒 |
| DEVICE_STEP_NUM | 100 | 设备每秒连接数量 | 
| DEVICE_NUM | 100 | 模拟设备总数，注意：需要小于或等于三元组文件提供的设备数量 |
| DEVICE_TRIAD_FILEPATH | "/etc/thingmocker/triads.txt" | 模拟设备三元组信息文件路径；文件格式为FogCloud后台导出的三元组文件： 设备->设备列表->批次管理->下载TXT|
| COMM_FILEPATH | "/etc/thingmocker/comm.csv"  |  模拟设备上报数据的文件路径，用来自定义设置设备上报数据 |

> 上报数据文件格式为JSON，文件示例：
```json
{
    "bf3de5ef": { // 产品product_key
        "*": [    // 设备device_name，当设置为"*"时，该产品下所有设备上报数据均设置为该配置
            {
                "topic": "fogcloud/+/+/thing/up/through", //上报MQTT主题
                "payload": "0100320001000001000431120000", //上报MQTT载荷
                "format": "hex" //上报MQTT载荷格式，可选：hex、plaintext
            }
        ]
    }
}
```

### 2. 编译
需要提前配置[Go1.16+](https://golang.org/dl/)开发环境
```bash
make build
```

### 3. 运行 
```bash
make run
```
或者直接运行编译得到的可执行文件：
```bash
thingmocker mock -c config.yaml -e defaults
```






