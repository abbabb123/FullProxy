package SOCKS5

import (
	"bytes"
	"github.com/shoriwe/FullProxy/src/ConnectionStructures"
	"github.com/shoriwe/FullProxy/src/Hashing"
	"github.com/shoriwe/FullProxy/src/Sockets"
)


func HandleUsernamePasswordAuthentication(clientConnectionReader ConnectionStructures.SocketReader,
	username *[]byte,
	passwordHash *[]byte) (bool, byte){
	numberOfReceivedBytes, credentials, connectionError := Sockets.Receive(clientConnectionReader, 1024)
	if connectionError != nil{
		return false, 0
	}
	if numberOfReceivedBytes < 4{
		return false, 0
	}

	if credentials[0] != BasicNegotiation {
		return false, 0
	}
	receivedUsernameLength := int(credentials[1])
	if receivedUsernameLength + 3  >= numberOfReceivedBytes{
		return false, 0
	}
	receivedUsername := credentials[2:2+receivedUsernameLength]
	if bytes.Equal(receivedUsername, *username){
		rawReceivedUsernamePassword := credentials[2+receivedUsernameLength+1:numberOfReceivedBytes]
		if bytes.Equal(Hashing.Sha3_512_256(rawReceivedUsernamePassword), *passwordHash){
			return true, BasicNegotiation
		}
	}
	return false, 0
}