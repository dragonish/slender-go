# slender-go

- English
- [简体中文](./README_zhCN.md)

Website navigation.

## Environment variables

| Name | Type | Default | Description |
| --- | --- | --- | --- |
| `SLENDER_ACCESS_PWD` | `string` | `""` | Access password |
| `SLENDER_ADMIN_PWD` | `string` | `""` | Admin password |
| `SLENDER_LOG_LEVEL` | `string` | `Info` | Log output level. Optional: `Debug`, `Info`, `Warn`, `Error` |
| `SLENDER_PORT` | `int` | `8080` | Web service running port |
| `SLENDER_TOKEN_AGE` | `int` | `30` | Token age (days) |
| `SLENDER_PERFORMANCE_MODE` | `int` | `0` | Performance mode. *It is recommended to enable only when the database update performance is bad* |

## Startup commands

*Startup commands take precedence over environment variables.*

| Name | Type | Description |
| --- | --- | --- |
| `--debug, -D` || Enable debug mode |
| `--version, -v` || Show application version |
| `--help, -h` || Show help document |
| `--performance, -P` || Enable performance mode. *It is recommended to enable only when the database update performance is bad* |
| `--access_pwd, -a` | `string` | Specify access password |
| `--admin_pwd, -d` | `string` | Specify admin password |
| `--token_age, -t` | `int` | Specify token age (days) |
| `--log, -l` | `string` | Specify log output level. Optional: `Debug`, `Info`, `Warn`, `Error` |
| `--port, -p` | `int` | Specify web service running port |

## Admin password

If admin password is unset, default to access password(not empty) or `p@$$w0rd`.

## Feature

### Dynamic link

Convert dynamic links and display them according to the network environment.

Suppose the homepage address of the Slender service is `https://192.168.0.1:8080/`, the following available parameters and their corresponding results:

| parameter | result |
| --- | --- |
| `host` | `192.168.0.1:8080` |
| `hostname` | `192.168.0.1` |
| `href` | `https://192.168.0.1:8080/` |
| `origin` | `https://192.168.0.1:8080` |
| `pathname` | `/` |
| `port` | `8080` |
| `protocol` | `https:` |

**Example**

If a bookmark URL is configured as `https://{hostname}:8888/test`:

- When the homepage of the Slender service is `https://192.168.0.1:8080/`, it appears as `https://192.168.0.1:8888/test`.
- When the homepage of the Slender service is `https://172.17.0.1:8080/`, it appears as `https://172.17.0.1:8888/test`.
- When the homepage of the Slender service is `https://link.example.com/`, it appears as `https://link.example.com:8888/test`.

## Credits

- [soulteary/flare](https://github.com/soulteary/flare)
- [simple-icons/simple-icons](https://github.com/simple-icons/simple-icons)
- [Templarian/MaterialDesign](https://github.com/Templarian/MaterialDesign)
