package exception

const (
	InsufficientNft  Code = "INSUFFICIENT_NFT"
	InvalidToken     Code = "INVALID_TOKEN"
	InvalidSignature Code = "INVALID_SIGNATURE"
	Unknown          Code = "UNKNOWN"
)

type Code string
