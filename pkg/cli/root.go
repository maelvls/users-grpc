package cli

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/lithammer/dedent"
	"github.com/maelvls/users-grpc/pkg/cli/logutil"
	"github.com/maelvls/users-grpc/schema/user"
	"github.com/mattn/go-isatty"
	"github.com/mgutz/ansi"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var cfgFile string
var verbose bool
var version Version
var cfg clientCfg

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "users-cli (list | search | create | get)",
	Short: "A nice CLI for querying users from the user-grpc microservice.",

	// https://github.com/spf13/cobra#prerun-and-postrun-hooks
	// This hook is also executed when subcommands are run.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cfg = clientCfg{
			address:    viper.GetString("address"),
			cacert:     viper.GetString("cacert"),
			cleartext:  viper.GetBool("cleartext"),
			servername: viper.GetString("servername"),
		}
		logutil.Debugf("config: %v", cfg)
		switch viper.GetString("color") {
		case "auto":
			ansi.DisableColors(!isatty.IsTerminal(os.Stdout.Fd()))
		case "always":
			ansi.DisableColors(false)
		case "never":
			ansi.DisableColors(true)
		default:
			logrus.Errorf("%s is not a valid value for --color; must be either 'auto', 'always' or 'never'", viper.GetString("color"))
			os.Exit(1)
		}
	},
	Long: dedent.Dedent(`
	Note: I could not add an '--insecure' flag that would disable TLS verification
	due to the lack of support in the go-grpc lib.
	`),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(v Version) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	logrus.SetFormatter(&logrus.TextFormatter{})
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.users-cli.yml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().String("address", ":8000", "The host:port to bind to. Alternatively, you can set ADDRESS or add 'address: localhost:8000' in $HOME/.users-cli.yml")
	rootCmd.PersistentFlags().String("color", "auto", "Supported are 'auto', 'always' and 'never'. In 'auto' mode, colors are enabled when stdout is a tty.")
	rootCmd.PersistentFlags().String("cacert", "", "CA certificate to verify the server's TLS certificate against.")
	rootCmd.PersistentFlags().Bool("cleartext", false, "Use HTTP/2 in cleartext mode (h2c).")
	rootCmd.PersistentFlags().String("servername", "", "Override server name when validating TLS certificate. Useful when testing locally.")
	err := viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("address"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("color", rootCmd.PersistentFlags().Lookup("color"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("cacert", rootCmd.PersistentFlags().Lookup("cacert"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("cleartext", rootCmd.PersistentFlags().Lookup("cleartext"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("servername", rootCmd.PersistentFlags().Lookup("servername"))
	if err != nil {
		panic(err)
	}
}

type clientCfg struct {
	address    string
	cacert     string
	servername string // Often used while testing.
	cleartext  bool
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if verbose {
		logutil.EnableDebug = true
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".users-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logutil.Debugf("using config file: %v", viper.ConfigFileUsed())
	}
}

func createClient(config clientCfg) (user.UserServiceClient, error) {
	var err error

	split := strings.SplitN(config.address, ":", 2)
	if len(split) != 2 {
		return nil, fmt.Errorf("address should be of the form host:port or :port")
	}
	servername := split[0]
	if servername == "" {
		servername = "127.0.0.1"
	}
	if config.servername != "" {
		servername = config.servername
	}
	creds := credentials.NewTLS(&tls.Config{ServerName: servername})
	if config.cacert != "" {
		creds, err = credentials.NewClientTLSFromFile(config.cacert, config.servername)
		if err != nil {
			return nil, fmt.Errorf("loading CA certificates: %w", err)
		}
	}

	var opts []grpc.DialOption

	if config.cleartext {
		opts = append(opts, grpc.WithInsecure())
	} else {
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	// c := credentials.CommonAuthInfo{SecurityLevel: credentials.NoSecurity}}

	cc, err := grpc.Dial(config.address, opts...)
	if err != nil {
		logutil.Errorf("%s", err)
		os.Exit(1)
	}

	client := user.NewUserServiceClient(cc)

	return client, err
}
