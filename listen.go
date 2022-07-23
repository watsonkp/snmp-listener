package main

import "fmt"
import "net"
import "encoding/asn1"

type VarBind struct {
	Name asn1.ObjectIdentifier
	Value asn1.RawValue
}

type GetRequest struct {
	RequestID int64
	ErrorStatus int64
	ErrorIndex int64
	VarBindList []VarBind
}

type TrapPDU struct {
	Enterprise asn1.ObjectIdentifier
	AgentAddr []byte `asn1:"application,tag:0"`
	GenericTrap int64
	SpecificTrap int64
	TimeStamp int64 `asn1:"application,tag:3"`
	VariableBindings []VarBind
}

type SNMPMessage struct {
	Version int64
	Community []byte
	Data TrapPDU `asn1:"tag:4"`
	//Data GetRequest `asn1:"tag:0"`
}

func main() {
	//addr := net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 10162}
	addr := net.UDPAddr{IP: net.IPv4(10, 42, 0, 156), Port: 162}

	conn, err := net.ListenUDP("udp4", &addr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	fmt.Printf("Listening at: %s\n", addr.String())
	channel := make(chan error, 1)

	go func() {
		buf := make([]byte, 1024)
		// TODO: Go func for each read?
		for {
			// TODO: Experiment with ReadMsgUDP
			n, address, err := conn.ReadFromUDP(buf)
			if err != nil {
				channel <- err
			}
			fmt.Printf("Read %d bytes from %v\n", n, address)
			fmt.Printf("%x\n", buf[0:n])
			if n == 0 {
				continue
			}
			val := new(SNMPMessage)
			_, err = asn1.Unmarshal(buf, val)
			if err != nil {
				fmt.Println(err)
				channel <- err
				continue
			}
			fmt.Printf("%v\n", val)
		}
	}()

	select {
	case err := <-channel:
		fmt.Println(err)
	}
}
