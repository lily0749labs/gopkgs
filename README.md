# gopkgs

Go 语言常用包。

可以通过根包门面使用稳定的公共能力：

```go
import "github.com/lily0749labs/gopkgs"

rsp := (gopkgs.ApiRsp{}).Ok()
engine := gopkgs.NewEngine(gopkgs.EngineConfig{})
```

也可以直接导入相应子包，以获得更明确的依赖边界。

- `config`：从 `fs.FS` 加载 YAML 配置到 Viper。
- `apiRsp`：统一 API 响应。
- `ginUtil`：Gin Engine、CORS、Context 与验证错误翻译工具。
- `jwtUtil`：JWT 与 Token 工具。
- `log`：基于 Zap 的日志工具。
