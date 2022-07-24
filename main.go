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

		ds := vmselectUrls
		dc := vminsertUrls

		// Generate grafana datasource file by populating vmselect URLs
		tctl.logger.Info("Generating grafana datasource file")

		t1, err := template.ParseFiles("templates/datasource.tmpl")
		if err != nil {
			tctl.logger.Error("Failed to parse datasource template", zap.Error(err))
			os.Exit(1)
		}

		file, _ := os.Create("provisioning/datasources/datasource.yaml")
		defer file.Close()

		err = t1.Execute(file, ds)
		if err != nil {
			tctl.logger.Error("Failed to write to datasource file", zap.Error(err))
			os.Exit(1)
		}

		// Generate docker compose file by populating vminsert URLs
		tctl.logger.Info("Generating docker compose file")

		t2, err := template.ParseFiles("templates/docker-compose.tmpl")
		if err != nil {
			tctl.logger.Error("Failed to parse template", zap.Error(err))
			os.Exit(1)
		}

		dockerCompose, _ := os.Create("docker-compose.yml")
		defer dockerCompose.Close()

		err = t2.Execute(dockerCompose, dc)
		if err != nil {
			tctl.logger.Error("Failed to write docker-compose file", zap.Error(err))
			os.Exit(1)
		}
	}

}
