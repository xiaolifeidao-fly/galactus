package vipper

import (
	"fmt"

	"github.com/spf13/viper"
)

// Init 初始化配置文件
func Init() {
	viper.SetConfigName("application")
	viper.SetConfigType("properties")
	viper.AddConfigPath("../configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
