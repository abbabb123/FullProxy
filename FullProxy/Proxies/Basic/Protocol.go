package Basic

import (
	"FullProxy/FullProxy/Sockets"
	"bufio"
<<<<<<< Updated upstream:FullProxy/Proxies/Basic/Protocol.go
=======
	"fmt"
	"github.com/shoriwe/FullProxy/src/Sockets"
>>>>>>> Stashed changes:src/Proxies/Basic/Protocol.go
	"net"
	"time"
)


func HandleReadWrite(
	sourceConnection net.Conn, sourceReader *bufio.Reader, destinationWriter *bufio.Writer, connectionAlive *bool){
	for {
		if !(*connectionAlive){
			break
		}
		_ = sourceConnection.SetReadDeadline(time.Now().Add(10 * time.Second))
		numberOfBytesReceived, buffer, ConnectionError := Sockets.Receive(sourceReader, 20480)
		if ConnectionError != nil {
			if ConnectionError, ok := ConnectionError.(net.Error); !(ok && ConnectionError.Timeout()) {
				break
			}
		}
		if numberOfBytesReceived > 0 {
			_, ConnectionError = Sockets.Send(destinationWriter, buffer[:numberOfBytesReceived])
			if ConnectionError != nil {
				break
			}}
		if numberOfBytesReceived > 0{
			fmt.Println(buffer[:numberOfBytesReceived])
		}
		buffer = nil
	}
	_ = sourceConnection.Close()
	*connectionAlive = false
}


func Proxy(clientConnection net.Conn,
	targetConnection net.Conn,
	clientConnectionReader *bufio.Reader,
	clientConnectionWriter *bufio.Writer,
	targetConnectionReader *bufio.Reader,
	targetConnectionWriter *bufio.Writer) {
	connectionAlive := true

	go HandleReadWrite(clientConnection, clientConnectionReader, targetConnectionWriter, &connectionAlive)
	go HandleReadWrite(targetConnection, targetConnectionReader, clientConnectionWriter, &connectionAlive)
}
