package contextkeys

type ContextKey string

const (
	UserID ContextKey = "user_id"
	Login  ContextKey = "login"
	Role   ContextKey = "role"
)
