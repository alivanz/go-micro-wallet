package microwallet

import (
	"crypto/ecdsa"
	"encoding/hex"
)

// SerializePubkey : Convert ECDSA Public Key to String
func SerializePubkey(pubkey *ecdsa.PublicKey) string {
	return hex.EncodeToString(append(pubkey.X.Bytes(), pubkey.Y.Bytes()...))
}
