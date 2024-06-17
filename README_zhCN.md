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

## 功能

### 动态链接

根据网络环境转换动态链接并展示。

假设 Slender 服务的首页地址为 `https://192.168.0.1:8080/`，以下可用的各参数及其对应解析结果：

| 参数名 | 解析结果 |
| --- | --- |
| `host` | `192.168.0.1:8080` |
| `hostname` | `192.168.0.1` |
| `href` | `https://192.168.0.1:8080/` |
| `origin` | `https://192.168.0.1:8080` |
| `pathname` | `/` |
| `port` | `8080` |
| `protocol` | `https:` |

**示例**

假设某书签网址配置为 `https://{hostname}:8888/test` 时:

- 当 Slender 服务的首页地址为 `https://192.168.0.1:8080/`，其显示为 `https://192.168.0.1:8888/test`。
- 当 Slender 服务的首页地址为 `https://172.17.0.1:8080/`，其显示为 `https://172.17.0.1:8888/test`。
- 当 Slender 服务的首页地址为 `https://link.example.com/`，其显示为 `https://link.example.com:8888/test`。

## 感谢

- [soulteary/flare](https://github.com/soulteary/flare)
- [simple-icons/simple-icons](https://github.com/simple-icons/simple-icons)
- [Templarian/MaterialDesign](https://github.com/Templarian/MaterialDesign)
