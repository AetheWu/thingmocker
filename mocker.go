package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"
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

func StartMocker(filepath string, addStep, msgNum, duration int) {
	triads, err := readTriadFromFile(filepath)
	if err != nil {
		panic(err)
	}
	things := initThingMockers(triads)
	things = initThingsConnByStep(things, addStep)
	communicate(things, msgNum, duration)
}

func initThingMockers(triads [][3]string) []*ThingMocker {
	things := make([]*ThingMocker, len(triads))
	for i := range triads {
		thing := NewDefalutThingMocker(triads[i][2], triads[i][0], triads[i][1])
		things[i] = thing
	}
	return things
}

func initThingsConnByStep(things []*ThingMocker, addStep int) []*ThingMocker {
	successList := make([]*ThingMocker, 0, len(things))
	tick := time.Tick(time.Second * 3)
	for left, right, i := 0, addStep, 0; left < len(things); left, right = i*addStep, (i+1)*addStep {
		if right > len(things) {
			right = len(things)
		}
		subs := initThingsConnConcurrency(things[left:right])
		successList = append(successList, subs...)
		i++
		<-tick
	}
	return successList
	//fmt.Printf("success rate: %.4f\n", len(successList)/len(things))
	//ch := make(chan struct{})
	//<- ch
}

func initThingsConnConcurrency(things []*ThingMocker) (successThings []*ThingMocker) {
	thingCh := make(chan *ThingMocker, len(successThings))
	failedCh := make(chan struct{}, len(successThings))

	connFn := func(wg *sync.WaitGroup, thing *ThingMocker) {
		defer wg.Done()
		err := thing.Conn()
		if err != nil {
			failedCh <- struct{}{}
		} else {
			err = thing.SubDefaultTopics()
			if err != nil {
				Debugf("SubDefaultTopics: %s", err)
			}
			thingCh <- thing
		}
	}

	countFn := func() {
		failedNum := 0
		for _ = range failedCh {
			failedNum++
		}
		Infof("failed num: %d\n", failedNum)
	}

	rFn := func() {
		for i := range thingCh {
			successThings = append(successThings, i)
		}
		Infof("success num: %d\n", len(successThings))
	}
	go countFn()
	go rFn()

	wg := new(sync.WaitGroup)
	for i := range things {
		wg.Add(1)
		go connFn(wg, things[i])
	}
	wg.Wait()
	close(failedCh)
	close(thingCh)
	return
}

func communicate(things []*ThingMocker, msgRate, duration int) {
	tick := time.NewTicker(time.Second)
	endTimer := time.After(time.Second * time.Duration(duration))
	Info("start thing communication mocking")
loop:
	for {
		select {
		case <-tick.C:
			mockCommunicationsConcurrency(things, msgRate)
		case <-endTimer:
			break loop
		}
	}
	Info("end thing communication mocking")
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
