package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	successList []*ThingMocker
)

func readTriadFromFile(filepath string) (triads []Triad, err error) {
	rawData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}

	triads = []Triad{}
	err = json.Unmarshal(rawData, &triads)
	// err = gocsv.UnmarshalBytes(rawData, &triads)
	return
}

func StartMocker(ifaddr, filepath string, addStep, msgNum, duration, devNum int) {
	Println("start mocking")
	chDone := make(chan struct{}, 1)
	ctx, done := context.WithCancel(context.Background())

	go runMocker(ctx, chDone, ifaddr, filepath, addStep, msgNum, duration, devNum)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL)

	select {
	case <-sig:
	case <-chDone:
	}
	done()
	disconnectThingsByStep(successList, addStep)
	Println("end mocking gracefully")
}

func runMocker(ctx context.Context, chDone chan struct{}, ifaddr, filepath string, addStep, msgNum, duration, devNum int) {
	triads, err := readTriadFromFile(filepath)
	if err != nil {
		panic(err)
	}
	if len(triads) > devNum {
		triads = triads[:devNum]
	}
	things := initThingMockers(triads, ifaddr)
	successList = connThingsByStep(ctx, things, addStep)

	communicate(ctx, chDone, successList, msgNum, duration, addStep)
}

func initThingMockers(triads []Triad, ifaddr string) []*ThingMocker {
	things := make([]*ThingMocker, len(triads))
	for i := range triads {
		thing := NewDefalutThingMocker(triads[i].ProductKey, triads[i].DeviceName, triads[i].DeviceSecret, ifaddr)
		things[i] = thing
	}
	return things
}

func connThingsByStep(ctx context.Context, things []*ThingMocker, addStep int) []*ThingMocker {
	successList := make([]*ThingMocker, 0, len(things))
	tick := time.Tick(time.Second * 3)
	Println("start thing connecting")
loop:
	for left, right, i := 0, addStep, 0; left < len(things); left, right = i*addStep, (i+1)*addStep {
		if right > len(things) {
			right = len(things)
		}
		subs := connThingsConcurrency(things[left:right])
		successList = append(successList, subs...)
		i++
		select {
		case <-tick:
			continue
		case <-ctx.Done():
			break loop
		}
	}
	Println("end thing connecting")
	return successList
}

func connThingsConcurrency(things []*ThingMocker) (successThings []*ThingMocker) {
	thingCh := make(chan *ThingMocker, len(successThings))
	doneCh := make(chan struct{})

	connFn := func(wg *sync.WaitGroup, thing *ThingMocker) {
		defer wg.Done()
		err := thing.Conn()
		if err != nil {
			Printf("Conn: %s", err)
			return
		} else {
			thing.SubDefaultTopics()
			// if err != nil {
			// 	Printf("SubDefaultTopics: %s", err)
			// }
			thingCh <- thing
		}
	}

	rFn := func() {
		for i := range thingCh {
			successThings = append(successThings, i)
		}
		Printf("successNum: %d, failedNum: %d", len(successThings), len(things)-len(successThings))
		doneCh <- struct{}{}
	}
	go rFn()

	wg := new(sync.WaitGroup)
	for i := range things {
		wg.Add(1)
		go connFn(wg, things[i])
	}
	wg.Wait()
	close(thingCh)
	<-doneCh
	return
}

func disconnectThingsByStep(things []*ThingMocker, addStep int) {
	tick := time.Tick(time.Second * 3)
	Println("start thing disconnecting")
	for left, right, i := 0, addStep, 0; left < len(things); left, right = i*addStep, (i+1)*addStep {
		if right > len(things) {
			right = len(things)
		}
		disconnectThingsConcurrency(things[left:right])
		i++
		<-tick
		Printf("disconnected num: %d", addStep)
	}
	Println("end thing disconnecting")
}

func disconnectThingsConcurrency(things []*ThingMocker) {
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

func communicate(ctx context.Context, chDone chan struct{}, things []*ThingMocker, msgRate, duration, step int) {
	tick := time.NewTicker(time.Second)
	endTimer := time.After(time.Second * time.Duration(duration))
	Println("start thing communication mocking")
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-tick.C:
			go mockCommunicationsConcurrency(things, msgRate)
		case <-endTimer:
			break loop
		}
	}
	Println("end thing communication mocking")
	chDone <- struct{}{}
}

func mockCommunicationsConcurrency(things []*ThingMocker, msgRate int) {
	thingsNum := len(things)
	if msgRate > thingsNum {
		Debug("message trans rate should be less than num of thing-mockers")
		msgRate = thingsNum
		// return
	}

	startIndex := rand.Int63n(int64(thingsNum))
	commFn := func(index int) {
		if err := things[index].PubProperties(); err != nil {
			Debugf("thing[%s] PubProperties: %s", things[index], err)
		} else {
			Debugf("thing[%s] pub property success", things[index])
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
