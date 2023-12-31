package utils

import (
	"crypto/ecdsa"
	"io/ioutil"
	"sync"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// Authorizer wraps an embedded bind.TransactOpts type with a mutex lock allowing for
// easier usage in concurrent programs. bind.TransactOpts is not thread-safe
// and as such must be used with mutex locks to prevent encountering any issues.
// Whenever using the embed bind.TransactOpts you must call Authorizer::Lock and
// Authorizer::Unlock to avoid any possible race conditions
type Authorizer struct {
	mx sync.Mutex
	*bind.TransactOpts
}

// NewAuthorizer returns an Authorizer object using a keyfile as the account source
func NewAuthorizer(keyFile, keyPass string) (*Authorizer, error) {
	fileBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}
	pk, err := keystore.DecryptKey(fileBytes, keyPass)
	if err != nil {
		return nil, errors.Wrap(err, "decrypt key")
	}
	return NewAuthorizerFromPK(pk.PrivateKey), nil
}

// NewAuthorizerFromPK returns an authorizer from a private key
func NewAuthorizerFromPK(pk *ecdsa.PrivateKey) *Authorizer {
	return &Authorizer{TransactOpts: bind.NewKeyedTransactor(pk)}
}

// NewKeyFile returns a new ethereum account as generated by `geth account new`
func NewKeyFile(keyFileDir, keyPass string) (accounts.Account, error) {
	return keystore.StoreKey(keyFileDir, keyPass, keystore.StandardScryptN, keystore.StandardScryptP)
}

// NewAccount creates a new ethereum private key, and associated transactor
func NewAccount() (*bind.TransactOpts, *ecdsa.PrivateKey, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, nil, err
	}
	return bind.NewKeyedTransactor(key), key, err
}

// Lock is used to claim a lock on the authorizer type
// and must be called before using it for transaction signing
func (a *Authorizer) Lock() {
	a.mx.Lock()
}

// Unlock is used to release a lock on the authorizer type
// and must be called after using it for transaction signing
func (a *Authorizer) Unlock() {
	a.mx.Unlock()
}
