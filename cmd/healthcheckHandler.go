package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	vm, _ := mem.VirtualMemory()

	memInfo := fmt.Sprintf("Total: %v, Available: %v, UsedPercent:%f%%\n", vm.Total, vm.Available, vm.UsedPercent)

	physicalCnt, _ := cpu.Counts(false)
	logicalCnt, _ := cpu.Counts(true)

	cpuInfo := fmt.Sprintf("physical count:%d logical count:%d\n", physicalCnt, logicalCnt)

	fmt.Fprintln(w, cpuInfo)

	totalPercent, _ := cpu.Percent(3*time.Second, false)
	perPercents, _ := cpu.Percent(3*time.Second, true)
	cpuUsage := fmt.Sprintf("total percent:%v per percents:%v", totalPercent, perPercents)

	fmt.Fprintln(w, cpuUsage)

	fmt.Fprintln(w, memInfo)

}
