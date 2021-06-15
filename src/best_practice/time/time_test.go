package time

// 你确信你了解时间吗    https://coolshell.cn/articles/5075.html
// 关于闰秒            https://coolshell.cn/articles/7804.html
// RFC 3339          https://tools.ietf.org/html/rfc3339

// 一定要使用 time.Time 和 time.Duration  这两个类型：
// - 在命令行上，flag 通过 time.ParseDuration 支持了 time.Duration。
// - JSON 中的 encoding/json 中也可以把time.Time 编码成 RFC 3339 的格式。
// - 数据库使用的 database/sql 也支持把 DATATIME 或 TIMESTAMP 类型转成 time.Time。
// - YAML 也可以使用 gopkg.in/yaml.v2 支持 time.Time 、time.Duration 和 RFC 3339 格式。
// - 如果要和第三方交互，实在没有办法也请使用 RFC 3339 的格式。
// - 如果要做全球化跨时区的应用，一定要把所有服务器和时间全部使用 UTC 时间。
