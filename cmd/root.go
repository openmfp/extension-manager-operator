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

	appConfig config.Config
	v         *viper.Viper
	log       *logger.Logger
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
	v, err = openmfpconfig.NewConfigFor(&appConfig)
	if err != nil {
		setupLog.Error(err, "Failed to create config")
		os.Exit(1)
	}

}

func initConfig() {
	// Parse environment variables into the Config struct
	if err := v.Unmarshal(&appConfig); err != nil {
		setupLog.Error(err, "Unable to decode into struct")
		os.Exit(1)
	}

	v.SetDefault("is-local", false)
	v.SetDefault("server-port", "8088")
	v.SetDefault("subroutines-content-configuration-enabled", true)
}

func initLog() { // coverage-ignore
	logcfg := logger.DefaultConfig()
	logcfg.Level = appConfig.Log.Level
	logcfg.NoJSON = appConfig.Log.NoJson

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
