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
	fd, _ := syscall.InotifyInit()
	/*
	  if err != nil {
	    log.Printf("error initializing inotify: %s", err)
	  }
	*/

	// add file to inotify watch list
	_, _ = syscall.InotifyAddWatch(fd, "/sys/class/backlight/intel_backlight/actual_brightness", syscall.IN_MODIFY)
	/*
	  if err != nil {
	    log.Printf("error adding item to watch list: %s", err)
	  }

	  log.Printf("wd: %d", wd)
	*/

	// read events
	var buf [syscall.SizeofInotifyEvent]byte
	for {
		_, err := syscall.Read(fd, buf[:])
		if err != nil {
			//log.Printf("error reading event: %s", err)
			break
		}
		block.FullText = strconv.Itoa(block.getBrightness())
		block.Update()
		//log.Printf("n: %d", n)
		//log.Printf("buf: % X", buf)
	}
	_ = syscall.Close(fd)
	/*
	  if err != nil {
	    log.Printf("error closing inotify: %s", err)
	  }
	*/
}

// calculates the current brightness percentage
func (block *BrightnessBlock) getBrightness() int {
	actualBrightness, _ := readFile(block.dir + "/actual_brightness")
	maxBrightness, _ := readFile(block.dir + "/max_brightness")
	block.actualBrightness, _ = strconv.Atoi(actualBrightness[:len(actualBrightness)-1])
	block.maxBrightness, _ = strconv.Atoi(maxBrightness[:len(maxBrightness)-1])

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
