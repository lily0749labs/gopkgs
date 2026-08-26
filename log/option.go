package log

type Option struct {
	IsConsole  bool   `json:"isConsole" toml:"isConsole"`
	IsFile     bool   `json:"isFile" toml:"isFile"`
	SavePath   string `json:"savePath" toml:"savePath"`     // 保存路径
	MaxSize    int64  `json:"maxSize" toml:"maxSize"`       // 文件大小限制,单位MB
	MaxBackups int64  `json:"maxBackups" toml:"maxBackups"` // 最大保留日志文件数量
	MaxAge     int64  `json:"maxAge" toml:"maxAge"`         // 日志文件保留天数
	IsCompress bool   `json:"isCompress" toml:"isCompress"` // 是否压缩处理
}
