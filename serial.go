package microwallet

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/jacobsa/go-serial/serial"
)

type bank struct {
	rwc io.ReadWriteCloser
}
type wallet struct {
	rwc   io.ReadWriteCloser
	index int
}

func OpenBank(options *serial.OpenOptions) (Bank, error) {
	var err error
	if options == nil {
		options, err = DefaultConfig()
		if err != nil {
			return nil, err
		}
	}
	// Open the port.
	rwc, err := serial.Open(*options)
	if err != nil {
		return nil, err
	}
	return &bank{rwc}, nil
}

func (b *bank) Open(index int) (WriteableWallet, error) {
	return &wallet{b.rwc, index}, nil
}

func (w *wallet) Request(req map[string]string) (map[string]string, error) {
	for k, v := range req {
		fmt.Fprintf(w.rwc, "%s %s\n", k, v)
	}
	fmt.Fprintf(w.rwc, "\n")
	// resp
	reader := bufio.NewReader(w.rwc)
	resp := make(map[string]string)
	for {
		bline, _, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		line := string(bline)
		if len(line) == 0 {
			break
		}
		if !strings.HasPrefix(line, ">>>") {
			continue
		} else {
			line = line[3:]
		}
		kv := strings.SplitN(line, " ", 2)
		if len(kv) == 1 {
			break
		}
		resp[kv[0]] = kv[1]
	}
	return resp, nil
}
func (w *wallet) PubKey() (*ecdsa.PublicKey, error) {
	resp, err := w.Request(map[string]string{
		"method": "getpubkey",
		"index":  fmt.Sprint(w.index),
	})
	if err != nil {
		return nil, err
	}
	errmsg, iserror := resp["error"]
	if iserror {
		return nil, fmt.Errorf("%s", errmsg)
	}
	pubkey := &ecdsa.PublicKey{}
	if resp["curve"] == "secp256k1" {
		pubkey.Curve = btcec.S256()
	} else {
		return nil, fmt.Errorf("unknown curve")
	}
	x := resp["pubkey"][:64]
	y := resp["pubkey"][64:]
	pubkey.X = big.NewInt(0)
	pubkey.Y = big.NewInt(0)
	pubkey.X.SetString(x, 16)
	pubkey.Y.SetString(y, 16)
	return pubkey, nil
}
func (w *wallet) Sign(msghash []byte) (*big.Int, *big.Int, error) {
	resp, err := w.Request(map[string]string{
		"method":  "signdeterministic",
		"index":   fmt.Sprint(w.index),
		"msghash": hex.EncodeToString(msghash),
	})
	if err != nil {
		return nil, nil, err
	}
	errmsg, iserror := resp["error"]
	if iserror {
		return nil, nil, fmt.Errorf("%s", errmsg)
	}
	hr := resp["signature"][:64]
	hs := resp["signature"][64:]
	r := big.NewInt(0)
	s := big.NewInt(0)
	r.SetString(hr, 16)
	s.SetString(hs, 16)
	return r, s, nil
}
func (w *wallet) Verify(msghash []byte, r, s *big.Int) bool {
	pubkey, err := w.PubKey()
	if err != nil {
		return false
	}
	resp, err := w.Request(map[string]string{
		"method":    "verify",
		"curve":     "secp256k1",
		"pubkey":    SerializePubkey(pubkey),
		"msghash":   hex.EncodeToString(msghash),
		"signature": hex.EncodeToString(append(r.Bytes(), s.Bytes()...)),
	})
	_, iserror := resp["error"]
	if iserror {
		return false
	}
	return true
}

func (w *wallet) SetPrivateKey([]byte) error {
	return nil
}
func (w *wallet) SetCurvePrivateKey(curvename string, pk []byte) error {
	resp, err := w.Request(map[string]string{
		"method":  "setprivkey",
		"index":   fmt.Sprint(w.index),
		"curve":   curvename,
		"privkey": hex.EncodeToString(pk),
	})
	if err != nil {
		return err
	}
	errmsg, iserror := resp["error"]
	if iserror {
		return fmt.Errorf("%s", errmsg)
	}
	return nil
}

func (w *wallet) CurveName() (string, error) {
	resp, err := w.Request(map[string]string{
		"method": "getcurve",
		"index":  fmt.Sprint(w.index),
	})
	if err != nil {
		return "", err
	}
	errmsg, iserror := resp["error"]
	if iserror {
		return "", fmt.Errorf("%s", errmsg)
	}
	return resp["curve"], nil
}
func (w *wallet) SetRandomPrivateKey(curvename string) error {
	resp, err := w.Request(map[string]string{
		"method": "genprivkey",
		"index":  fmt.Sprint(w.index),
		"curve":  curvename,
	})
	if err != nil {
		return err
	}
	errmsg, iserror := resp["error"]
	if iserror {
		return fmt.Errorf("%s", errmsg)
	}
	return nil
}
