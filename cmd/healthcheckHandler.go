package main

import (
	"net/http"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	vm, _ := mem.VirtualMemory()

	physicalCnt, _ := cpu.Counts(false)
	logicalCnt, _ := cpu.Counts(true)

	totalPercent, _ := cpu.Percent(3*time.Second, false)
	perPercents, _ := cpu.Percent(3*time.Second, true)

	device := make(map[string]interface{})
	device["mem:total"] = vm.Total
	device["mem:available"] = vm.Available
	device["mem:used_percent"] = vm.UsedPercent
	device["cpu:physical_core"] = physicalCnt
	device["cpu:logical_core"] = logicalCnt
	device["cpu:total_percent"] = totalPercent
	device["cpu:per_percent"] = perPercents

	//TBD: change a package to handle JSON for higher perfermance
	err := writeJSON(w, http.StatusOK, device, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Server fail to process your Request", http.StatusInternalServerError)
	}
}
