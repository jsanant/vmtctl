package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/knadh/koanf"
	"go.uber.org/zap"
)

type tctl struct {
	logger *zap.Logger
	config *koanf.Koanf
}

func main() {

	k, err := initConfig()
	if err != nil {
		fmt.Println("Error when initializing config: ", err)
		os.Exit(1)
	}

	var logger = initLogger()

	tctl := tctl{
		logger: logger,
		config: k,
	}

	logger.Info("Generating endpoints")

	vminsertUrls, vmselectUrls, err := tctl.TenantGen()
	if err != nil {
		tctl.logger.Error("Error when generating endpoints", zap.Error(err))
		os.Exit(1)
	}

	tctl.logger.Debug("Showing vm tenant data on terminal")

	tctl.RenderToTerminal(vminsertUrls, vmselectUrls)

	// Generate only endpoints and exit, no need to generate grafana datasource and docker-compose files
	if !tctl.config.Bool("gen-endpoints") {

		// vmselect endpoints for Grafana datasource file
		ds := vmselectUrls

		// vminsert endpoints for vminsert container in docker-compose file
		dc := vminsertUrls

		// Generate grafana datasource file
		tctl.logger.Info("Generating grafana datasource file")

		datasourceTmpl, err := template.ParseFiles("templates/datasource.tmpl")
		if err != nil {
			tctl.logger.Error("Failed to parse datasource template", zap.Error(err))
			os.Exit(1)
		}

		datasourceFile, _ := os.Create("provisioning/datasources/datasource.yaml")
		defer datasourceFile.Close()

		err = datasourceTmpl.Execute(datasourceFile, ds)
		if err != nil {
			tctl.logger.Error("Failed to write to datasource file", zap.Error(err))
			os.Exit(1)
		}

		// Generate docker compose file
		tctl.logger.Info("Generating docker compose file")

		dockercomposeTmpl, err := template.ParseFiles("templates/docker-compose.tmpl")
		if err != nil {
			tctl.logger.Error("Failed to parse docker-compose template", zap.Error(err))
			os.Exit(1)
		}

		dockerComposeFile, _ := os.Create("docker-compose.yml")
		defer dockerComposeFile.Close()

		err = dockercomposeTmpl.Execute(dockerComposeFile, dc)
		if err != nil {
			tctl.logger.Error("Failed to write to docker-compose file", zap.Error(err))
			os.Exit(1)
		}
	}

}
