package idworker

import (
	"log"
	"os"
	"sync"
	"time"
)

// 雪花算法格式：时间戳（42） + 数据中心（5）+ 机器（5）+ 序列号（12）

const (
	dataCenterBits uint8 = 0  // 数据中心的ID位数，5位最大可以有2^5-1=31个节点
	workerBits     uint8 = 14 // 机器的ID位数 5位最大可以有2^5-1=31个节点
	sequenceBits   uint8 = 8  // 表示每个集群下的每个节点，1毫秒内可生成的id序号的二进制位数 即每毫秒可生成 2^12-1=4096个唯一ID

	//dataCenterMax int64 = -1 ^ (-1 << dataCenterBits) // 数据中心Id的最大值，用于防止溢出
	workerMax   int64 = -1 ^ (-1 << workerBits)   // 机器ID的最大值，用于防止溢出
	sequenceMax int64 = -1 ^ (-1 << sequenceBits) // 序列号的最大值，用于防止溢出

	workerShift     = sequenceBits                     // 机器ID向左的偏移量
	dataCenterShift = workerBits + workerShift         // 数据中心ID向左的偏移量
	timeShift       = dataCenterBits + dataCenterShift // 时间戳向左的偏移量

	// 41位字节作为时间戳数值的话 大约68年就会用完
	// 雪花算法初始时间戳，一旦设置不允许修改
	epoch int64 = 1631082836000
)

type WorkerId int64
type DataCenterId int64

// SnowflakeWorker 雪花算法
type SnowflakeWorker struct {
	mu           sync.Mutex // 添加互斥锁 确保并发安全
	dataCenterId int64      // 数据中心Id
	workerId     int64      // 机器ID
	epoch        int64      // 雪花算法初始时间戳

	timestamp int64 // 记录时间戳
	sequence  int64 // 当前毫秒已经生成的id序列号(从0开始累加) 1毫秒内最多生成4096个ID
}

// NewSnowflakeWorkerForPid 使用当前进程号除最大的worker数量的余数作者workerId
func NewSnowflakeWorkerForPid() *SnowflakeWorker {
	wid := int64(os.Getpid()) % MaxWorkId()
	return NewSnowflakeWorker(WorkerId(wid))
}

// NewSnowflakeWorker 实例化一个雪花算法
func NewSnowflakeWorker(workId WorkerId) *SnowflakeWorker {
	//dataCenterId := conf.DataCenterId
	//if dataCenterId < 0 || dataCenterId > dataCenterMax {
	//	log.Fatalf("dataCenterId exceeded maximum value %v", dataCenterMax)
	//}

	workerId := int64(workId)
	if workerId < 0 || workerId > workerMax {
		log.Fatalf("workerId exceeded maximum value %v", workerMax)
	}

	return &SnowflakeWorker{
		epoch: epoch,
		//dataCenterId: dataCenterId,
		workerId:  workerId,
		timestamp: 0,
		sequence:  0,
	}
}

// NextId 生成下一个Id
func (w *SnowflakeWorker) NextId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 获取生成时的时间戳
	now := time.Now().UnixNano() / 1e6 // 纳秒转毫秒
	if w.timestamp == now {
		w.sequence++

		// 判断当前节点是否在1毫秒内已经生成sequenceMax个ID
		if w.sequence > sequenceMax {
			// 如果当前工作节点在1毫秒内生成的ID已经超过上限 需要等待1毫秒再继续生成
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 如果当前时间与工作节点上一次生成ID的时间不一致 则需要重置工作节点生成ID的序号
		w.sequence = 0
		w.timestamp = now
	}

	// 计算时间已经流失了xxx毫秒
	timeFly := now - w.epoch
	return timeFly<<timeShift | (w.dataCenterId << dataCenterShift) | (w.workerId << workerShift) | (w.sequence)
}

func MaxWorkId() int64 {
	return workerMax
}
