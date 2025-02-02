package utils

type Config struct {
	Servers       []*ServerConfig `json:"servers"`
	MaxNickLength int             `json:"maxNickLength"`
}

type ServerConfig struct {
	Host            string   `json:"host"`
	Port            int      `json:"port"`
	SSL             bool     `json:"ssl"`
	Nick            string   `json:"nick"`
	Username        string   `json:"username"`
	RealName        string   `json:"realName"`
	DefaultChannels []string `json:"defaultChannels"`
}
