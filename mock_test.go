package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMocker(t *testing.T) {
	filepath := "/Users/zhw/Documents/Triad.csv"

	t.Run("Start Mocker", func(t *testing.T) {
		StartMocker(filepath, 800, 10, 10)
	})

	t.Run("single thing mock", func(t *testing.T) {
		mustLoad("defaults", "./configs/config.yaml")
		dn, pk, ds := "00000000", "a1bb9d14", "4450c02f5fe642e206c39980a6629ad8"
		thing := NewDefalutThingMocker(pk, dn, ds)
		thing.Conn()
		ch := make(chan struct{})
		<- ch
	})

	t.Run("thing mocker", func(t *testing.T) {
		triads, err := readTriadFromFile(filepath)
		assert.NoError(t, err)
		thing := NewDefalutThingMocker(triads[0][2], triads[0][0], triads[0][1])
		err = thing.Conn()
		if err == nil {
			err = thing.SubDefaultTopics()
		}
		assert.NoError(t, err)
		err = thing.PubProperties()
		assert.NoError(t, err)
		ch:=make(chan struct{})
		<-ch
	})

	t.Run("things mocker", func(t *testing.T) {
		triads, err := readTriadFromFile(filepath)
		assert.NoError(t, err)
		things := initThingMockers(triads)
		initThingsConnConcurrency(things[:10])
		//time.Sleep(time.Second*3)
		ch:=make(chan struct{})
	    <-ch
	})
}