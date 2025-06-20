package model

type EnvConfig struct {
	LogLvl    string `env:"LOGLEVEL" env-default:"info"`
	Token     string `env:"TOKEN" env-required:"true"`
	ChannelId int    `env:"CHANNEL_ID" env-required:"true"`
}

func NewEnvConfig() *EnvConfig {
	return &EnvConfig{}
}

type LocaleConfig struct {
	Commands `yaml:"commands"`
	Infos    `yaml:"infos"`
}

type Commands struct {
	StartMessage string `yaml:"start_message"`
	HelpMessage  string `yaml:"help_message"`
}
type Infos struct {
	TemplateMessageToChannel string `yaml:"template_message_to_channel"`
	SuccessMessage           string `yaml:"success_message"`
	WarnSendPhotoMessage     string `yaml:"warn_send_photo_message"`
}

func NewLocaleConfig() *LocaleConfig {
	return &LocaleConfig{}
}
