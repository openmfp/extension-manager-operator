package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	corev1alpha1 "github.com/openmfp/extension-manager-operator/api/v1alpha1"
	"github.com/openmfp/extension-manager-operator/internal/config"
	openmfpconfig "github.com/openmfp/golang-commons/config"
	"github.com/openmfp/golang-commons/logger"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")

	operatorCfg config.OperatorConfig
	serverCfg   config.ServerConfig
	defaultCfg  openmfpconfig.CommonServiceConfig
	v           *viper.Viper
	log         *logger.Logger
)

var rootCmd = &cobra.Command{
	Use:   "extension-manager-operator",
	Short: "operator to reconcile ContentConfiguration",
}

func init() { // coverage-ignore
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(corev1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme

	rootCmd.AddCommand(operatorCmd)
	rootCmd.AddCommand(serverCmd)

	cobra.OnInitialize(initConfig, initLog)

	var err error
	v, err = openmfpconfig.NewDefaultConfig(rootCmd)
	if err != nil {
		setupLog.Error(err, "Failed to create config")
		os.Exit(1)
	}

	err = openmfpconfig.BindConfigToFlags(v, operatorCmd, &operatorCfg)
	if err != nil {
		setupLog.Error(err, "Failed to bind config to flags")
		os.Exit(1)
	}
	err = openmfpconfig.BindConfigToFlags(v, serverCmd, &serverCfg)
	if err != nil {
		setupLog.Error(err, "Failed to bind config to flags")
		os.Exit(1)
	}

}

func initConfig() {

	v.SetDefault("is-local", false)
	v.SetDefault("server-port", "8088")
	v.SetDefault("subroutines-content-configuration-enabled", true)

	// Parse environment variables into the Config struct
	if err := v.Unmarshal(&defaultCfg); err != nil {
		setupLog.Error(err, "Unable to decode into struct")
		os.Exit(1)
	}

	// Parse environment variables into the Config struct
	if err := v.Unmarshal(&operatorCfg); err != nil {
		setupLog.Error(err, "Unable to decode into struct")
		os.Exit(1)
	}

	// Parse environment variables into the Config struct
	if err := v.Unmarshal(&serverCfg); err != nil {
		setupLog.Error(err, "Unable to decode into struct")
		os.Exit(1)
	}
}

func initLog() { // coverage-ignore
	logcfg := logger.DefaultConfig()
	logcfg.Level = defaultCfg.Log.Level
	logcfg.NoJSON = defaultCfg.Log.NoJson

	var err error
	log, err = logger.New(logcfg)
	if err != nil {
		setupLog.Error(err, "unable to create logger")
		os.Exit(1)
	}
}

func Execute() { // coverage-ignore
	cobra.CheckErr(rootCmd.Execute())
}
