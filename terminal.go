package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func (tctl *tctl) RenderToTerminal(vminsertUrls []VminsertData, vmselectUrls []VmselectData) {

	tctl.logger.Debug("Generating terminal output")

	vminsertTenant := [][]string{}
	vmselectTenant := [][]string{}

	for _, data := range vminsertUrls {
		vminsertTenant = append(vminsertTenant, []string{"vminsert", fmt.Sprintf("%d", data.AccountId), fmt.Sprintf("%d", data.ProjectId), data.PrometheusUrl, data.InfluxDbUrl})
	}

	for _, data := range vmselectUrls {
		vmselectTenant = append(vmselectTenant, []string{"vmselect", fmt.Sprintf("%d", data.AccountId), fmt.Sprintf("%d", data.ProjectId), data.PrometheusUrl, data.InfluxDbUrl})
	}

	bulkData := append(vminsertTenant, vmselectTenant...)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Component", "accountID", "projectID", "Prometheus URL", "InfluxDB URL"})
	table.SetAutoMergeCellsByColumnIndex([]int{0, 1})
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(true)
	table.SetNoWhiteSpace(false)
	table.AppendBulk(bulkData)
	table.Render()

	tctl.logger.Debug("Generated terminal output")
}
