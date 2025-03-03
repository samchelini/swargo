package bar

import (
	"strconv"
	"syscall"
)

// block for displaying brightness
type BrightnessBlock struct {
	BlockTemplate
	dir              string
	actualBrightness int
	maxBrightness    int
}

// gets and watches for changes to brightness file
func (block *BrightnessBlock) Run() {
	block.FullText = strconv.Itoa(block.getBrightness())

	// initialize inotify instance
	fd, err := syscall.InotifyInit()
	if err != nil {
		block.LogError("error initializing inotify: " + err.Error())
		return
	}

	// add file to inotify watch list
	_, err = syscall.InotifyAddWatch(fd, block.dir+"/actual_brightness", syscall.IN_MODIFY)
	if err != nil {
		block.LogError("error adding item to watch list: " + err.Error())
		return
	}

	// read events
	var buf [syscall.SizeofInotifyEvent]byte
	for {
		_, err := syscall.Read(fd, buf[:])
		if err != nil {
			block.LogError("error reading event: " + err.Error())
			return
		}
		block.FullText = strconv.Itoa(block.getBrightness())
		block.Update()
	}

	// close inotify
	err = syscall.Close(fd)
	if err != nil {
		block.LogError("error closing inotify: " + err.Error())
	}
}

// calculates the current brightness percentage
func (block *BrightnessBlock) getBrightness() int {
	actualBrightness, err := readFile(block.dir + "/actual_brightness")
	if err != nil {
		block.LogError("error reading actual_brightness: " + err.Error())
	}

	maxBrightness, err := readFile(block.dir + "/max_brightness")
	if err != nil {
		block.LogError("error reading max_brightness: " + err.Error())
	}

	block.actualBrightness, err = strconv.Atoi(actualBrightness[:len(actualBrightness)-1])
	if err != nil {
		block.LogError("error converting actual_brightness to int: " + err.Error())
	}

	block.maxBrightness, err = strconv.Atoi(maxBrightness[:len(maxBrightness)-1])
	if err != nil {
		block.LogError("error converting max_brightness to int: " + err.Error())
	}

	// math trick to round to nearest integer
	brightness := int((float64(block.actualBrightness)/float64(block.maxBrightness)*100 + 0.5))
	return brightness
}

// set brightness directory
func (block *BrightnessBlock) SetDir(dir string) {
	block.dir = dir
}

// opens/reads a file and returns the string
func readFile(filePath string) (string, error) {
	var data string // initialize return string

	// open file as read-only
	fd, err := syscall.Open(filePath, syscall.O_RDONLY, 0)
	if err != nil {
		return data, err
	}
	defer syscall.Close(fd)

	// create a buffer and read the file contents
	var buf [8]byte
	n, err := syscall.Read(fd, buf[:])
	if err != nil {
		return data, err
	}

	// convert file contents to a string
	data = string(buf[:n])
	return data, nil
}
