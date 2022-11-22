package main

import (
	"fmt"
	"os"

	"encoding/csv"

	"go.uber.org/zap"
)

func (tctl *tctl) RenderToCsv(vminsertUrls []VminsertData, vmselectUrls []VmselectData) {

	tctl.logger.Debug("Generating CSV data")

	vminsertTenant := [][]string{}
	vmselectTenant := [][]string{}
	headers := [][]string{}

	for _, data := range vminsertUrls {
		vminsertTenant = append(vminsertTenant, []string{"vminsert", fmt.Sprintf("%d", data.AccountId), fmt.Sprintf("%d", data.ProjectId), data.PrometheusUrl, data.InfluxDbUrl})
	}

	for _, data := range vmselectUrls {
		vmselectTenant = append(vmselectTenant, []string{"vmselect", fmt.Sprintf("%d", data.AccountId), fmt.Sprintf("%d", data.ProjectId), data.PrometheusUrl, data.InfluxDbUrl})
	}

	vmData := append(vminsertTenant, vmselectTenant...)
	headers = append(headers, []string{"COMPONENT", "ACCOUNTID", "PROJECTID", "PROMETHEUS URL", "INFLUXDB URL"})
	bulkData := append(headers, vmData...)

	tctl.logger.Info("Creating CSV file")
	csvFile, err := os.Create("vmtctl.csv")
	if err != nil {
		tctl.logger.Error("Failed to create CSV file: %s", zap.Error(err))
	}

	csvWriter := csv.NewWriter(csvFile)

	for _, data := range bulkData {
		_ = csvWriter.Write(data)
	}

	csvWriter.Flush()
	csvFile.Close()

	tctl.logger.Debug("Generated CSV data")
}
