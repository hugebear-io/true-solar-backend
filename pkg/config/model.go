package config

type ApplicationConfig struct {
	Timezone      string              `mapstructure:"timezone"`
	API           APIConfig           `mapstructure:"api"`
	Database      DatabaseConfig      `mapstructure:"db"`
	SNMP          SNMPConfig          `mapstructure:"snmp"`
	Redis         RedisConfig         `mapstructure:"redis"`
	ElasticSearch ElasticSearchConfig `mapstructure:"elasticsearch"`
	Logger        LoggerConfig        `mapstructure:"logger"`
	Cronjob       CronjobConfig       `mapstructure:"cronjob"`
}

type APIConfig struct {
	Port      string `mapstructure:"port"`
	Host      string `mapstructure:"host"`
	LogLevel  int    `mapstructure:"loglevel"`
	SecretKey string `mapstructure:"secretkey"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type SNMPConfig struct {
	TargetHost string `mapstructure:"target_host"`
	TargetPort string `mapstructure:"target_port"`
	AgentHost  string `mapstructure:"agent_host"`
}

type ElasticSearchConfig struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Index    string `mapstructure:"index"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type LoggerConfig struct {
	Level      int  `mapstructure:"level"`
	SkipCaller int  `mapstructure:"skipcaller"`
	Size       int  `mapstructure:"size"`
	Age        int  `mapstructure:"age"`
	Backup     int  `mapstructure:"backup"`
	Compress   bool `mapstructure:"compress"`
}

type CronjobConfig struct {
	PerformanceAlarmLow    string `mapstructure:"performance_alarm_low"`
	SumPerformanceAlarmLow string `mapstructure:"sum_performance_alarm_low"`
	Collector              string `mapstructure:"collector"`
}
