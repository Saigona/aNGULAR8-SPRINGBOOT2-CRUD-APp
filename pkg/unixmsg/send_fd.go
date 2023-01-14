package unixmsg

import (
	"fmt"
	"net"
	"syscall"
)

// https://github.com/mindreframer/golang-stuff/blob/master/github.com/youtube/vitess/go/umgmt/fdpass.go
// see also TestPassFD

func SendFd(conn *net.UnixConn, fd 