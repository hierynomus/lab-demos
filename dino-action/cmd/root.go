package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hierynomus/autobind"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/stackvista/demo/kubecon2024/poi/internal/config"
	"gitlab.com/stackvista/demo/kubecon2024/poi/version"
)

const (
	VerboseFlag      = "verbose"
	VerboseFlagShort = "v"
)

func RootCommand(cfg *config.Config) *cobra.Command {
	var verbosity int

	cmd := &cobra.Command{
		Use:   "poi",
		Short: "PoI - Points of Interest Service",
		Long:  "PoI - Points of Interest Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			switch verbosity {
			case 0:
				// Nothing to do
			case 1:
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			default:
				zerolog.SetGlobalLevel(zerolog.TraceLevel)
			}

			logger := log.Ctx(cmd.Context())

			vp := viper.New()
			vp.SetConfigName("config")
			vp.AddConfigPath(".")
			vp.AddConfigPath("/config")
			vp.AddConfigPath("/etc/poi")
			vp.SetConfigType("yaml")
			binder := &autobind.Autobinder{UseNesting: true, EnvPrefix: "POI", ConfigObject: cfg, Viper: vp, SetDefaults: true}

			if err := vp.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					logger.Warn().Msg("No config file found... Continuing with defaults")
					// Config file not found; ignore error if desired
				} else {
					fmt.Printf("%s", err)
					os.Exit(1)
				}
			}

			binder.Bind(cmd.Context(), cmd, []string{})

			logger.Info().Str("version", version.Version).Str("commit", version.Commit).Str("date", version.Date).Msg("PoI Version")

			return nil
		},
	}

	cmd.PersistentFlags().CountVarP(&verbosity, VerboseFlag, VerboseFlagShort, "Print verbose logging to the terminal (use multiple times to increase verbosity)")
	return cmd
}

func Execute(ctx context.Context) {
	config := &config.Config{}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cmd := RootCommand(config)
	cmd.AddCommand(NewStartCmd(config))

	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
