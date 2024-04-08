package thingmocker

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func Run() error {
	log.Println("Mocker run")
	sched, err := NewMockerScheduler(Conf)
	if err != nil {
		return err
	}
	go sched.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-c
	sched.Stop()
	log.Println("Mocker stop")
	return nil
}

func readTriadFromFile(filepath string) (triads []Triad, err error) {
	rawData, err := os.ReadFile(filepath)
	if err != nil {
		return
	}

	triads = []Triad{}
	err = json.Unmarshal(rawData, &triads)
	// err = gocsv.UnmarshalBytes(rawData, &triads)
	return
}

func readMsgFromFile(filepath string) (map[string]map[string][]MockerMsg, error) {
	rawData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	msgs := make(map[string]map[string][]MockerMsg)
	err = json.Unmarshal(rawData, &msgs)
	return msgs, err
}

type MockerScheduler struct {
	mockers, connectedMockers []*ThingMocker
	ch                        chan *ThingMocker
	closeCh                   chan struct{}
	cfg                       ConfigData
}

func NewMockerScheduler(cfg ConfigData) (*MockerScheduler, error) {
	s := &MockerScheduler{
		cfg:     cfg,
		closeCh: make(chan struct{}),
		ch:      make(chan *ThingMocker, 100),
	}
	triads, err := readTriadFromFile(cfg.DEVICE_TRIAD_FILEPATH)
	if err != nil {
		return nil, err
	}
	msgs, err := readMsgFromFile(cfg.COMM_FILEPATH)
	if err != nil {
		return nil, err
	}
	s.mockers = newThingMockers(triads, cfg.IF_ADDR, msgs)
	return s, nil
}

func (s *MockerScheduler) Run() {
	s.mockConnect()
	s.mockPublish()
}

func (s *MockerScheduler) Stop() {
	close(s.closeCh)
	s.mockDisconnect()
}

func newThingMockers(triads []Triad, ifaddr string, msgs map[string]map[string][]MockerMsg) []*ThingMocker {
	things := make([]*ThingMocker, len(triads))
	for i := range triads {
		pkMockerMsgs := msgs[triads[i].ProductKey]
		var thingMockerMsgs []MockerMsg
		if pkMockerMsgs != nil {
			thingMockerMsgs = pkMockerMsgs[triads[i].DeviceName]
			msgs, ok := pkMockerMsgs["*"]
			if ok {
				thingMockerMsgs = append(thingMockerMsgs, msgs...)
			}
		}
		thing := NewDefalutThingMocker(triads[i].ProductKey, triads[i].DeviceName, triads[i].DeviceSecret, ifaddr, thingMockerMsgs)
		things[i] = thing
	}
	return things
}

func (s *MockerScheduler) mockConnect() {
	tick := time.NewTicker(time.Second * 3)
	defer tick.Stop()
	log.Println("start thing connecting")
	thingNum := len(s.mockers)
	if s.cfg.DEVICE_NUM > thingNum {
		s.cfg.DEVICE_NUM = thingNum
	}
	go func() {
	loop:
		for left, right, i := 0, s.cfg.DEVICE_STEP_NUM, 0; left < s.cfg.DEVICE_NUM; left, right = i*s.cfg.DEVICE_STEP_NUM, (i+1)*s.cfg.DEVICE_STEP_NUM {
			if right > s.cfg.DEVICE_NUM {
				right = s.cfg.DEVICE_NUM
			}
			s.connThingsConcurrency(s.mockers[left:right])
			i++
			select {
			case <-tick.C:
				continue
			case <-s.ch:
				break loop
			}
		}
		close(s.ch)
	}()
	s.recvConnectedThingMocker()
	log.Println("end thing connecting")
}

func (s *MockerScheduler) recvConnectedThingMocker() {
	for i := range s.ch {
		s.connectedMockers = append(s.connectedMockers, i)
		go s.mockThingPublish(i)
	}
	log.Printf("Things all connected num: %d", len(s.connectedMockers))
}

func (s *MockerScheduler) mockThingPublish(thing *ThingMocker) {
	tick := time.NewTicker(time.Second)
	for i := 0; i < 5; i++ {
		for i := range thing.mockerMsgs {
			topic := thing.mockerMsgs[i].GetTopic(thing.productKey, thing.deviceName)
			payload := thing.mockerMsgs[i].GetPayload()
			thing.PubMsg(topic, 0, payload)
		}
		<-tick.C
	}
}

func (s *MockerScheduler) connThingsConcurrency(things []*ThingMocker) {
	connFn := func(wg *sync.WaitGroup, thing *ThingMocker) {
		defer wg.Done()
		err := thing.Conn()
		if err != nil {
			Printf("Conn: %s", err)
			return
		} else {
			thing.SubDefaultTopics()
			s.ch <- thing
		}
	}
	wg := new(sync.WaitGroup)
	for i := range things {
		wg.Add(1)
		go connFn(wg, things[i])
	}
	wg.Wait()
	log.Printf("connected num: %d", len(things))
}

func (s *MockerScheduler) mockDisconnect() {
	tick := time.NewTicker(time.Second * 3)
	defer tick.Stop()
	log.Println("start thing disconnecting")
	for left, right, i := 0, s.cfg.DEVICE_STEP_NUM, 0; left < len(s.connectedMockers); left, right = i*s.cfg.DEVICE_STEP_NUM, (i+1)*s.cfg.DEVICE_STEP_NUM {
		if right > len(s.connectedMockers) {
			right = len(s.connectedMockers)
		}
		s.mockDisconnectOnConcurrency(s.connectedMockers[left:right])
		i++
		<-tick.C
		Printf("disconnected num: %d", s.cfg.DEVICE_STEP_NUM)
	}
	log.Println("end thing disconnecting")
}

func (s *MockerScheduler) mockDisconnectOnConcurrency(things []*ThingMocker) {
	wg := new(sync.WaitGroup)
	for i := range things {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			things[i].DisConn()
		}(i)
	}

	wg.Wait()
}

func (s *MockerScheduler) mockPublish() {
	tick := time.NewTicker(time.Second)
	Println("start thing communication mocking")
loop:
	for {
		select {
		case <-tick.C:
			go s.mockPublishOnConcurrency()
		case <-s.closeCh:
			break loop
		}
	}
	Println("end thing communication mocking")
}

func (s *MockerScheduler) mockPublishOnConcurrency() {
	thingsNum := len(s.connectedMockers)
	msgRate := s.cfg.MESSAGE_RATE
	if msgRate > thingsNum {
		msgRate = thingsNum
	}
	startIndex := rand.Int63n(int64(thingsNum))
	commFn := func(index int) {
		for i := range s.connectedMockers[index].mockerMsgs {
			if !strings.Contains(s.connectedMockers[index].mockerMsgs[i].Payload, "login") {
				topic := s.connectedMockers[index].mockerMsgs[i].GetTopic(s.connectedMockers[index].productKey, s.connectedMockers[index].deviceName)
				payload := s.connectedMockers[index].mockerMsgs[i].GetPayload()
				s.connectedMockers[index].PubMsg(topic, 0, payload)
			}
		}
	}

	endIndex := int(startIndex) + msgRate
	if endIndex > thingsNum {
		for i := int(startIndex); i < thingsNum; i++ {
			commFn(i)
		}
		for i := 0; i < endIndex-thingsNum; i++ {
			commFn(i)
		}
	} else {
		for i := int(startIndex); i < endIndex; i++ {
			commFn(i)
		}
	}
}
