package pidfile

import (
        "os"
	"testing"
)

func TestNotSetPidfileAndRead(t *testing.T) {
    SetPidfilePath("")
_, err := Read()
         if err != errNotConfigured {
             t.Fatalf("should not read before calls SetPidfilePath")
         }
}

func TestNotSetPidfileAndWrite(t *testing.T) {
    SetPidfilePath("")
err := Write()
         if err != errNotConfigured {
             t.Fatalf("should not write before calls SetPidfilePath")
         }
}

func TestWriteTwice(t *testing.T) {
	SetPidfilePath("./.pid")

	err := Write()
	if err != nil {
		t.Fatalf("error write pidfile: %s", err)
	}

	err = Write()
	if err == nil {
		t.Fatalf("error write pidfile twice")
	} else {
		t.Log(err)
	}

	pid, err := Read()
	    if err != nil {
            t.Fatalf("error read pidfile: %s", err)
        }

    if pid != os.Getpid() {
        t.Fatalf("pidfile's pid:%d not match os.GetPid():%d", pid, os.Getpid())
    } else {
        t.Log("pid:", pid)
    }

    os.Remove("./.pid")
}