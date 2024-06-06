package events

import (
	"fmt"
	"sync"
	"time"
)

const (
	ringBufferCapacity = 16
	timeWindowLength   = 60 // second
)

// FuncRequestFrequencyManager 管理所有函数的最近请求记录，并能够计算最近的请求频率
// 其开放的接口AddOneRequest和GetAllRecentRequestFrequency都是**线程安全**的
type FuncRequestFrequencyManager struct {
	/* fun name -> request time log */
	funcRequestTimeLog map[string]*timeRingBuffer
	mutex              sync.Mutex
}

func NewFuncRequestFrequencyManager() *FuncRequestFrequencyManager {
	fmt.Printf("New FuncRequestFrequencyManager\n")
	return &FuncRequestFrequencyManager{}
}

func (frfm *FuncRequestFrequencyManager) Init() {
	fmt.Printf("Init FuncRequestFrequencyManager\n")
	frfm.funcRequestTimeLog = make(map[string]*timeRingBuffer)
	frfm.mutex = sync.Mutex{}

}

// 这个线程安全的保证，可以确保先加入ring buffer的记录时间戳更靠前
// 也可以说，这样的ring buffer中的log遵守时间序
func (frfm *FuncRequestFrequencyManager) AddOneRequest(funcname string) {
	frfm.mutex.Lock()
	if va, exist := frfm.funcRequestTimeLog[funcname]; exist {
		// 如果此函数已经分配ring buffer, 则直接add
		va.add(time.Now())
	} else {
		// 否则，先分配ring buffer,后add
		frfm.funcRequestTimeLog[funcname] = newAndInitTimeRingBuffer()
		frfm.funcRequestTimeLog[funcname].add(time.Now())
	}
	frfm.mutex.Unlock()
}

// GetAllRecentRequestFrequency 返回所有记录函数的的最近每分钟请求数(timeWindowLength定下了“最近”的标准)
//
//	@receiver frfm
func (frfm *FuncRequestFrequencyManager) GetAllRecentRequestFrequency() map[string]float64 {
	result := map[string]float64{}
	frfm.mutex.Lock()
	get_time := time.Now()
	for func_name, func_request_log := range frfm.funcRequestTimeLog {
		result[func_name] = func_request_log.getRecentRequestFrequency(get_time)
	}
	frfm.mutex.Unlock()
	return result
}

// timeRingBuffer 以ring buffer方式记录函数的最近请求记录，并可以zero-copy方式计算请求频率
// timeRingBuffer**不是**线程安全的
type timeRingBuffer struct {
	buffer      [ringBufferCapacity]time.Time
	current_ptr int
	len         int
}

func newAndInitTimeRingBuffer() *timeRingBuffer {
	return &timeRingBuffer{
		buffer: [ringBufferCapacity]time.Time{},
		/* current_ptr按照正向顺序（index递增），指向第一个**有效的**time所在的index */
		current_ptr: 0,
		len:         0,
	}
}

// add时，ring buffer ptr后退一步，如果退到-1,则更改为（ringBufferCapacity-1）
func (trb *timeRingBuffer) add(newt time.Time) {
	trb.len += 1
	if trb.len > ringBufferCapacity {
		trb.len = ringBufferCapacity
	}

	trb.current_ptr -= 1
	if trb.current_ptr == -1 {
		trb.current_ptr += ringBufferCapacity
	}
	trb.buffer[trb.current_ptr] = newt
}

func (trb *timeRingBuffer) getRecentRequestFrequency(get_time time.Time) float64 {
	if trb.len == 0 {
		return 0
	}

	max_len := trb.len
	left_len := 0
	time_interval := float64(0)
	copy_ptr := trb.current_ptr

	// 计算time window内的数据，在遇到第一个不在time window内的log时停止
	for i := 0; i < max_len; i++ {
		time_interval = get_time.Sub(trb.buffer[copy_ptr]).Seconds()
		if time_interval <= float64(timeWindowLength) {
			left_len += 1
			copy_ptr = (copy_ptr + 1) % ringBufferCapacity
		} else {
			break
		}
	}
	// 丢弃不在time window中的数据
	trb.len = left_len
	// 计算并返回结果
	result := float64(left_len) / (time_interval / float64(60))
	return result
}

// NOTE: 下面的函数对一个ring buffer是有帮助的，但是目前暂时用不到
// func (trb *timeRingBuffer) length() int {
// 	return trb.len
// }

// func (trb *timeRingBuffer) deleteToLeft(leftlen int) {
// 	if(leftlen>=0&&leftlen<=ringBufferCapacity){
// 		trb.len=leftlen
// 	}
// }

// func (trb *timeRingBuffer) getAllTime() []time.Time {
// 	if(trb.len==ringBufferCapacity){
// 		result:=trb.buffer
// 		return result[:]
// 	}

// 	result:=[]time.Time{}
// 	ptr_copy:=trb.current_ptr
// 	for i:=0;i<trb.len;i++{
// 		result=append(result, trb.buffer[ptr_copy])
// 		ptr_copy+=1
// 		if(ptr_copy==ringBufferCapacity){
// 			ptr_copy-=ringBufferCapacity
// 		}
// 	}
// 	return result
// }
