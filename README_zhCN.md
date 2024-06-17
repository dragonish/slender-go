# slender-go

- [English](./README.md)
- 简体中文

网址导航程序。

## 环境变量

| 变量名 | 值类型 | 默认值 | 描述 |
| --- | --- | --- | --- |
| `SLENDER_ACCESS_PWD` | `string` | `""` | 访问密码 |
| `SLENDER_ADMIN_PWD` | `string` | `""` | 管理员密码 |
| `SLENDER_LOG_LEVEL` | `string` | `Info` | 日志输出级别，可选值：`Debug`、`Info`、`Warn`、`Error` |
| `SLENDER_PORT` | `int` | `8080` | Web 服务运行端口 |
| `SLENDER_TOKEN_AGE` | `int` | `30` | 令牌保存期限 (天) |
| `SLENDER_PERFORMANCE_MODE` | `int` | `0` | 性能模式。*建议仅在遇到数据库更新性能不佳时才开启* |

## 启动命令

*启动命令的优先级高于环境变量。*

| 命令名 | 值类型 | 描述 |
| --- | --- | --- |
| `--debug, -D` || 开启调试模式 |
| `--version, -v` || 显示版本信息 |
| `--help, -h` || 显示帮助文档 |
| `--performance, -P` || 启用性能模式。*建议仅在遇到数据库更新性能不佳时才开启* |
| `--access_pwd, -a` | `string` | 指定访问密码 |
| `--admin_pwd, -d` | `string` | 指定管理员密码 |
| `--token_age, -t` | `int` | 指定令牌保存期限 (天) |
| `--log, -l` | `string` | 指定日志输出级别，可选值：`Debug`、`Info`、`Warn`、`Error` |
| `--port, -p` | `int` | 指定 Web 服务运行端口 |

## 管理员密码

如果未设置管理员密码，则默认为访问密码(非空时)或者 `p@$$w0rd`。

## 感谢

- [soulteary/flare](https://github.com/soulteary/flare)
- [simple-icons/simple-icons](https://github.com/simple-icons/simple-icons)
- [Templarian/MaterialDesign](https://github.com/Templarian/MaterialDesign)
