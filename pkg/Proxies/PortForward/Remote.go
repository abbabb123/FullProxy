package PortForward

import (
	"bufio"
	"crypto/tls"
	"errors"
	"github.com/shoriwe/FullProxy/pkg/Proxies/RawProxy"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"github.com/shoriwe/FullProxy/pkg/Templates"
	"github.com/shoriwe/FullProxy/pkg/Templates/Types"
	"net"
	"time"
)

type RemoteForward struct {
	MasterHost       string
	MasterPort       string
	TLSConfiguration *tls.Config
	LoggingMethod    Types.LoggingMethod
	Tries            int
	Timeout          time.Duration
	InboundFilter    Types.IOFilter
}

func (remoteForward *RemoteForward) SetLoggingMethod(loggingMethod Types.LoggingMethod) error {
	remoteForward.LoggingMethod = loggingMethod
	return nil
}

func (remoteForward *RemoteForward) SetTries(tries int) error {
	remoteForward.Tries = tries
	return nil
}

func (remoteForward *RemoteForward) SetTimeout(timeout time.Duration) error {
	remoteForward.Timeout = timeout
	return nil
}

func (remoteForward *RemoteForward) SetInboundFilter(filter Types.IOFilter) error {
	remoteForward.InboundFilter = filter
	return nil
}

func (remoteForward *RemoteForward) Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error {
	if !Templates.FilterInbound(remoteForward.InboundFilter, Templates.ParseIP(clientConnection.RemoteAddr().String())) {
		errorMessage := "Connection denied to: " + clientConnection.RemoteAddr().String()
		Templates.LogData(remoteForward.LoggingMethod, errorMessage)
		_ = clientConnection.Close()
		return errors.New(errorMessage)
	}
	Templates.LogData(remoteForward.LoggingMethod, "Connection Received from: ", clientConnection.RemoteAddr().String())
	targetConnection, connectionError := Sockets.TLSConnect(
		&remoteForward.MasterHost,
		&remoteForward.MasterPort,
		(*remoteForward).TLSConfiguration)
	if connectionError != nil {
		Templates.LogData(remoteForward.LoggingMethod, connectionError)
		_ = clientConnection.Close()
		return connectionError
	}
	targetConnectionReader, targetConnectionWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
	rawProxy := RawProxy.RawProxy{
		TargetConnection:       targetConnection,
		TargetConnectionReader: targetConnectionReader,
		TargetConnectionWriter: targetConnectionWriter,
	}
	_ = rawProxy.SetTries(remoteForward.Tries)
	_ = rawProxy.SetTimeout(remoteForward.Timeout)
	_ = rawProxy.SetLoggingMethod(remoteForward.LoggingMethod)
	return rawProxy.Handle(clientConnection, clientConnectionReader, clientConnectionWriter)
}

func (remoteForward *RemoteForward) SetAuthenticationMethod(_ Types.AuthenticationMethod) error {
	return errors.New("This kind of proxy doesn't support authentication methods")
}
