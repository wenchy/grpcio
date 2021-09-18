// +build !windows

package pidfile

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var (
	errNotConfigured = errors.New("pidfile not configured")
	errStopTimeout   = errors.New("stop process timeout")
	pidfile          = flag.String("pid-file", ".mypid", "if specified, write pid to this file")
	fileHandle       *os.File
)

// GetPidfilePath returns the configured pidfile path.
func GetPidfilePath() string {
	return *pidfile
}

// SetPidfilePath sets the pidfile path.
func SetPidfilePath(p string) {
	*pidfile = p
}

// Write the pidfile based on the flag.
// It is an error if the pidfile hasn't been configured.
func Write() error {
	if *pidfile == "" {
		return errNotConfigured
	}

	if err := os.MkdirAll(filepath.Dir(*pidfile), os.FileMode(0755)); err != nil {
		return err
	}

	file, err := os.OpenFile(*pidfile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening pidfile %s: %s", *pidfile, err)
	}
	// intentionally not closed

	// try to lock file
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return fmt.Errorf("error lock pidfile %s: %s", *pidfile, err)
	}

	// truncate file
	err = syscall.Ftruncate(int(file.Fd()), 0)
	if err != nil {
		return err
	}

	// write pid to file
	_, err = fmt.Fprintf(file, "%d\n", os.Getpid())
	if err != nil {
		return err
	}
	file.Sync()

	fileHandle = file

	return nil
}

// Release pidfile resource
func Release() error {
	defer fileHandle.Close()
	// // write 3276800 to file
	// _, err := fmt.Fprintf(fileHandle, "%d", 3276800)
	// if err != nil {
	// 	return err
	// }

	// unlock
	err := syscall.Flock(int(fileHandle.Fd()), syscall.LOCK_UN)
	if err != nil {
		return err
	}
	return nil
}

// Read the pid from the configured file.
// It is an error if the pidfile hasn't been configured.
func Read() (int, error) {
	if *pidfile == "" {
		return 0, errNotConfigured
	}

	d, err := ioutil.ReadFile(*pidfile)
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(string(bytes.TrimSpace(d)))
	if err != nil {
		return 0, fmt.Errorf("failed to parsing pid from %s: %s", *pidfile, err)
	}

	return pid, nil
}

func GetExecutableByPid(pid int) (string, error) {
	return os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
}

func StopProcess() error {
	pid, err := Read()
	if err != nil {
		return err
	}
	// On Unix systems, FindProcess always succeeds and returns a Process
	// for the given pid, regardless of whether the process exists.
	process, err := os.FindProcess(pid)
	if err != nil {
		// process already finished
		return nil
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		// fmt.Println("signal SIGTERM error: ", err)
		// process already finished? or permission denied(app process impossivle)
		// NOTE: we cannot compare the error type to check
		return nil
	}
	// wait 3 seconds, send 30
	for i := 0; i < 30; i++ {
		// If sig is 0, then no signal is sent, but error checking is still performed.
		// this can be used to check for the existence of a process ID or process group ID.
		// refer: https://unix.stackexchange.com/questions/169898/what-does-kill-0-do
		err = process.Signal(syscall.Signal(0))
		if err != nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errStopTimeout
}

func ReloadProcess() error {
	pid, err := Read()
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = process.Signal(syscall.SIGHUP)
	if err != nil {
		return err
	}
	// wait 3 seconds, send 30
	for i := 0; i < 30; i++ {
		// If sig is 0, then no signal is sent, but error checking is still performed.
		// this can be used to check for the existence of a process ID or process group ID.
		// refer: https://unix.stackexchange.com/questions/169898/what-does-kill-0-do
		err = process.Signal(syscall.Signal(0))
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}
