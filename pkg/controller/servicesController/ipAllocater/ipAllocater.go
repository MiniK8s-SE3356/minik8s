package ipAllocater

import (
	"fmt"
	"math"
)

type IpAllocater struct {
	Subnet     []int16
	MaskBitNum int8
	IPBitMap   []int8
}

func (ia *IpAllocater) Init() {
	// 获取子网和子网掩码，转换为本地数据
	ia.Subnet = []int16{10, 100, 0, 0}
	ia.MaskBitNum = 16

	// 创建位图
	length := int(math.Pow(2, float64(32-ia.MaskBitNum)))
	ia.IPBitMap = make([]int8, length)

	// 位图全部赋0
	for i := 0; i < length; i++ {
		ia.IPBitMap[i] = 0
	}
}

func (ia *IpAllocater) AllocateIP() string {
	length := len(ia.IPBitMap)
	for i := 0; i < length; i++ {
		if ia.IPBitMap[i] == 0 {
			ia.IPBitMap[i] = 1

			subnet := make([]int16, 4)
			copy(subnet, ia.Subnet)

			suboff := 3
			bitmapoff := i
			for suboff >= 0 {
				subnet[suboff] = int16(bitmapoff % 256)
				bitmapoff = bitmapoff / 256
				if bitmapoff == 0 {
					break
				}
				suboff--
			}
			ip_str := fmt.Sprintf("%d.%d.%d.%d", subnet[0], subnet[1], subnet[2], subnet[3])
			fmt.Printf("Allocate New Virtual IP %s\n", ip_str)
			return ip_str
		}
	}
	return "error"
}

func (ia *IpAllocater) DeallocateIP(ip string) error {
	subnet := make([]int16, 4)
	_, err := fmt.Sscanf(ip, "%d.%d.%d.%d", &(subnet[0]), &(subnet[1]), &(subnet[2]), &(subnet[3]))
	if err != nil {
		return err
	}

	// 计算子网掩码后的剩余子网段
	maskbitnum := ia.MaskBitNum
	suboff := 0
	for suboff < 4 {
		if maskbitnum >= 8 {
			subnet[suboff] = 0

			maskbitnum -= 8
			if (maskbitnum) <= 0 {
				break
			}
		} else {
			leftbitnum := 8 - maskbitnum
			mask := int16(1<<leftbitnum) - 1

			subnet[suboff] = subnet[suboff] & mask
			break
		}
		suboff++
	}

	// 计算bitmap偏移量
	var offset int = 0
	for _, value := range subnet {
		offset = (offset << 8) + int(value)
	}

	ia.IPBitMap[offset] = 0

	return nil
}
