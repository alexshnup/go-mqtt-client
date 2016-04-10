package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"rpi.mqtt.client/conf"
	"rpi.mqtt.client/service"
	"rpi.mqtt.client/raspberry"
)

// check error
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// config
	config := conf.GetConfig()

	// interrupt
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// open mqtt connection
	client, err := service.NewMqttClient(
		config.Mqtt.Protocol,
		config.Mqtt.Address,
		config.Mqtt.Port,
		0,
	)
	CheckError(err)

	rpi := raspberry.NewRaspberry(client, config.Name, config.Debug)

	// Run publisher
	rpi.Cpu.Themperature.Publish(config.Timeout, 0)
	rpi.Cpu.CoreVolt.Publish(config.Timeout, 0)

	//go rpi.CpuTemp(config.Timeout, 0)
	//go rpi.CpuCoreVolt(config.Timeout, 0)
	go rpi.SystemMemory(config.Timeout, 0)

	// Run subscribe
	err = rpi.Led0(0)
	CheckError(err)

	// wait for terminating
	for {
		select {
		case <-interrupt:
			log.Println("Clean and terminating...")

			// disconnecting
			client.Disconnect(250)

			os.Exit(0)
		}
	}
}
