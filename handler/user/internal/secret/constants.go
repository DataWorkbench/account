package secret

// for access key.
const (
	accessKeyIdLength  = 20
	secretKeyLength    = 40
	accessKeyIdLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	secretKeyLetters   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

const (
	SessionCacheSeconds = 3600 * 24
	sessionLength       = 50
	sessionLetters      = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

const (
	ProviderEnfi = "enfi"
)
