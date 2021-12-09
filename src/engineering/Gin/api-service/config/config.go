package config

import (
	"api-service/pkg/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{Name: cfg}

	// 初始化配置文件。
	if err := c.initConfig(); err != nil {
		return err
	}

	// 初始化日志包。
	c.initLog()
	return nil
}

// initConfig 初始化配置加载。
func (c *Config) initConfig() error {
	// 如果指定了配置文件，则解析指定的配置文件，否则解析默认的配置文件。
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath("conf")
		viper.SetConfigName("config")
	}
	// 设置配置文件格式为 YAML。
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// 读取匹配的环境变量前缀，比如 export API_SERVICE_URL=http://127.0.0.1:7777。
	viper.SetEnvPrefix("API_SERVICE")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 解析配置文件，启动后出现变化也支持热加载。
	if err := viper.ReadInConfig(); err != nil { // viper 解析配置文件
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})

	// 后续以 viper.GetXxx("a.b.c") 的方式读取配置项。
	return nil
}

// initLog 初始化日志
func (c *Config) initLog() {
	passLagerCfg := log.PassLagerCfg{
		Writers:        viper.GetString("log.writers"),
		LoggerLevel:    viper.GetString("log.logger_level"),
		LoggerFile:     viper.GetString("log.logger_file"),
		LogFormatText:  viper.GetBool("log.log_format_text"),
		RollingPolicy:  viper.GetString("log.rollingPolicy"),
		LogRotateDate:  viper.GetInt("log.log_rotate_date"),
		LogRotateSize:  viper.GetInt("log.log_rotate_size"),
		LogBackupCount: viper.GetInt("log.log_backup_count"),
	}

	_ = log.InitWithConfig(&passLagerCfg)
}
