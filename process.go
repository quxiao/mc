package mc

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	PS_PID_IDX                   = iota
	PS_COMM_IDX                  // The filename of the executable, in parentheses. This is visible whether or not the executable is swapped out.
	PS_STATE_IDX                 // The state of process
	PS_PPID_IDX                  // The PID of the parent of this process.
	PS_PGRP_IDX                  // The process group ID of the process.
	PS_SESSION_IDX               // The session ID of the process.
	PS_TTY_NR_IDX                // The controlling terminal of the process.  (The minor device number is contained in the combination of bits 31 to 20 and 7 to 0; the major device number is in bits 15 to 8.)
	PS_PGDID_IDX                 // The ID of the foreground process group of the controlling terminal of the process.
	PS_FLAGS_IDX                 // The kernel flags word of the process.
	PS_MINFLT_IDX                // The number of minor faults the process has made which have not required loading a memory page from disk.
	PS_CMINFLT_IDX               // The number of minor faults that the process's waited-for children have made.
	PS_MAJFLT_IDX                // The number of major faults the process has made which have required loading a memory page from disk.
	PS_CMAJFLT_IDX               // The number of major faults that the process's waited-for children have made.
	PS_UTIME_IDX                 // Amount of time that this process has been scheduled in user mode, measured in clock ticks (divide by sysconf(_SC_CLK_TCK)).
	PS_STIME_IDX                 // Amount of time that this process has been scheduled in kernel mode, measured in clock ticks (divide by sysconf(_SC_CLK_TCK)).
	PS_CUTIME_IDX                // Amount of time that this process's waited-for children have been scheduled in user mode, measured in clock ticks.
	PS_CSTIME_IDX                // Amount of time that this process's waited-for children have been scheduled in kernel mode, measured in clock ticks (divide by sysconf(_SC_CLK_TCK)).
	PS_PRIORITY_IDX              // (Explanation for Linux 2.6) For processes running a real-time scheduling policy (policy below; see sched_setscheduler(2)), this is the negated scheduling priority, minus one; that is, a number in the range -2 to -100, corresponding to real-time priorities 1 to 99.  For processes running under a non-real-time scheduling policy, this is the raw nice value (setpriority(2)) as represented in the kernel.  The kernel stores nice values as numbers in the range 0 (high) to 39 (low), corresponding to the user-visible nice range of -20 to 19. Before Linux 2.6, this was a scaled value based on the scheduler weighting given to this process.
	PS_NICE_IDX                  // The nice value (see setpriority(2)), a value in the range 19 (low priority) to -20 (high priority).
	PS_NUM_THREADS_IDX           // Number of threads in this process (since Linux 2.6). Before kernel 2.6, this field was hard coded to 0 as a placeholder for an earlier removed field.
	PS_ITVEALVALUE_IDX           //The time in jiffies before the next SIGALRM is sent to the process due to an interval timer.  Since kernel 2.6.17, this field is no longer maintained, and is hard coded as 0.
	PS_STARTTIME_IDX             // The time the process started after system boot.  In kernels before Linux 2.6, this value was expressed in jiffies.  Since Linux 2.6, the value is expressed in clock ticks (divide by sysconf(_SC_CLK_TCK)).
	PS_VSIZE_IDX                 // Virtual memory size in bytes.
	PS_RSS_IDX                   // Resident Set Size: number of pages the process has in real memory.  This is just the pages which count toward text, data, or stack space.  This does not include pages which have not been demand-loaded in, or which are swapped out.
	PS_RSSLIM_IDX                // Current soft limit in bytes on the rss of the process.
	PS_STARTCODE_IDX             // The address above which program text can run.
	PS_ENDCODE_IDX               // The address below which program text can run.
	PS_STARTSTACK_IDX            // The address of the start (i.e., bottom) of the stack.
	PS_KSTKESP_IDX               // The current value of ESP (stack pointer), as found in the kernel stack page for the process.
	PS_KSTKEIP_IDX               // The current EIP (instruction pointer).
	PS_SIGNAL_IDX                //Obsolete
	PS_BLOCKED_IDX               //Obsolete
	PS_SIGIGNORE_IDX             //Obsolete
	PS_SIGCATCH_IDX              //Obsolete
	PS_WCHAN_IDX                 //This is the "channel" in which the process is waiting.  It is the address of a location in the kernel where the process is sleeping.  The corresponding symbolic name can be found in /proc/[pid]/wchan.
	PS_NSWAP_IDX                 // Number of pages swapped (not maintained)
	PS_CNSWAP_IDX                // Cumulative nswap for child processes (not maintained)
	PS_EXIT_SIGNAL_IDX           // Signal to be sent to parent when we die.
	PS_PROCESSOR_IDX             // CPU number last executed on.
	PS_RT_PRIORITY_IDX           // Real-time scheduling priority, a number in the range 1 to 99 for processes scheduled under a real-time policy, or 0, for non-real-time processes
	PS_POLICY_IDX                // Scheduling policy (see sched_setscheduler(2)). Decode using the SCHED_* constants in linux/sched.h.
	PS_DELAYACCT_BLKIO_TICKS_IDX // Aggregated block I/O delays, measured in clock ticks (centiseconds).
	PS_GUEST_TIME_IDX            // Guest time of the process (time spent running a virtual CPU for a guest operating system), measured in clock ticks (divide by sysconf(_SC_CLK_TCK)).
	PS_CGUEST_TIME_IDX           // Guest time of the process's children, measured in clock ticks (divide by sysconf(_SC_CLK_TCK)).
	PS_FIELDS_LEN
	/*
		PS_START_DATA_IDX (since Linux 3.3)
		PS_END_DATA_IDX
		PS_START_BRK_IDX
		PS_ARG_START_IDX
		PS_ARG_END_IDX
		PS_ENV_START_IDX
		PS_ENV_END_IDX
		PS_EXIT_CODE_IDX
	*/

)

type ProcessStat struct {
	Utime     uint64
	Stime     uint64
	ThreadNum uint64
	MinFltNum uint64
	MajFltNum uint64
	Vsize     uint64
	Rss       uint64
	RssLimit  uint64
}

func GetProcessUtime(pid int) (uint64, error) {
	ps, err := getProcessStat(pid)
	if err != nil {
		return 0, err
	}
	return ps.Utime, nil
}

func GetProcessStime(pid int) (uint64, error) {
	ps, err := getProcessStat(pid)
	if err != nil {
		return 0, err
	}
	return ps.Stime, nil
}

func GetProcessThreadNum(pid int) (uint64, error) {
	ps, err := getProcessStat(pid)
	if err != nil {
		return 0, err
	}
	return ps.ThreadNum, nil
}

func GetProcessMinFltNum(pid int) (uint64, error) {
	ps, err := getProcessStat(pid)
	if err != nil {
		return 0, err
	}
	return ps.MinFltNum, nil
}

func GetProcessMajFltNum(pid int) (uint64, error) {
	ps, err := getProcessStat(pid)
	if err != nil {
		return 0, err
	}
	return ps.MajFltNum, nil
}

func GetProcessVsize(pid int) (uint64, error) {
	ps, err := getProcessStat(pid)
	if err != nil {
		return 0, err
	}
	return ps.Vsize, nil
}

func GetProcessRss(pid int) (uint64, error) {
	ps, err := getProcessStat(pid)
	if err != nil {
		return 0, err
	}
	return ps.Rss, nil
}

func GetProcessRssLimit(pid int) (uint64, error) {
	ps, err := getProcessStat(pid)
	if err != nil {
		return 0, err
	}
	return ps.RssLimit, nil
}

func getProcessStat(pid int) (ps ProcessStat, _ error) {
	file, err := os.Open(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return ps, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	if scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < PS_FIELDS_LEN {
			return ps, fmt.Errorf("/proc/%d/stat format is invalid", pid)
		}
		var err error
		ps.Utime, err = strconv.ParseUint(fields[PS_UTIME_IDX], 10, 64)
		if err != nil {
			return ps, err
		}
		ps.Stime, err = strconv.ParseUint(fields[PS_STIME_IDX], 10, 64)
		if err != nil {
			return ps, err
		}
		ps.ThreadNum, err = strconv.ParseUint(fields[PS_NUM_THREADS_IDX], 10, 64)
		if err != nil {
			return ps, err
		}
		ps.MinFltNum, err = strconv.ParseUint(fields[PS_MINFLT_IDX], 10, 64)
		if err != nil {
			return ps, err
		}
		ps.MajFltNum, err = strconv.ParseUint(fields[PS_MAJFLT_IDX], 10, 64)
		if err != nil {
			return ps, err
		}
		ps.Vsize, err = strconv.ParseUint(fields[PS_VSIZE_IDX], 10, 64)
		if err != nil {
			return ps, err
		}
		ps.Rss, err = strconv.ParseUint(fields[PS_RSS_IDX], 10, 64)
		if err != nil {
			return ps, err
		}
		ps.RssLimit, err = strconv.ParseUint(fields[PS_RSSLIM_IDX], 10, 64)
		if err != nil {
			return ps, err
		}
	} else {
		return ps, fmt.Errorf("/proc/%d/stat is empty", pid)
	}

	return ps, nil
}
