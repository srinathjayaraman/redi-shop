package cmd

import (
	"fmt"
	"os"

	"github.com/martijnjanssen/redi-shop/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string
	service string

	rootCmd = &cobra.Command{
		Use:   "redi",
		Short: "redi is a microservice shop implementation",
		Long: `A fast and resilient microservice implementation with different
                backends.`,
		Run: func(cmd *cobra.Command, args []string) {
			server.Start()
		},
	}
)

// Initialize commands
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./redi.yaml)")
	rootCmd.Flags().StringVarP(&service, "service", "s", "", "Service to start (user, stock, order, payment)")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath("./")
		viper.SetConfigName("redi")
	}

	// Bind service flag to env values
	err := viper.BindPFlag("service", rootCmd.Flags().Lookup("service"))
	if err != nil {
		logrus.WithError(err).Fatal("unable to bind service flag to config value")
	}

	// Default config values
	viper.SetDefault("postgres.url", "localhost")
	viper.SetDefault("postgres.port", "5432")
	viper.SetDefault("postgres.username", "postgres")
	viper.SetDefault("postgres.password", "postgres")
	viper.SetDefault("postgres.database", "redi")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.WithField("file", viper.ConfigFileUsed()).Info("Loaded config")
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
