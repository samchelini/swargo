package main

import (
	"github.com/samchelini/swargo/bar"
	"github.com/samchelini/swargo/netlink"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stderr) // log to stderr

	// initialize bar
	b := bar.NewBar()
	//b.EnablePrettyPrint()
	b.EnableLogging()

	// create and add a DateTimeBlock
	dateTimeBlock := new(bar.DateTimeBlock)
	b.AddBlock(dateTimeBlock)

	// create and add a BrightnessBlock
	brightnessBlock := new(bar.BrightnessBlock)
	brightnessBlock.SetDir("/sys/class/backlight/intel_backlight")
	brightnessBlock.SetPrefix("\u2600")
	b.AddBlock(brightnessBlock)

	/*
		// create and add a BatteryBlock
		batteryBlock := new(bar.BatteryBlock)
		batteryBlock.SetDir("/sys/class/power_supply/BAT0")
		batteryBlock.SetChargingPrefix("CHG")
		batteryBlock.SetDischargingPrefix("BAT")
		b.AddBlock(batteryBlock)
	*/

  // netlink testing
	log.Println("creating generic netlink connection...")
	nl, err := netlink.Dial(netlink.Generic)
	if err != nil {
		log.Printf("error connecting to netlink: %s\n", err)
	}
	log.Printf("created generic netlink socket with fd: %d\n", nl.GetFd())

	err = nl.GetFamilyId("acpi_event")
	if err != nil {
		log.Printf("error getting family ID: %s", err)
	}

	// run the bar
	b.Run()
}
