package idgen

import (
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"net"
)

// GetMachineID 根据网络信息生成唯一的机器ID
func GetMachineID() (uint16, error) {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return 0, err
	}

	// 寻找第一个有效的MAC地址
	var mac net.HardwareAddr
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) > 0 && (iface.Flags&net.FlagUp) != 0 && (iface.Flags&net.FlagLoopback) == 0 {
			mac = iface.HardwareAddr
			break
		}
	}

	if mac == nil {
		return 0, errors.New("no valid network interface found")
	}

	// 使用SHA1哈希MAC地址
	hash := sha1.New()
	hash.Write([]byte(mac.String()))
	hashBytes := hash.Sum(nil)

	// 取哈希值的前2个字节作为machineID
	machineID := binary.BigEndian.Uint16(hashBytes[:2])

	return machineID, nil
}
