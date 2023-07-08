package operator

import (
	"cook-robot-controller-go/data"
	"cook-robot-controller-go/modbus"
	"time"
)

type Writer struct {
	tcpServer *modbus.TCPServer
}

func NewWriter(tcpServer *modbus.TCPServer) *Writer {
	writer := &Writer{
		tcpServer: tcpServer,
	}
	return writer
}

func (w *Writer) Send(successChan chan bool, addressValueList []*data.AddressValue) {
	time.Sleep(50 * time.Millisecond)
	for _, addressValue := range addressValueList {
		w.tcpServer.Write(addressValue.Address, uint64(addressValue.Value))
	}
	successChan <- true
}
