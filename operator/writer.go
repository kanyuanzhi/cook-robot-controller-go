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
	for _, addressValue := range addressValueList[1:] {
		w.tcpServer.Write(addressValue.Address, uint64(addressValue.Value))
		//time.Sleep(20 * time.Millisecond) // 每条指令间隔20ms发送
	}
	w.tcpServer.Write(addressValueList[0].Address, uint64(addressValueList[0].Value)) // 最后发送控制字
	successChan <- true
}
