package rsa

type KeyPair struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type SignatureRequest struct {
	PrivateKey string `json:"privateKey" binding:"required"`
	Message    string `json:"message" binding:"required"`
	Algorithm  string `json:"algorithm" binding:"required"`
}

type VerificationRequest struct {
	PublicKey string `json:"publicKey" binding:"required"`
	Message   string `json:"message" binding:"required"`
	Signature string `json:"signature" binding:"required"`
	Algorithm string `json:"algorithm" binding:"required"`
}

type PublicKeyGenerationRequest struct {
	PrivateKey string `json:"privateKey" binding:"required"`
}
