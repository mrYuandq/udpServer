package udpServer

import (
	"net"
	"os"
	"udpServer/log"
)

var logger = log.NewLogger(os.Stdout)

type UdpServer struct {
	Addr    *net.UDPAddr
	Conn    *net.UDPConn
	NetWork string
	Address string

	Handler func()
}

func (us *UdpServer) StartUdpListen() (err error) {
	err = us.beginResolveUDPAddr()
	errorHandler(err)
	err = us.beginListenUDP()
	errorHandler(err)

	logger.Info("listen udp address: ", us.Address, " network: ", us.NetWork)

	us.Handler()
	return
}

/**
 * 监听端口
 */
func (us *UdpServer) beginResolveUDPAddr() (err error) {
	us.Addr, err = net.ResolveUDPAddr(us.NetWork, us.Address)
	return
}

/**
 * 添加 Udp network 监听
 */
func (us *UdpServer) beginListenUDP() (err error) {
	us.Conn, err = net.ListenUDP(us.NetWork, us.Addr)
	return
}

/**
 * 写入到
 */
func (us *UdpServer) WriteToUDP(b []byte, addr *net.UDPAddr) (i int, e error) {
	i, e = us.Conn.WriteToUDP(b, addr)
	return
}

/**
 * 读出
 * 每次读出后, 便更新设备信息
 */
func (us *UdpServer) ReadFromUDP(b []byte) (i int, addr *net.UDPAddr, err error) {
	i, addr, err = us.Conn.ReadFromUDP(b)
	// Devices.UpdateDevice(addr)
	return
}

/**
 * 开启心跳检查
 */
func (us *UdpServer) StartHeart() {
	Devices.heartSend(us)
}

/**
 * err 方法
 */
func errorHandler(err error) {
	if nil != err {
		logger.Fatalf("error :%s", err.Error())
	}
}
