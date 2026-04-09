# gologging-ext

![Go Version](https://img.shields.io/badge/go-1.26+-blue?style=flat-square)
![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)
[![Go Report Card](https://goreportcard.com/badge/github.com/cnmax/gologging-ext)](https://goreportcard.com/report/github.com/cnmax/gologging-ext)
![Release](https://img.shields.io/github/v/release/cnmax/gologging-ext)

一个支持 **异步、批量、重试、多通道分发** 的 Go logging 扩展库，适用于高吞吐日志处理与消息通知场景。

---

## ✨特性

- 🚀 **异步处理**：解耦业务线程，提高性能
- 📦 **批量写入**：减少 IO / 网络开销
- 🔁 **失败重试**：提升日志可靠性
- 🔀 **多通道分发**：支持同时写入多个 sink
- 🔌 **插件化 adapters**：支持 zap、zerolog 等扩展

---

## 📦安装

```bash
go get github.com/cnmax/gologging-ext
```

## 🔌可选依赖（Adapters）

- **zap**
    ```bash
    go get github.com/cnmax/gologging-ext/adapters/zap
    ```
- **zerolog**
    ```bash
    go get github.com/cnmax/gologging-ext/adapters/zerolog
    ```

## 🚀快速开始

### slog（标准库）

```go
package tests

import (
    "context"
    "log/slog"
    "os"
    "testing"
    "time"

    slogAdapter "github.com/cnmax/gologging-ext/adapters/slog"
    "github.com/cnmax/gologging-ext/pipeline"
    "github.com/cnmax/gologging-ext/sinks/cls"
    "github.com/cnmax/gologging-ext/sinks/telegram"
    "github.com/cnmax/gologging-ext/sinks/wechat"
)

func TestSlogAdapter(t *testing.T) {

    ctx := context.Background()

    clsSink := cls.New(
        ctx,
        "ap-shanghai.cls.tencentcs.com",
        "your-secret-id",
        "your-secret-key",
        "your-topic-id",
    )

    wechatSink := wechat.New(
        ctx,
        "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx",
    )

    telegramSink := telegram.New(
        ctx,
        "your-bot-token",
        your-chat-id,
    )

    dispatcher := pipeline.NewDispatcher(
        pipeline.Build(clsSink),
        //pipeline.Build(telegramSink),
        pipeline.NewAsync(
            pipeline.NewRetry(telegramSink, 3, time.Second),
        ),
    )

    //logger := slog.New(slogAdapter.NewHandler(dispatcher, slog.LevelInfo))

    logger := slog.New(slog.NewMultiHandler(
        slog.NewTextHandler(os.Stdout, nil),
        slogAdapter.NewHandler(dispatcher, slog.LevelInfo),
        slogAdapter.NewHandler(pipeline.NewAsync(wechatSink), slog.LevelError),
    ))
    slog.SetDefault(logger)

    slog.Error("发生错误了", "user", "xxx", "age", 18)

    time.Sleep(10 * time.Second)
}
```

### zap

```go
package tests

import (
    "context"
    "os"
    "testing"
    "time"

    zapAdapter "github.com/cnmax/gologging-ext/adapters/zap"
    "github.com/cnmax/gologging-ext/pipeline"
    "github.com/cnmax/gologging-ext/sinks/cls"
    "github.com/cnmax/gologging-ext/sinks/telegram"
    "github.com/cnmax/gologging-ext/sinks/wechat"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func TestZapAdapter(t *testing.T) {

    ctx := context.Background()

    clsSink := cls.New(
        ctx,
        "ap-shanghai.cls.tencentcs.com",
        "your-secret-id",
        "your-secret-key",
        "your-topic-id",
    )

    wechatSink := wechat.New(
        ctx,
        "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx",
    )

    telegramSink := telegram.New(
        ctx,
        "your-bot-token",
        your-chat-id,
    )

    dispatcher := pipeline.NewDispatcher(
        pipeline.Build(clsSink),
        //pipeline.Build(telegramSink),
        pipeline.NewAsync(
            pipeline.NewRetry(telegramSink, 3, time.Second),
        ),
    )

    core := zapcore.NewTee(
        zapcore.NewCore(
            zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
            zapcore.AddSync(os.Stdout),
            zapcore.InfoLevel,
        ),
        zapAdapter.NewCore(dispatcher, zapcore.InfoLevel),
        zapAdapter.NewCore(pipeline.NewAsync(wechatSink), zapcore.ErrorLevel),
    )
    logger := zap.New(
        core,
        zap.AddCaller(),
        //zap.AddStacktrace(zapcore.ErrorLevel),
    )
    zap.ReplaceGlobals(logger)

    zap.L().Error(
        "发生错误了",
        zap.String("user", "xxx"),
        zap.Int("age", 18),
    )

    time.Sleep(10 * time.Second)
}
```

### zerolog

```go
package tests

import (
    "context"
    "os"
    "testing"
    "time"

    zerologAdapter "github.com/cnmax/gologging-ext/adapters/zerolog"
    "github.com/cnmax/gologging-ext/pipeline"
    "github.com/cnmax/gologging-ext/sinks/cls"
    "github.com/cnmax/gologging-ext/sinks/telegram"
    "github.com/cnmax/gologging-ext/sinks/wechat"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func TestZerologAdapter(t *testing.T) {

    ctx := context.Background()

    clsSink := cls.New(
        ctx,
        "ap-shanghai.cls.tencentcs.com",
        "your-secret-id",
        "your-secret-key",
        "your-topic-id",
    )

    wechatSink := wechat.New(
        ctx,
        "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx",
    )

    telegramSink := telegram.New(
        ctx,
        "your-bot-token",
        your-chat-id,
    )

    dispatcher := pipeline.NewDispatcher(
        pipeline.Build(clsSink),
        //pipeline.Build(telegramSink),
        pipeline.NewAsync(
            pipeline.NewRetry(telegramSink, 3, time.Second),
        ),
    )

    multi := zerolog.MultiLevelWriter(
        zerolog.ConsoleWriter{Out: os.Stdout},
        zerologAdapter.NewWriter(dispatcher, zerolog.InfoLevel),
        zerologAdapter.NewWriter(pipeline.NewAsync(wechatSink), zerolog.ErrorLevel),
    )

    logger := zerolog.
        New(multi).
        With().Timestamp().Logger().
        With().Caller().Logger()
    log.Logger = logger

    log.Error().
        Str("user", "xxx").
        Int("v", 18).
        Msg("发生错误了")

    time.Sleep(10 * time.Second)
}
```

## 🧩Sink

支持扩展多种输出目标（可结合 async / batch / retry 使用）：

### 腾讯云 CLS

用于日志上报至 CLS（Cloud Log Service）

```go
func New(ctx context.Context, endpoint, secretID, secretKey, topicID string, opts ...Option) core.Writer
```

### 参数
- `endpoint` CLS接入地址
- `secretID` 腾讯云密钥
- `secretKey` 腾讯云密钥
- `topicID` 日志主题 ID

### 微信机器人

用于企业微信群告警通知

```go
func New(ctx context.Context, webhook string, opts ...Option) core.Writer
```

### 参数
- `webhook` 企业微信机器人地址

### Telegram

用于机器人消息推送

```go
func New(ctx context.Context, token string, chatID int, opts ...Option) core.Writer
```

### 参数
- `token` Telegram Bot Token
- `chatID` 聊天 ID

### 钉钉机器人

用于企业告警通知

```go
func New(ctx context.Context, accessToken, secret string, opts ...Option) core.Writer
```

### 参数
- `accessToken` Webhook URL 中的 access_token
- `secret` 签名密钥（开启加签时必填）

### 飞书机器人

用于团队消息推送

```go
func New(ctx context.Context, webhook, secret string, opts ...Option) core.Writer
```

### 参数
- `webhook` 飞书机器人地址
- `secret` 签名密钥（开启加签时必填）

### 🔧自定义 Sink

实现 core.Writer 接口即可：

```go
type Writer interface {
    Write(entry *Entry) error
}
```

示例：

```go
type MySink struct{}

func (s *MySink) Write(entry *core.Entry) error {
    // 自定义处理逻辑
    return nil
}
```

## 📌适用场景
- 日志聚合 / 转发
- 告警 / 通知系统（飞书、钉钉、Webhook 等）
- 多目标日志输出（文件 + 云服务）
