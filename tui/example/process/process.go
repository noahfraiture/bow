package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Process struct {
	PID    int
	Name   string
	CPU    float64
	Memory float64
	Status string
}

func (p Process) String() string {
	name := p.Name
	if len(name) > 20 {
		name = name[:17] + "..."
	}
	return fmt.Sprintf("%-8d %-20s %-6.1f %-6.1f %s", p.PID, name, p.CPU, p.Memory, p.Status)
}

func FetchProcesses() ([]Process, error) {
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var processes []Process
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	first := true
	for scanner.Scan() {
		line := scanner.Text()
		if first {
			first = false
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue
		}
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		cpu, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			continue
		}
		mem, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			continue
		}
		name := fields[10]
		if len(fields) > 11 {
			name = strings.Join(fields[10:], " ")
		}
		processes = append(processes, Process{
			PID:    pid,
			Name:   name,
			CPU:    cpu,
			Memory: mem,
			Status: fields[7],
		})
	}
	return processes, nil
}
