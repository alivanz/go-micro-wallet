package microwallet

import crypto "github.com/alivanz/go-crypto"

type WriteableWallet interface {
	crypto.Wallet
	CurveName() (string, error)
	SetRandomPrivateKey(curvename string) error
	SetCurvePrivateKey(curvename string, pk []byte) error
}

type Bank interface {
	Open(int) (WriteableWallet, error)
}
