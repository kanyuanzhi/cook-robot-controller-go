package modbus

import (
	"cook-robot-controller-go/logger"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TCPServer struct {
	host string
	port uint16
	conn net.Conn
	//RealtimeValueMap     map[string]uint32
	RealtimeValueSyncMap *sync.Map

	PauseReadChan chan bool

	//mu sync.Mutex

	debugMode bool
}

func NewTCPServer(host string, port uint16, debugMode bool) *TCPServer {
	return &TCPServer{
		host: host,
		port: port,
		//RealtimeValueMap:     make(map[string]uint32),
		RealtimeValueSyncMap: new(sync.Map),

		PauseReadChan: make(chan bool),
		debugMode:     debugMode,
	}
}

func (t *TCPServer) Run() {
	if t.debugMode {
		//t.RealtimeValueMap["DD232"] = 4501
		//t.RealtimeValueMap["DD234"] = 300
		t.SetRealtimeValue("DD232", uint32(4510))
		t.SetRealtimeValue("DD234", uint32(300))
		logger.Log.Println("TCP服务以测试模式启动，无TCP连接建立")
		return
	}

	timeout := 2 * time.Second
	// 创建一个 Dialer，并设置超时时间
	dialer := net.Dialer{
		Timeout: timeout,
	}
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", t.host, t.port))
	if err != nil {
		logger.Log.Println("无法建立TCP连接:", err)
	}
	t.conn = conn
	logger.Log.Println("TCP服务启动，与下位机建立TCP连接")

	//for i := 2050; i < 2173; i += 2 {
	//	t.Write(fmt.Sprintf("DD%d", i), 100000)
	//}

	//t.Write("DD0", 1)
	//t.Write("DD2", 1)
	//t.Read("DS2050", 120)

	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			t.Read("DS200", 120)
			//t.Read("DS2050", 120)
		case <-t.PauseReadChan:
			<-t.PauseReadChan
		}
	}
	//t.Read("Dd500", 20)
}

func (t *TCPServer) Write(prefixAddress string, value uint64) {
	t.PauseReadChan <- true
	defer func() {
		t.PauseReadChan <- true
	}()
	transmissionIdentifier := "0422"
	protocolIdentifier := "0000" // 协议标识符
	salveNum := "01"

	var bytesLen string
	var order string
	var valueStr string
	var valueStrLow string
	var valueStrHigh string
	var registerNum string
	var dataBytesLen string

	valueStr = strconv.FormatUint(value, 16)
	if len(valueStr) < 8 {
		zeros := strings.Repeat("0", 8-len(valueStr))
		valueStr = zeros + valueStr
	}
	valueStr = strings.ToUpper(valueStr)
	//logger.Log.Println(valueStr)

	prefixAddress = strings.ToUpper(prefixAddress)
	prefix := prefixAddress[:2]

	if strings.Index(prefix, "S") > -1 { // 单字
		bytesLen = "0006"
		order = "06"
	} else { // 双字
		bytesLen = "000B"
		order = "10"
		valueStrLow = valueStr[4:]
		valueStrHigh = valueStr[:4]
		registerNum = "0002" // 写入的寄存器个数
		dataBytesLen = "04"  // 写入的数值字节数
	}

	address := prefixAddress[2:]
	addressNum, _ := strconv.ParseUint(address, 10, 16)
	var addressStr string
	if strings.Index(prefix, "H") > -1 { // 断电寄存器
		addressStr = strconv.FormatUint(addressNum+41088, 16)
	} else { // 普通寄存器
		addressStr = strconv.FormatUint(addressNum, 16)
	}
	if len(addressStr) < 4 {
		zeros := strings.Repeat("0", 4-len(addressStr))
		addressStr = zeros + addressStr
	}
	addressStr = strings.ToUpper(addressStr)

	CMD := encode(transmissionIdentifier, 2)
	CMD = append(CMD, encode(protocolIdentifier, 2)...)
	CMD = append(CMD, encode(bytesLen, 2)...)
	CMD = append(CMD, encode(salveNum, 1)...)
	CMD = append(CMD, encode(order, 1)...)
	CMD = append(CMD, encode(addressStr, 2)...)
	if strings.Index(prefix, "S") > -1 {
		CMD = append(CMD, encode(valueStr, 2)...)
	} else {
		CMD = append(CMD, encode(registerNum, 2)...)
		CMD = append(CMD, encode(dataBytesLen, 1)...)
		CMD = append(CMD, encode(valueStrLow, 2)...)
		CMD = append(CMD, encode(valueStrHigh, 2)...)
	}
	//logger.Log.Println(CMD)
	_, err := t.conn.Write(CMD)
	if err != nil {
		logger.Log.Println(err)
		return
	}

	buffer := make([]byte, 12)
	_, err = t.conn.Read(buffer)
	//logger.Log.Println(string(buffer))
	if err != nil {
		logger.Log.Printf("读取数据失败：%v\n", err)
		return
	}
	// 处理接收到的数据
	//data := buffer[:n]
	//fmt.Printf("接收到的数据：%s\n", string(data))
	//fmt.Println(data)
}

// prefixAddress DD21 DS21 HD21 HS21
func (t *TCPServer) Read(prefixAddress string, size uint64) {
	transmissionIdentifier := "0422"
	protocolIdentifier := "0000" // 协议标识符
	salveNum := "01"
	bytesLen := "0006" // 字节长度
	order := "03"      // 功能码，寄存器读

	prefixAddress = strings.ToUpper(prefixAddress)
	prefix := prefixAddress[:2]
	address := prefixAddress[2:]
	addressNum, _ := strconv.ParseUint(address, 10, 16)

	var addressStr string
	if strings.Index(prefix, "H") > -1 { // 断电寄存器
		addressStr = strconv.FormatUint(addressNum+41088, 16)
	} else { // 普通寄存器
		addressStr = strconv.FormatUint(addressNum, 16)
	}
	if len(addressStr) < 4 {
		zeros := strings.Repeat("0", 4-len(addressStr))
		addressStr = zeros + addressStr
	}
	addressStr = strings.ToUpper(addressStr)

	var sizeStr string
	if strings.Index(prefix, "S") > -1 { // 单字
		sizeStr = strconv.FormatUint(size, 16)
	} else { // 双字
		size = 2
		sizeStr = strconv.FormatUint(size, 16)
	}

	if len(sizeStr) < 4 {
		zeros := strings.Repeat("0", 4-len(sizeStr))
		sizeStr = zeros + sizeStr
	}
	sizeStr = strings.ToUpper(sizeStr)

	CMD := encode(transmissionIdentifier, 2)
	CMD = append(CMD, encode(protocolIdentifier, 2)...)
	CMD = append(CMD, encode(bytesLen, 2)...)
	CMD = append(CMD, encode(salveNum, 1)...)
	CMD = append(CMD, encode(order, 1)...)
	CMD = append(CMD, encode(addressStr, 2)...)
	CMD = append(CMD, encode(sizeStr, 2)...)

	_, err := t.conn.Write(CMD)
	if err != nil {
		logger.Log.Println(err)
		return
	}

	buffer := make([]byte, 9+size*2)
	_, err = t.conn.Read(buffer)
	bufferHexStr := hex.EncodeToString(buffer)

	// logger.Log.Printf("buffer长度%d", len(buffer))

	if err != nil {
		logger.Log.Printf("读取数据失败：%v\n", err)
		return
	}
	// 处理接收到的数据
	//data := buffer[:n]
	//logger.Log.Println(data)
	var i uint64
	//t.mu.Lock()
	//defer t.mu.Unlock()
	for i = 0; i < size; i = i + 2 {
		value, _ := strconv.ParseInt(bufferHexStr[18+4*(i+1):18+4*(i+1)+4]+bufferHexStr[18+4*i:18+4*i+4], 16, 64)
		//value, _ := strconv.ParseInt(string(data[9+2*i:9+2*i+2]), 16, 64)
		//t.RealtimeValueMap[fmt.Sprintf("DD%d", addressNum+i)] = uint32(value)
		t.SetRealtimeValue(fmt.Sprintf("DD%d", addressNum+i), uint32(value))
		//t.RealtimeValueMap[fmt.Sprintf("%s%d", prefix, addressNum+i)] = uint32(value)
		//logger.Log.Printf("%s%d:%d", prefix, addressNum+i, value)
	}
	//logger.Log.Printf("%d", t.RealtimeValueMap[fmt.Sprintf("DD%d", 212)])
}

func encode(numStr string, length uint8) []byte {
	num, _ := strconv.ParseInt(numStr, 16, 32)
	if length == 1 {
		return []byte{uint8(num)}
	}
	bytes := make([]byte, length)
	binary.BigEndian.PutUint16(bytes, uint16(num))
	return bytes
}

func (t *TCPServer) SetRealtimeValue(address string, value uint32) {
	t.RealtimeValueSyncMap.Store(address, value)
}

func (t *TCPServer) GetRealtimeValue(address string) (uint32, bool) {
	valueAny, has := t.RealtimeValueSyncMap.Load(address)
	if !has {
		logger.Log.Printf("当前RealtimeValueSyncMap中不存在%f值", address)
		return 0, false
	}
	if value, ok := valueAny.(uint32); ok {
		return value, true
	} else {
		logger.Log.Printf("当前RealtimeValueSyncMap中的%f值不为uint32类型", address)
		return 0, false
	}
}
