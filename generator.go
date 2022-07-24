package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

type VminsertData struct {
	AccountId     uint32
	ProjectId     uint32
	PrometheusUrl string
	InfluxDbUrl   string
}

type VmselectData struct {
	AccountId     uint32
	ProjectId     uint32
	PrometheusUrl string
	InfluxDbUrl   string
}

func (tctl *tctl) TenantGen() (vminsertUrls []VminsertData, vmselectUrls []VmselectData, err error) {

	// Generate seed value based on current time so that we get different values for accountId and projectId
	rand.Seed(time.Now().UnixNano())

	// Get values from config file
	numOfAccountIds := tctl.config.Int("app_num_account_ids")
	numOfProjectIds := tctl.config.Int("app_num_project_ids")

	vminsertUrlScheme := tctl.config.String("vminsert_url_scheme")
	vminsertHost := tctl.config.String("vminsert_host")
	vminsertPort := tctl.config.String("vminsert_port")

	vmselectUrlScheme := tctl.config.String("vmselect_url_scheme")
	vmselectHost := tctl.config.String("vmselect_host")
	vmselectPort := tctl.config.String("vmselect_port")

	i := 0

	if numOfAccountIds == 0 {
		return nil, nil, errors.New("account id must be one or more than one")
	}

	for i < numOfAccountIds {

		j := 0

		accountId := rand.Uint32()

		// If ProjectID is not set, keep default as 0. More Info: https://docs.victoriametrics.com/Cluster-VictoriaMetrics.html#multitenancy
		if numOfProjectIds == 0 {

			projectId := uint32(0)

			vminsertUrl := fmt.Sprintf("%s://%s:%s/insert/%d:%d/prometheus/", vminsertUrlScheme, vminsertHost, vminsertPort, accountId, projectId)
			vmselectUrl := fmt.Sprintf("%s://%s:%s/select/%d:%d/prometheus", vmselectUrlScheme, vmselectHost, vmselectPort, accountId, projectId)

			vminsertInfluxUrl := fmt.Sprintf("%s://0.0.0.0:%s/insert/%d:%d/influx/write", vminsertUrlScheme, vminsertPort, accountId, projectId)

			tctl.logger.Debug("Generating vminsert URL for ingestion", zap.String("vminsertUrl", vminsertUrl))
			tctl.logger.Debug("Generating vmselect URL for querying", zap.String("vmselectUrl", vmselectUrl))

			vminsertUrls = append(vminsertUrls, VminsertData{
				accountId,
				projectId,
				vminsertUrl,
				vminsertInfluxUrl,
			})

			vmselectUrls = append(vmselectUrls, VmselectData{
				accountId,
				projectId,
				vmselectUrl,
				"",
			})

		} else {

			for j < numOfProjectIds {

				projectId := rand.Uint32()

				vminsertUrl := fmt.Sprintf("%s://%s:%s/insert/%d:%d/prometheus/", vminsertUrlScheme, vminsertHost, vminsertPort, accountId, projectId)
				vmselectUrl := fmt.Sprintf("%s://%s:%s/select/%d:%d/prometheus", vmselectUrlScheme, vmselectHost, vmselectPort, accountId, projectId)

				vminsertInfluxUrl := fmt.Sprintf("%s://0.0.0.0:%s/insert/%d:%d/influx/write", vminsertUrlScheme, vminsertPort, accountId, projectId)

				tctl.logger.Debug("Generating vminsert URL for ingestion", zap.String("vminsertUrl", vminsertUrl))
				tctl.logger.Debug("Generating vmselect URL for querying", zap.String("vmselectUrl", vmselectUrl))

				vminsertUrls = append(vminsertUrls, VminsertData{
					accountId,
					projectId,
					vminsertUrl,
					vminsertInfluxUrl,
				})

				vmselectUrls = append(vmselectUrls, VmselectData{
					accountId,
					projectId,
					vmselectUrl,
					"",
				})

				j += 1
			}
		}
		i += 1
	}
	return vminsertUrls, vmselectUrls, nil
}
