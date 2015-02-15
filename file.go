package metric

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	ALLOCATE_FD_IDX = iota
	FREE_FD_IDX
	MAX_FD_IDX
	FD_FIELD_LEN
)

type MetricFileDescription struct {
	AllocateFdNum uint64
	FreeFdNum     uint64
	MaxFdNum      uint64
}

func GetFdInfo() (fdInfo MetricFileDescription, _ error) {
	file, err := os.Open("/proc/sys/fs/file-nr")
	if err != nil {
		return fdInfo, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != FD_FIELD_LEN {
			return fdInfo, fmt.Errorf("/proc/sys/fs/file-nr invalid")
		}
		var err error
		fdInfo.AllocateFdNum, err = strconv.ParseUint(fields[ALLOCATE_FD_IDX], 10, 64)
		if err != nil {
			return fdInfo, err
		}
		fdInfo.FreeFdNum, err = strconv.ParseUint(fields[FREE_FD_IDX], 10, 64)
		if err != nil {
			return fdInfo, err
		}
		fdInfo.MaxFdNum, err = strconv.ParseUint(fields[MAX_FD_IDX], 10, 64)
		if err != nil {
			return fdInfo, err
		}

		return fdInfo, nil
	} else {
		return fdInfo, fmt.Errorf("/proc/sys/fs/file-nr is empty")
	}
}
