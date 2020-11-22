package PortForward

import (
	"bufio"
	"fmt"
	"github.com/shoriwe/FullProxy/pkg/ConnectionHandlers"
	"github.com/shoriwe/FullProxy/pkg/Proxies/Basic"
	"github.com/shoriwe/FullProxy/pkg/Sockets"
	"net"
)

type LocalForward struct {
	TargetHost string
	TargetPort string
}

func (localForward *LocalForward) SetAuthenticationMethod(_ ConnectionHandlers.AuthenticationFunction) error {
	return nil
}

func (localForward *LocalForward)Handle(
	clientConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer) error{
	targetConnection, connectionError := Sockets.Connect(&localForward.TargetHost, &localForward.TargetPort)
	if connectionError == nil {
		targetReader, targetWriter := Sockets.CreateSocketConnectionReaderWriter(targetConnection)
		Basic.Proxy(
			clientConnection, targetConnection,
			clientConnectionReader, clientConnectionWriter,
			targetReader, targetWriter)
	} else {
		_ = clientConnection.Close()
	}
	return connectionError
}

func StartLocalPortForward(targetHost *string, targetPort *string, masterAddress *string, masterPort *string) {
	if !(*targetHost == "" || *targetPort == "" || *masterAddress == "" || *masterPort == "") {
		localForward := LocalForward{
			TargetHost: *targetHost,
			TargetPort: *targetPort,
		}
		ConnectionHandlers.GeneralSlave(masterAddress, masterPort, &localForward)
	} else {
		fmt.Println("All flags need to be in use")
	}
}
