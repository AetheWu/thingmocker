package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"net"
	"sync/atomic"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Triad struct {
	ProductKey   string `csv:"product_key" json:"fog_pk"`
	DeviceName   string `csv:"device_name" json:"fog_dn"`
	DeviceSecret string `csv:"device_secret" json:"fog_ds"`
}

func NewDefalutThingMocker(productKey, deviceName, deviceSecret, ifaddr string) *ThingMocker {
	return &ThingMocker{
		productKey:   productKey,
		deviceName:   deviceName,
		deviceSecret: deviceSecret,
		signMethod:   "hmacsha1",
		clientId:     deviceName + "&" + productKey,

		subTopics: fillTopics(SubTopics, productKey, deviceName),
		pubTopics: fillTopics(PubTopics, productKey, deviceName),

		thingModel: getExampleThingModel(),

		ifaddr: ifaddr,
	}
}

type ThingMocker struct {
	client mqtt.Client

	ifaddr string

	deviceName   string
	productKey   string
	deviceSecret string
	thingModel   *Metadata

	clientId   string
	signMethod string

	subTopics []string
	pubTopics []string

	msgId uint32
}

func (t *ThingMocker) Conn() error {
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d/mqtt", Conf.MQTT_HOST, Conf.MQTT_PORT)).
		SetUsername(t.getUsername()).
		SetClientID(t.getClientId()).
		SetPassword(t.getPassword())

	if t.ifaddr != "" {
		dialer, err := t.newDialerWithIfaddr(t.ifaddr)
		if err != nil {
			return err
		}
		opts.SetDialer(dialer)
	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	t.client = c
	return nil
}

func (t *ThingMocker) newDialerWithIfaddr(ifaddr string) (*net.Dialer, error) {
	laddr, err := net.ResolveTCPAddr("", ifaddr)
	if err != nil {
		return nil, err
	}
	dial := &net.Dialer{
		Timeout:   time.Second * 10,
		LocalAddr: laddr,
	}
	return dial, nil
}

func (t *ThingMocker) DisConn() {
	t.client.Disconnect(0)
}

func (t *ThingMocker) String() string {
	return fmt.Sprintf("productKey[%s],deviceName[%s]", t.productKey, t.deviceName)
}

func (t *ThingMocker) getUsername() string {
	return t.deviceName + "&" + t.productKey
}

func (t *ThingMocker) getClientId() string {
	return fmt.Sprintf("%s|securemode=3,signmethod=%s|", t.clientId, t.signMethod)
}

func (t *ThingMocker) getPassword() string {
	sign, _ := authDeviceSign(t.deviceName, t.productKey, t.clientId, "", t.deviceSecret, t.signMethod)
	return sign
}

func (t *ThingMocker) PubMsg(topic string, qos byte, payload interface{}) error {
	if token := t.client.Publish(topic, qos, false, payload); token.Wait() {
		return token.Error()
	}
	return nil
}

func (t *ThingMocker) SubDefaultTopics() error {
	topics := make(map[string]byte, len(t.subTopics))
	for i := range t.subTopics {
		topics[t.subTopics[i]] = 0
	}

	tk := t.client.SubscribeMultiple(topics, func(client mqtt.Client, message mqtt.Message) {
		//Debugf("connected: %s", message.Payload())
	})
	if tk.Wait() && tk.Error() != nil {
		return tk.Error()
	}
	return nil
}

func (t *ThingMocker) PubProperties() error {
	rawData := generateExampleProperties(t.getId(), time.Now().Unix())
	return t.PubMsg(t.pubTopics[IndexThingPropertyPost], 0, rawData)
}

func (t *ThingMocker) PubEvents() error {
	rawData := generateExampleEvents(t.getId(), time.Now().Unix())
	return t.PubMsg(t.pubTopics[IndexThingEventPost], 0, rawData)
}

func (t *ThingMocker) getId() uint32 {
	return atomic.AddUint32(&t.msgId, 1)
}

func authDeviceSign(deviceName, productKey, clientId, timestamp, deviceSecret, signMethod string) (string, error) {
	src := ""
	src = fmt.Sprintf("clientId%sdeviceName%sproductKey%s", clientId, deviceName, productKey)
	if timestamp != "" {
		src = src + "timestamp" + timestamp
	}
	var h hash.Hash
	switch signMethod {
	case "hmacsha1":
		h = hmac.New(sha1.New, []byte(deviceSecret))
	case "hmacmd5":
		h = hmac.New(md5.New, []byte(deviceSecret))
	default:
		return "", errors.New("invalid sign method")
	}

	_, err := h.Write([]byte(src))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
