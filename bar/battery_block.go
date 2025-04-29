package bar

import (
	"github.com/samchelini/swargo/filewatcher"
	"os"
	"strconv"
)

// block for displaying brightness
type BatteryBlock struct {
	BlockTemplate
	dir               string
	chargingPrefix    string
	dischargingPrefix string
}

// gets and watches for changes to brightness file
func (block *BatteryBlock) Run() {
	// initialize text to current charge and status
	block.SetFullText(block.getStatus(), strconv.Itoa(block.getCharge()))
	block.Update()

	// initialize filewatcher
	fw, err := filewatcher.NewFileWatcher()
	if err != nil {
		block.LogError("error initializing inotify: " + err.Error())
		return
	}

	// add file to watch list for modify events
	err = fw.AddWatch(block.dir+"/status", filewatcher.IN_ALL_EVENTS)
	if err != nil {
		block.LogError("error adding item to watch list: " + err.Error())
		return
	}

	// watch for events
	for {
		err := fw.Watch()
		if err != nil {
			block.LogError("error reading event: " + err.Error())
			return
		}
		block.SetFullText(block.getStatus(), strconv.Itoa(block.getCharge()))
		block.Update()
	}

	// close filewatcher
	err = fw.Close()
	if err != nil {
		block.LogError("error closing filewatcher: " + err.Error())
	}
}

// get battery charging/discharging status
func (block *BatteryBlock) getStatus() string {
	var buf []byte
	var status string

	// get status
	buf, err := os.ReadFile(block.dir + "/status")
	if err != nil {
		block.LogError("error reading battery status: " + err.Error())
	}
	status = string(buf[:len(buf)-1])

	// return the specified prefix
	if status == "Charging" {
		return block.chargingPrefix
	} else {
		return block.dischargingPrefix
	}
}

// get charge
func (block *BatteryBlock) getCharge() int {
	var buf []byte
	var capacity int

	// get charge
	buf, err := os.ReadFile(block.dir + "/capacity")
	if err != nil {
		block.LogError("error reading battery charge: " + err.Error())
	}

	// convert charge to integer
	capacity, err = strconv.Atoi(string(buf[:len(buf)-1]))
	if err != nil {
		block.LogError("error converting battery charge to int: " + err.Error())
	}

	return capacity
}

// set battery directory
func (block *BatteryBlock) SetDir(dir string) {
	block.dir = dir
}

// set charging prefix
func (block *BatteryBlock) SetChargingPrefix(prefix string) {
	block.chargingPrefix = prefix
}

// set discharging prefix
func (block *BatteryBlock) SetDischargingPrefix(prefix string) {
	block.dischargingPrefix = prefix
}
