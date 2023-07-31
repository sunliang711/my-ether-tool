package hd

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	btcutil "github.com/FactomProject/btcutilecc"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

const (
	PublicKeyCompressedLength = 33
	QUOTE_PREFIX              = 0x80000000
)

type Key struct {
	Path            string `json:"path"`
	PrivateKey      string `json:"secret_key"`
	PublicKey       string `json:"public_key"`
	EthereumAddress string `json:"ethereum_address"`
}

type OutputKey struct {
	Seed string `json:"seed,omitempty"`
	Keys []Key  `json:"keys"`
}

func (outputKey *OutputKey) JsonString() (string, error) {
	bytes, err := json.Marshal(outputKey)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (outputKey *OutputKey) String() string {
	var result string
	if outputKey.Seed != "" {
		result += fmt.Sprintf("seed: %s\n", outputKey.Seed)
	}
	for _, key := range outputKey.Keys {
		result += fmt.Sprintf(`
path: %s
private key: %s
public key: %s
ethereum address: %s`,
			key.Path, key.PrivateKey, key.PublicKey, key.EthereumAddress)
	}

	return result
}

func Derive(mnemonic string, passphrase string, path string, start, count uint) (*OutputKey, error) {
	seed := bip39.NewSeed(mnemonic, passphrase)

	outputKey := OutputKey{Seed: hexutil.Encode(seed)}

	rootKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("create master key error: %s", err)
	}

	if path == "" {
		privateKey := hexutil.Encode(rootKey.Key)
		pubkey := hexutil.Encode(rootKey.PublicKey().Key)
		address, err := PubkeyToAddress(pubkey)
		if err != nil {
			return nil, err
		}
		outputKey.Keys = append(outputKey.Keys, Key{Path: path, PrivateKey: privateKey, PublicKey: pubkey, EthereumAddress: address})
	} else {
		keys, paths, err := DerivesByPath(rootKey, path, start, count)
		if err != nil {
			return nil, err
		}
		for i := range keys {
			privateKey := hexutil.Encode(keys[i].Key)
			pubkey := hexutil.Encode(keys[i].PublicKey().Key)
			address, err := PubkeyToAddress(pubkey)
			if err != nil {
				return nil, err
			}
			outputKey.Keys = append(outputKey.Keys, Key{Path: paths[i], PrivateKey: privateKey, PublicKey: pubkey, EthereumAddress: address})
		}
	}
	return &outputKey, nil
}

func PrivateKeyToPublicKey(privateKey string) (string, error) {
	bz, err := hexutil.Decode(privateKey)
	if err != nil {
		return "", nil
	}
	pubBytes := PublicKeyForPrivateKey(bz)
	return hexutil.Encode(pubBytes), nil
}

func PublicKeyForPrivateKey(key []byte) []byte {
	curve := btcutil.Secp256k1()
	return CompressPublicKey(curve.ScalarBaseMult(key))
}

func CompressPublicKey(x *big.Int, y *big.Int) []byte {
	var key bytes.Buffer

	// Write header; 0x2 for even y value; 0x3 for odd
	key.WriteByte(byte(0x2) + byte(y.Bit(0)))

	// Write X coord; Pad the key so x is aligned with the LSB. Pad size is key length - header size (1) - xBytes size
	xBytes := x.Bytes()
	for i := 0; i < (PublicKeyCompressedLength - 1 - len(xBytes)); i++ {
		key.WriteByte(0x0)
	}
	key.Write(xBytes)

	return key.Bytes()
}

func IsCompressedPublicKey(pubkey []byte) (bool, error) {
	if len(pubkey) == 33 {
		if pubkey[0] == 2 || pubkey[0] == 3 {
			return true, nil
		} else {
			return false, errors.New("invalid pubkey")
		}
	} else if len(pubkey) == 65 {
		if pubkey[0] == 4 {
			return false, nil
		} else {
			return false, errors.New("invalid pubkey")
		}
	}
	return false, errors.New("invalid pubkey")
}

func PubkeyToAddress(pubkey string) (string, error) {
	var (
		publicKey *ecdsa.PublicKey
	)
	p, err := hexutil.Decode(pubkey)
	if err != nil {
		return "", err
	}
	compressed, err := IsCompressedPublicKey(p)
	if err != nil {
		return "", err
	}

	if compressed {
		publicKey, err = crypto.DecompressPubkey(p)
		if err != nil {
			return "", err
		}
	} else {
		publicKey, err = crypto.UnmarshalPubkey(p)
		if err != nil {
			return "", err
		}
	}
	address := crypto.PubkeyToAddress(*publicKey)
	return address.Hex(), nil

}

func DeriveByPath(key *bip32.Key, path string) (*bip32.Key, error) {
	if !strings.HasPrefix(path, "m/") {
		return nil, errors.New("invalid path prefix")
	}

	path = strings.TrimPrefix(path, "m/")
	pathIndice := strings.Split(path, "/")
	childKey := key
	for _, pathIndex := range pathIndice {
		quote := false
		if strings.HasSuffix(pathIndex, "'") {
			pathIndex = strings.TrimRight(pathIndex, "'")
			quote = true
		}
		index, err := strconv.Atoi(pathIndex)
		if err != nil {
			return nil, errors.New("invalid path field")
		}

		if quote {
			index += QUOTE_PREFIX
		}
		childKey, err = childKey.NewChildKey(uint32(index))
		if err != nil {
			return nil, err
		}
	}

	return childKey, nil
}

// path format: m/60'/44'/x/0
// x is place holder to derive from start to start + count - 1  secret key
func DerivesByPath(key *bip32.Key, path string, start, count uint) (keys []*bip32.Key, paths []string, err error) {
	if strings.Count(path, "x") == 0 {
		key, err = DeriveByPath(key, path)
		keys = append(keys, key)
		paths = append(paths, path)
		return
	}

	if strings.Count(path, "x") != 1 {
		return nil, nil, errors.New("invalid path format,path must has one x")
	}

	newPath := strings.Replace(path, "x", "%d", 1)
	for i := start; i < start+count; i++ {
		path = fmt.Sprintf(newPath, i)
		key, err := DeriveByPath(key, path)
		if err != nil {
			return nil, nil, err
		}
		keys = append(keys, key)
		paths = append(paths, path)
	}
	return
}

func CheckHdPath(path string) error {
	path = strings.Replace(path, "x", "0", 1)
	_, err := accounts.ParseDerivationPath(path)
	return err
}
