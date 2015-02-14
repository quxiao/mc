package metric

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

var CORE_MEM_METRIC_MAP = make(map[string]struct{})

type MetricMemory struct {
	MemTotal uint64
	RealFree uint64 // = (MemFree + Buffers + Cached)
	MemFree  uint64
	Buffers  uint64
	Cached   uint64
}

func init() {
	CORE_MEM_METRIC_MAP["MemTotal"] = struct{}{}
	CORE_MEM_METRIC_MAP["MemFree"] = struct{}{}
	CORE_MEM_METRIC_MAP["Buffers"] = struct{}{}
	CORE_MEM_METRIC_MAP["Cached"] = struct{}{}
}

func GetMemoryInfo() (mem MetricMemory, _ error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return mem, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	metricMap := make(map[string]uint64)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 || len(fields[0]) == 0 {
			continue
		}
		metric := fields[0][:len(fields[0])-1]
		if _, ok := CORE_MEM_METRIC_MAP[metric]; !ok {
			continue
		}
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		metricMap[metric] = value
		if len(metricMap) == len(CORE_MEM_METRIC_MAP) {
			break
		}
	}

	mem.MemTotal = metricMap["MemTotal"]
	mem.MemFree = metricMap["MemFree"]
	mem.Buffers = metricMap["Buffers"]
	mem.Cached = metricMap["Cached"]
	mem.RealFree = mem.MemFree + mem.Buffers + mem.Cached

	return mem, nil
}
