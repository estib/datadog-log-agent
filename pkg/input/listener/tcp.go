// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

package listener

import (
	"fmt"
	"log"
	"net"

	"github.com/DataDog/datadog-log-agent/pkg/config"
	"github.com/DataDog/datadog-log-agent/pkg/pipeline"
)

// A TcpListener listens to bytes on a tcp connection and sends log lines to
// an output channel
type TcpListener struct {
	listener net.Listener
	anl      *AbstractNetworkListener
}

// NewTcpListener returns an initialized NewTcpListener
func NewTcpListener(pp *pipeline.PipelineProvider, source *config.IntegrationConfigLogSource) (*AbstractNetworkListener, error) {
	log.Println("Starting TCP forwarder on port", source.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", source.Port))
	if err != nil {
		return nil, err
	}
	tcpListener := &TcpListener{
		listener: listener,
	}
	anl := &AbstractNetworkListener{
		listener: tcpListener,
		pp:       pp,
		source:   source,
	}
	tcpListener.anl = anl
	return anl, nil
}

// run lets the listener handle incoming tcp connections
func (tcpListener *TcpListener) run() {
	for {
		conn, err := tcpListener.listener.Accept()
		if err != nil {
			log.Println("Can't listen:", err)
			return
		}
		go tcpListener.anl.handleConnection(conn)
	}
}

func (tcpListener *TcpListener) readMessage(conn net.Conn, inBuf []byte) (int, error) {
	return conn.Read(inBuf)
}
