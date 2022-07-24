package main

import (
	"fmt"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
)

// Inject version at build time
var version string

func initLogger() *zap.Logger {

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Print("Failed to initialise logger: ", err)
		os.Exit(1)
	}
	defer logger.Sync()

	return logger
}

func initConfig() (*koanf.Koanf, error) {

	var k = koanf.New("_")
	f := flag.NewFlagSet("vmtctl", flag.ContinueOnError)

	// Print CLI flags when -h is used
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	// Define CLI flags
	f.String("config", "sample_config.toml", "Path to the .toml config file")
	f.Bool("gen-endpoints", false, "Generate only endpoints")
	f.BoolP("version", "v", false, "Show version")

	// Parse the defined flags
	err := f.Parse(os.Args[1:])
	if err != nil {
		fmt.Print("Failed to parse args")
		os.Exit(1)
	}

	// Print version and exit
	if ok, _ := f.GetBool("version"); ok {
		fmt.Println(version)
		os.Exit(0)
	}

	// Load the config file provided in the command line
	c, _ := f.GetString("config")
	if err := k.Load(file.Provider(c), toml.Parser()); err != nil {
		return nil, err
	}

	// Overwrite from command line
	if err := k.Load(posflag.Provider(f, "_", k), nil); err != nil {
		return nil, err
	}

	return k, err
}
