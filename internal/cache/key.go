package cache

const (
	Prefix            = "sky-take-out:"
	KeyTokenSetPrefix = "employee:" // set; 保存登录用户及token
)

// 给redis key加上前缀
func GetRedisKey(key string) string {
	return Prefix + key
}
