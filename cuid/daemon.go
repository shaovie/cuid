package main

import (
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"syscall"
)

var pidFd *os.File

func outputPid(pidPath string) {
	if len(pidPath) == 0 {
		return
	}

	os.Mkdir(pidPath, 0755)
	_, name := path.Split(os.Args[0])
	pidPath += "/" + name + ".pid"
	pidFd, err := os.OpenFile(pidPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error - %s output pid %s [%s]", name, pidPath, err)
	}

	// Luanch only one instance.
	if err := syscall.Flock(int(pidFd.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		log.Fatalf("Error - %s flock failed [%s]", name, err)
	}
	pidFd.Truncate(0)
	pidFd.Write([]byte(strconv.Itoa(os.Getpid())))
}

func daemon() {
	isDarwin := runtime.GOOS == "darwin"

	// already a daemon
	if syscall.Getppid() == 1 {
		return
	}

	// fork off the parent process
	ret, ret2, errno := syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if errno != 0 {
		log.Fatalf("fork error! [errno=%d]", errno)
	}

	// failure
	if ret2 < 0 {
		os.Exit(1)
	}

	// handle exception for darwin
	if isDarwin && ret2 == 1 {
		ret = 0
	}

	// if we got a good PID, then we call exit the parent process.
	if int(ret) > 0 {
		os.Exit(0)
	}

	syscall.Umask(0)

	// create a new SID for the child process
	_, err := syscall.Setsid()
	if err != nil {
		log.Fatalf("Error - syscall.Setsid error: %d", err.Error())
	}

	//os.Chdir("/")

	f, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
	if err == nil {
		fd := f.Fd()
		syscall.Dup2(int(fd), int(os.Stdin.Fd()))
		syscall.Dup2(int(fd), int(os.Stdout.Fd()))
		syscall.Dup2(int(fd), int(os.Stderr.Fd()))
	}
}
