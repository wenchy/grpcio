// +build !windows

package daemon

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wenchy/grpcio/internal/daemon/pidfile"
)

// SignalHandlerFunc is the interface for signal handler functions.
// type SignalHandlerFunc func(sig os.Signal) (err error)
type SignalHandlerFunc func(sig os.Signal) error

type arguments struct {
	Conf   string
	ID     string
	Daemon bool
	Action string
}

var Args arguments

func init() {
	flag.StringVar(&Args.Conf, "conf", "./conf.yaml", "server conf path")
	flag.StringVar(&Args.ID, "id", "0.0.0.0", "server node string id")
	flag.BoolVar(&Args.Daemon, "daemon", false, "run as a daemon")
	flag.StringVar(&Args.Action, "action", "start", "action options: start, stop, restart and reload")
}

var globalStdoutFile *os.File

func initStdoutFile(path string) error {
	// log.Println("init stdout file in unix mode")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	globalStdoutFile = f
	if err != nil {
		println(err)
		return err
	}
	if err = syscall.Dup2(int(f.Fd()), int(os.Stdout.Fd())); err != nil {
		return err
	}
	if err = syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd())); err != nil {
		return err
	}
	return nil
}

func Start() {
	flag.Parse()

	if err := initStdoutFile("./stdout.log"); err != nil {
		log.Fatalln("initStdoutFile failed", err)
	}

	switch Args.Action {
	case "start":
		fmt.Println("starting...")
	case "stop":
		err := pidfile.StopProcess()
		if err != nil {
			fmt.Printf("stop failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "restart":
		err := pidfile.StopProcess()
		if err != nil {
			fmt.Printf("stop failed: %v\n", err)
			os.Exit(1)
		}
		// restart: continue to start
		for i, arg := range os.Args {
			if arg == "-action=restart" {
				os.Args[i] = "-action=start"
				break
			}
		}
	case "reload":
		fmt.Println("reloading...")
		err := pidfile.ReloadProcess()
		if err != nil {
			fmt.Printf("reload failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	default:
		fmt.Printf("unknown action: %s\n", Args.Action)
		os.Exit(1)
	}

	if Args.Daemon {
		// 守护进程
		daemonize(1, 1)
	}
	err := pidfile.Write()
	if err != nil {
		fmt.Println("already running!", err)
		os.Exit(2)
	}
}

func Fini() {
	pidfile.Release()
}

func ProcessSignal(ctx context.Context, stopHandler, reloadHandler SignalHandlerFunc) {
	// Wait for signals
	// 	1. gracefully shutdown the server with a timeout of 5 seconds.
	// 	2. SIGHUP: reload config
	sigChan := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	defer signal.Stop(sigChan)
	for {
		var sig os.Signal
		stopped := false
		select {
		case sig = <-sigChan:
			fmt.Println("received signal: ", sig)
			switch sig {
			case syscall.SIGINT:
				stopped = true
			case syscall.SIGTERM:
				stopped = true
			case syscall.SIGQUIT:
				stopped = true
			case syscall.SIGPIPE:
				stopped = true
			case syscall.SIGHUP:
				// fmt.Println("need reload")
				reloadHandler(sig)
			default:
				fmt.Println("unknown signal: ", sig)
			}
		case <-ctx.Done():
			stopped = true
			break
		}
		if stopped {
			err := stopHandler(sig)
			if err != nil {
				fmt.Println("stop success, but trigger error: ", err)
				os.Exit(2)
			}
			fmt.Println("stop success")
			break
		}
	}
}

func daemonize(nochdir, noclose int) (int, error) {
	// already a daemon
	if syscall.Getppid() == 1 {
		/* Change the file mode mask */
		syscall.Umask(0)

		if nochdir == 0 {
			os.Chdir("/")
		}

		return 0, nil
	}

	files := make([]*os.File, 3, 6)
	if noclose == 0 {
		nullDev, err := os.OpenFile("/dev/null", 0, 0)
		if err != nil {
			return 1, err
		}
		files[0], files[1], files[2] = nullDev, nullDev, nullDev
	} else {
		files[0], files[1], files[2] = os.Stdin, os.Stdout, os.Stderr
	}

	dir, _ := os.Getwd()
	sysattrs := syscall.SysProcAttr{Setsid: true}
	attrs := os.ProcAttr{Dir: dir, Env: os.Environ(), Files: files, Sys: &sysattrs}

	proc, err := os.StartProcess(os.Args[0], os.Args, &attrs)
	if err != nil {
		return -1, fmt.Errorf("can't create process %s: %s", os.Args[0], err)
	}
	proc.Release()
	os.Exit(0)

	return 0, nil
}
