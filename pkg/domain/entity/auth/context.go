package auth

const (
	sessionContextKey = "SESSION"
)

type contextKey struct {
	name string
}

var SessionContextKey = contextKey{sessionContextKey}
