package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMocker(t *testing.T) {
	mustLoad("defaults", "./configs/config.yaml")
	filepath := "/Users/zhw/Documents/Triad.csv"
	t.Run("Start Mocker", func(t *testing.T) {
		StartMocker(filepath, 800, 10, 10)
	})

	t.Run("single thing mock", func(t *testing.T) {
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
}

func TestTransStrBytes(t *testing.T) {
	bytes := []byte{0x80, 0x81}
	str := string(bytes)
	nBytes := []byte(str)
	log.Printf("%s, %v", str, nBytes)

	rawStr := "�"
	nBytes = []byte(rawStr)
	log.Printf("%s, %v", rawStr, nBytes)

}