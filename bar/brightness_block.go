package bar


import (
  "syscall"
  "strconv"
)

// block for displaying brightness
type BrightnessBlock struct {
	BlockTemplate
  dir string
  actualBrightness int
  maxBrightness int
}

// gets and watches for changes to brightness file
func (b *BrightnessBlock) Run() {
  b.init()

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
    b.init()
    b.Update()
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

func (b *BrightnessBlock) init() {
  actualBrightness, _ := readFile(b.dir + "/actual_brightness")
  maxBrightness, _ := readFile(b.dir + "/max_brightness")
  b.actualBrightness, _ = strconv.Atoi(actualBrightness[:len(actualBrightness)-1])
  b.maxBrightness, _ = strconv.Atoi(maxBrightness[:len(maxBrightness)-1])
  brightness := (float64(b.actualBrightness) / float64(b.maxBrightness) * 100 + 0.5)
  b.FullText = strconv.Itoa(int(brightness))
}

// set brightness directory
func (b *BrightnessBlock) SetDir(dir string) {
  b.dir = dir
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

