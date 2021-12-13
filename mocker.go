package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
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

func readTriadFromFile(filepath string) (triad [][3]string, err error) {
	rawData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}

	lines, err := csv.NewReader(bytes.NewReader(rawData)).ReadAll()
	if err != nil {
		return
	}

	if len(lines) == 0 {
		err = errors.New("empty triad")
		return
	}

	triad = make([][3]string, len(lines)-1)
	for i := 1; i < len(lines); i++ {
		if len(lines[i]) != 3 {
			err = errors.New("invalid csv format")
			return
		}
		for j := range triad[i-1] {
			triad[i-1][j] = lines[i][j]
		}
	}
	return
}

func StartMocker(filepath string, addStep, msgNum, duration, devNum int) {
	Println("start mocking")
	chDone := make(chan struct{}, 1)
	ctx, done := context.WithCancel(context.Background())

	go runMocker(ctx, chDone, filepath, addStep, msgNum, duration, devNum)

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

func runMocker(ctx context.Context, chDone chan struct{}, filepath string, addStep, msgNum, duration, devNum int) {
	triads, err := readTriadFromFile(filepath)
	if err != nil {
		panic(err)
	}
	if len(triads) > devNum {
		triads = triads[:devNum]
	}
	things := initThingMockers(triads)
	successList = connThingsByStep(ctx, things, addStep)

	communicate(ctx, chDone, successList, msgNum, duration, addStep)
}

func initThingMockers(triads [][3]string) []*ThingMocker {
	things := make([]*ThingMocker, len(triads))
	for i := range triads {
		thing := NewDefalutThingMocker(triads[i][2], triads[i][0], triads[i][1])
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
			err = thing.SubDefaultTopics()
			if err != nil {
				Printf("SubDefaultTopics: %s", err)
			}
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
		return
	}

	startIndex := rand.Int63n(int64(thingsNum))
	commFn := func(wg *sync.WaitGroup, index int) {
		defer wg.Done()
		if err := things[index].PubProperties(); err != nil {
			Debugf("thing[%s] PubProperties: %s", things[index], err)
		}
	}

	wg := new(sync.WaitGroup)
	endIndex := int(startIndex) + msgRate
	if endIndex > thingsNum {
		for i := int(startIndex); i < thingsNum; i++ {
			wg.Add(1)
			go commFn(wg, i)
		}
		for i := 0; i < endIndex-thingsNum; i++ {
			wg.Add(1)
			go commFn(wg, i)
		}
	} else {
		for i := int(startIndex); i < endIndex; i++ {
			wg.Add(1)
			go commFn(wg, i)
		}
	}

	wg.Wait()
}
