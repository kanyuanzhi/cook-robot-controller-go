package operator

import (
	"cook-robot-controller-go/data"
	"net"
	"time"
)

type Writer struct {
	connection net.Conn
}

func NewWriter() *Writer {
	writer := &Writer{}
	return writer
}

func (w *Writer) Run() {
	//serverAddr := "127.0.0.1:10001"
	//conn, err := net.Dial("tcp", serverAddr)
	//if err != nil {
	//	logger.Log.Println("无法建立TCP连接:", err)
	//}
	//w.connection = conn
}

func (w *Writer) Send(successChan chan bool, addressValueList []*data.AddressValue) {
	//addressValueList := a.GetAddressValueList()
	//logger.Log.Println(a.ShowAddressValueList())
	//w.connection.Write([]byte("123"))
	time.Sleep(100 * time.Millisecond)
	//logger.Log.Println(addressValueList)
	successChan <- true
}
