package main

import (
	"context"

	altsrc "github.com/urfave/cli-altsrc/v3"
	"github.com/urfave/cli-altsrc/v3/toml"
	"github.com/urfave/cli-altsrc/v3/yaml"
	"github.com/urfave/cli/v3"
)

var (
	configFile   string
	configSource = altsrc.NewStringPtrSourcer(&configFile)

	source        string
	captureDevice string
	output        string
	destination   string
	loggingLevel  string

	processors    []string
	pipeline      []string
	configuration []string
	config        []string
	processorsDir string
	modelsDir     string

	downloadProcessors bool
	downloadModels     bool

	mcpEnabled bool
	mcpPort    string

	runFlags = []cli.Flag{
		&cli.StringFlag{Name: "file",
			Aliases:     []string{"f"},
			Usage:       "TOML or YAML file with configuration",
			Destination: &configFile,
		},
		&cli.StringFlag{Name: "source",
			Aliases:     []string{"s"},
			Value:       "0",
			Usage:       "video capture source to use. webcam id, file name, stream, testpattern (0 is the default webcam on most systems)",
			Sources:     cli.NewValueSourceChain(toml.TOML("main.source", configSource), yaml.YAML("main.source", configSource)),
			Destination: &source,
			Action: func(ctm context.Context, cmd *cli.Command, val string) error {
				if val == "testpattern" {
					// special case for testpattern, uses gstreamer to generate a test pattern
					source = "videotestsrc ! videoconvert ! appsink"
					captureDevice = "gstreamer"
				}
				return nil
			},
		},
		&cli.StringFlag{Name: "capture",
			Value:       "auto",
			Usage:       "video capture type to use (auto, ffmpeg, gstreamer, webcam)",
			Sources:     cli.NewValueSourceChain(toml.TOML("main.capture", configSource), yaml.YAML("main.capture", configSource)),
			Destination: &captureDevice,
		},
		&cli.StringFlag{Name: "output",
			Aliases:     []string{"o"},
			Value:       "mjpeg",
			Usage:       "output type (mjpeg, file)",
			Sources:     cli.NewValueSourceChain(toml.TOML("main.output", configSource), yaml.YAML("main.output", configSource)),
			Destination: &output,
		},
		&cli.StringFlag{Name: "destination",
			Aliases:     []string{"d"},
			Usage:       "output destination (port, file path)",
			Sources:     cli.NewValueSourceChain(toml.TOML("main.destination", configSource), yaml.YAML("main.destination", configSource)),
			Destination: &destination,
		},
		&cli.StringFlag{Name: "logging",
			Usage:       "logging level to use (error, warn, info, debug)",
			Value:       "warn",
			Sources:     cli.NewValueSourceChain(toml.TOML("main.logging", configSource), yaml.YAML("main.logging", configSource)),
			Destination: &loggingLevel,
		},
		&cli.StringSliceFlag{
			Name:        "processor",
			Aliases:     []string{"p"},
			Usage:       "wasm module to use for processing frames. Format: -processor /path/processor1.wasm -processor /path2/processor2.wasm",
			Destination: &processors,
		},
		&cli.StringSliceFlag{
			Name:        "pipeline",
			Sources:     cli.NewValueSourceChain(toml.TOML("processing.pipeline", configSource), yaml.YAML("processing.pipeline", configSource)),
			Destination: &pipeline,
			Hidden:      true,
		},
		&cli.StringFlag{Name: "processors-dir",
			Usage:       "directory for processor loading (default to $home/processors)",
			Sources:     cli.NewValueSourceChain(toml.TOML("processing.directory", configSource), yaml.YAML("processing.directory", configSource), cli.EnvVar("WASMVISION_PROCESSORS_DIR")),
			Destination: &processorsDir,
		},
		&cli.BoolFlag{Name: "processor-download",
			Value:       true,
			Usage:       "automatically download known processors (default: true)",
			Sources:     cli.NewValueSourceChain(toml.TOML("processing.download", configSource), yaml.YAML("processing.download", configSource)),
			Destination: &downloadProcessors,
		},
		&cli.StringSliceFlag{
			Name:        "config",
			Aliases:     []string{"c"},
			Usage:       "configuration for processors. Format: -config key1=val1 -config key2=val2",
			Sources:     cli.NewValueSourceChain(toml.TOML("config", configSource), yaml.YAML("config", configSource)),
			Destination: &config,
		},
		&cli.StringSliceFlag{
			Name:        "configuration",
			Sources:     cli.NewValueSourceChain(toml.TOML("processing.configuration", configSource), yaml.YAML("processing.configuration", configSource)),
			Destination: &configuration,
			Hidden:      true,
		},
		&cli.StringFlag{Name: "models-dir",
			Usage:       "directory for model loading (default to $home/models)",
			Sources:     cli.NewValueSourceChain(toml.TOML("models.directory", configSource), yaml.YAML("models.directory", configSource), cli.EnvVar("WASMVISION_MODELS_DIR")),
			Destination: &modelsDir,
		},
		&cli.BoolFlag{Name: "models-download",
			Aliases:     []string{"download"},
			Value:       true,
			Usage:       "automatically download known models (default: true)",
			Sources:     cli.NewValueSourceChain(toml.TOML("models.downloads", configSource), yaml.YAML("models.downloads", configSource)),
			Destination: &downloadModels,
		},
		&cli.BoolFlag{Name: "mcp-server",
			Value:       false,
			Usage:       "enable MCP server",
			Sources:     cli.NewValueSourceChain(toml.TOML("server.mcp-enable", configSource), yaml.YAML("server.mcp-enable", configSource)),
			Destination: &mcpEnabled,
		},
		&cli.StringFlag{Name: "mcp-port",
			Value:       ":5001",
			Usage:       "port for MCP server (default: :5001)",
			Sources:     cli.NewValueSourceChain(toml.TOML("server.mcp-port", configSource), yaml.YAML("server.mcp-port", configSource)),
			Destination: &mcpPort,
		},
	}

	downloadFlags = []cli.Flag{
		&cli.StringFlag{Name: "models-dir",
			Aliases: []string{"models"},
			Usage:   "directory for model downloading (default to $home/models)",
			Sources: cli.EnvVars("WASMVISION_MODELS_DIR"),
		},
		&cli.StringFlag{Name: "processors-dir",
			Aliases: []string{"processors"},
			Usage:   "directory for processor downloading (default to $home/processors)",
			Sources: cli.EnvVars("WASMVISION_PROCESSORS_DIR"),
		},
	}
)
