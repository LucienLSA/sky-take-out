package global

var (
	ServiceName string
)

func InitConst() {
	ServiceName = Config.Jaeger.ServiceName
}
