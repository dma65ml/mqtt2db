package main

import (
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/robertkrimen/otto"
)

var wg sync.WaitGroup // syncro goroutine

var vm *otto.Otto

func main() {

	// configuration
	cfgPath, err := ParseFlags()
	if err != nil {
		log.Error(err)
	}
	config, err := NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	vm = otto.New()

	client := initMQTT(config)
	defer disconnectMQTT(client)

	wg.Add(1)
	go db()

	wg.Wait()

}
