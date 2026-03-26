package constant

type ctxKey string

const CtxKeyRequestID ctxKey = "__ctx_reqid"

const (
	MetadataKeyRequestID = "x-request-id"
	MetadataKeuAuth      = "authorization"
	MetadataKeyUserId    = "x-user-id"
	MetadataKeyUsername  = "x-user-username"
)
