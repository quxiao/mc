package mc

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	CPU_NAME_IDX = iota
	USER_IDX
	NICE_IDX
	SYSTEM_IDX
	IDLE_IDX
	IO_WAIT_IDX
	IRQ_IDX
	SOFT_IRQ_IDX
	STEAL_IDX // since Linux 2.6.11
	GUEST_IDX // since Linux 2.6.24
	// GUEST_NICE_IDX // since Linux 2.6.33
	CPU_FIELD_LEN
)

type MetricCpu struct {
	User    uint64
	System  uint64
	Nice    uint64
	Idle    uint64
	IoWait  uint64
	Irq     uint64
	SoftIrq uint64
	Steal   uint64
	Guest   uint64
}

func GetCpuInfo() (cpu MetricCpu, _ error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return cpu, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < CPU_FIELD_LEN {
			return cpu, fmt.Errorf("/proc/stat format is invalid")
		}
		var err error
		cpu.User, err = strconv.ParseUint(fields[USER_IDX], 10, 64)
		if err != nil {
			return
		}
		cpu.Nice, err = strconv.ParseUint(fields[NICE_IDX], 10, 64)
		if err != nil {
			return
		}
		cpu.System, err = strconv.ParseUint(fields[SYSTEM_IDX], 10, 64)
		if err != nil {
			return
		}
		cpu.Idle, err = strconv.ParseUint(fields[IDLE_IDX], 10, 64)
		if err != nil {
			return
		}
		cpu.IoWait, err = strconv.ParseUint(fields[IO_WAIT_IDX], 10, 64)
		if err != nil {
			return
		}
		cpu.Irq, err = strconv.ParseUint(fields[IRQ_IDX], 10, 64)
		if err != nil {
			return
		}
		cpu.SoftIrq, err = strconv.ParseUint(fields[SOFT_IRQ_IDX], 10, 64)
		if err != nil {
			return
		}
		cpu.Steal, err = strconv.ParseUint(fields[STEAL_IDX], 10, 64)
		if err != nil {
			return
		}
		cpu.Guest, err = strconv.ParseUint(fields[GUEST_IDX], 10, 64)
		if err != nil {
			return
		}
	} else {
		return cpu, fmt.Errorf("no content in /proc/stat")
	}

	return cpu, nil
}
