package cmd

import (
	"galactus/blade/routers"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Server struct {
	// 配置文件路径
	ConfigPath string
}

func init() {
	StartCmd.PersistentFlags().StringVarP(&cmd.ConfigPath, "config", "c", "configs/application.properties", "使用提供的配置文件启动服务器")
}

var (
	cmd      Server // 参数
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "lodge server web",
		Example: "server -c configs/application.properties",
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func setup() {
	viper.SetConfigFile(cmd.ConfigPath)
	viper.ReadInConfig()

}

func run() error {

	routers.Run()
	return nil
}
