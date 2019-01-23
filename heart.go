/**
 * 在线状态检测包
 * 检测方式: 复合
 * 1: 被动检测 (超时验证)
 *   .1 在每次设备发送数据时, 更新设备状态
 * 2: 主要检测 (心跳机制)
 *   .2 主要向设备发送状态, 接收回复判断设备状态.
 */
package udpServer

import (
	"errors"
	"net"
	"time"
)

var Devices = devices{
	Devices: make(map[int]*device),
}

type devices struct {
	Devices map[int]*device
}

type device struct {
	Addr *net.UDPAddr

	IsOnline    bool  // 是否在线
	LostCount   int   // 临时的, 丢失次数
	OffLineTime int64 // 离线时间点
	LastTime    int64 // 上次通过验证的时间点
}

func newDevice(addr *net.UDPAddr) *device {
	return &device{
		Addr:        addr,
		IsOnline:    true,
		LostCount:   0,
		OffLineTime: 0,
		LastTime:    time.Now().Unix(),
	}
}

/**
 * 添加一个设备
 */
func (d *devices) UpdateDevice(addr *net.UDPAddr) {
	// 先更新设备, 如果更新失败, 则添加这个设备
	if err := d.updateDeviceState(addr); nil != err {
		d.Devices[addr.Port] = newDevice(addr)
	}
}

/**
 * 重置设备状态
 */
func (d *devices) updateDeviceState(addr *net.UDPAddr) (err error) {
	device := d.Devices[addr.Port]
	if nil == device {
		err = errors.New("not find device")
		return
	}
	device.IsOnline = true
	device.LastTime = time.Now().Unix()
	return
}

/**
 * 获取设备状态
 */
func (d *devices) GetLineState(port int) (result bool) {
	result = d.Devices[port].IsOnline
	return
}

/**
 * 发送心跳包.
 */
func (d *devices) heartSend(us *UdpServer) {
	for _, o := range d.Devices {
		_, err := us.WriteToUDP([]byte("heart"), o.Addr)
		if nil != err {
			o.LostCount += 1
			continue
		}
	}
}
