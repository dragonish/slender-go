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

## Startup commands

*Startup commands take precedence over environment variables.*

| Name | Type | Description |
| --- | --- | --- |
| `--debug, -D` || Enable debug mode |
| `--version, -v` || Show application version |
| `--help, -h` || Show help document |
| `--access_pwd, -a` | `string` | Specify access password |
| `--admin_pwd, -d` | `string` | Specify admin password |
| `--token_age, -t` | `int` | Specify token age (days) |
| `--log, -l` | `string` | Specify log output level. Optional: `Debug`, `Info`, `Warn`, `Error` |
| `--port, -p` | `int` | Specify web service running port |

## Admin password

If admin password is unset, default to access password(not empty) or `p@$$w0rd`.

## Credits

- [soulteary/flare](https://github.com/soulteary/flare)
- [simple-icons/simple-icons](https://github.com/simple-icons/simple-icons)
- [Templarian/MaterialDesign](https://github.com/Templarian/MaterialDesign)
