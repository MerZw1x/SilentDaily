package abstract

type HandlerContext interface {
	Get(key string, defaultValue ...string) string
	Status(status int) HandlerContext
	BindJSON(data any) error
}
