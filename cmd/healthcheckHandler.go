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

	mem := make(map[string]interface{})
	cpu := make(map[string]interface{})
	mem["total"] = vm.Total
	mem["available"] = vm.Available
	mem["used_percent"] = vm.UsedPercent
	cpu["physical_core"] = physicalCnt
	cpu["logical_core"] = logicalCnt
	cpu["total_percent"] = totalPercent
	cpu["per_percent"] = perPercents

	//TBD: change a package to handle JSON for higher perfermance
	err := writeJSON(w, http.StatusOK, envelope{"mem": mem, "cpu": cpu}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
