package bar

import (
	"github.com/samchelini/swargo/filewatcher"
	"os"
	"strconv"
)

// block for displaying brightness
type BrightnessBlock struct {
	BlockTemplate
	dir    string
	prefix string
}

// gets and watches for changes to brightness file
func (block *BrightnessBlock) Run() {
	// initialize text to current brightness
	block.SetFullText(block.prefix, strconv.Itoa(block.getBrightness()))
	block.Update()

	// initialize filewatcher
	fw, err := filewatcher.NewFileWatcher()
	if err != nil {
		block.LogError("error initializing inotify: " + err.Error())
		return
	}

	// add file to watch list for modify events
	err = fw.AddWatch(block.dir+"/actual_brightness", filewatcher.IN_MODIFY)
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
		block.SetFullText(block.prefix, strconv.Itoa(block.getBrightness()))
		block.Update()
	}

	// close filewatcher
	err = fw.Close()
	if err != nil {
		block.LogError("error closing filewatcher: " + err.Error())
	}
}

// calculates the current brightness percentage
func (block *BrightnessBlock) getBrightness() int {
	var buf []byte
	var actualBrightness float64
	var maxBrightness float64

	// get actual_brightness
	buf, err := os.ReadFile(block.dir + "/actual_brightness")
	if err != nil {
		block.LogError("error reading actual_brightness: " + err.Error())
	}

	actualBrightness, err = strconv.ParseFloat(string(buf[:len(buf)-1]), 64)
	if err != nil {
		block.LogError("error converting actual_brightness to float: " + err.Error())
	}

	// get max_brightness
	buf, err = os.ReadFile(block.dir + "/max_brightness")
	if err != nil {
		block.LogError("error reading max_brightness: " + err.Error())
	}

	maxBrightness, err = strconv.ParseFloat(string(buf[:len(buf)-1]), 64)
	if err != nil {
		block.LogError("error converting max_brightness to float: " + err.Error())
	}

	// math trick to round to nearest integer
	return int((actualBrightness/maxBrightness)*100 + 0.5)
}

// set brightness directory
func (block *BrightnessBlock) SetDir(dir string) {
	block.dir = dir
}

// set prefix before the brightness percentage
func (block *BrightnessBlock) SetPrefix(prefix string) {
	block.prefix = prefix
}
