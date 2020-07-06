package zetabase

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/zetabase/zetabase-client/zbprotocol"
	"math/big"
	"github.com/spaolacci/murmur3"
)


// PUBLIC-FACING FUNCTIONS

func MakeZetabaseSignature(uid string, nonce int64, relData []byte, pk *ecdsa.PrivateKey) (string, string) {
	sbs := signingBytes(uid, nonce)
	return signMessageBytes(sbs, relData, pk)
}

func ValidateZetabaseSignature(uid string, nonce int64, relData []byte, pk *ecdsa.PublicKey, r, s string) bool {
	sbs := signingBytes(uid, nonce)
	return validateSignatureBytes(pk, sbs, relData, r, s)
}

// INTERNAL CRYPTOGRAPHY FUNCTIONS

func signingBytes(uid string, nonce int64) []byte {
	sbs := fmt.Sprintf("%s%d", uid, nonce)
	return []byte(sbs)
}

func permissionSigningBytes(entry zbprotocol.PermissionsEntry) []byte {
	b1 := []byte{byte(entry.AudienceType), byte(entry.Level)}
	b2 := []byte(entry.Id + entry.TableId + entry.AudienceId)
	var b3 []byte
	for _, p := range entry.Constraints {
		b3 = append(b3, byte(p.FieldConstraint.ConstraintType))
		b3 = append(b3, byte(p.FieldConstraint.ValueType))
		s := []byte(p.FieldConstraint.FieldKey + p.FieldConstraint.RequiredValue)
		b3 = append(b3, s...)
	}
	var bs []byte
	bs = b1
	bs = append(bs, b2...)
	bs = append(bs, b3...)
	return bs
}

func permissionSetSigningBytes(perms []*zbprotocol.PermissionsEntry) []byte {
	var bs []byte
	for _, p := range perms {
		if p != nil {
			b := permissionSigningBytes(*p)
			bs = append(bs, b...)
		}
	}
	return bs
}

func PermissionsEntrySigningBytes(perm *zbprotocol.PermissionsEntry) []byte {
	return permissionSetSigningBytes([]*zbprotocol.PermissionsEntry{perm})
}

func TableCreateSigningBytes(tblId string, perms []*zbprotocol.PermissionsEntry) []byte {
	bs := []byte(tblId)
	bs2 := permissionSetSigningBytes(perms)
	bs = append(bs, bs2...)
	return bs
}

// MULTI-PUT EXTRA SIGNING BYTES AND VARIANTS

func MultiPutExtraSigningBytesMd5(pairs []*zbprotocol.DataPair) []byte {
	hash := md5.New()
	for _, v := range pairs {
		hash.Write([]byte(v.GetKey()))
		hash.Write(v.GetValue())
	}
	bs := hash.Sum(nil)
	return bs
}

func MultiPutExtraSigningBytesMurmur3(pairs []*zbprotocol.DataPair) []byte {
	hash := murmur3.New32WithSeed(1234)
	for _, v := range pairs {
		hash.Write([]byte(v.GetKey()))
		hash.Write(v.GetValue())
	}
	bs := hash.Sum(nil)
	return bs
}

func MultiPutExtraSigningBytesMurmur3Sliding(pairs []*zbprotocol.DataPair) []byte {
	hash := murmur3.New32WithSeed(1234)
	stdLen := 64
	for i, v := range pairs {
		startIdx := i % len(v.GetValue())
		endIdx := startIdx + stdLen
		if endIdx > len(v.GetValue()) {
			endIdx = len(v.GetValue())
		}
		valu := v.GetValue()[startIdx:endIdx]
		hash.Write([]byte(v.GetKey()))
		hash.Write(valu)
	}
	bs := hash.Sum(nil)
	return bs
}

func MultiPutExtraSigningBytes(pairs []*zbprotocol.DataPair) []byte {
	hash := md5.New()
	for _, v := range pairs {
		hash.Write(v.GetValue())
	}
	bs := hash.Sum(nil)
	return bs
}

func MultiPutExtraSigningBytesAbbrev(pairs []*zbprotocol.DataPair) []byte {
	//hash := murmur3.New32WithSeed(1234)
	hash := md5.New()
	nBytes := 64
	everyNth := 128
	//hash := md5.New()
	for i, v := range pairs {
		if i % everyNth != 0 {
			continue
		}
		valu := v.GetValue()
		if len(valu) > nBytes {
			valu = valu[:nBytes]
		}
		hash.Write(valu)
	}
	bs := hash.Sum(nil)
	return bs
}
// -----------------------------------------------

func TablePutExtraSigningBytes(key string, valu []byte) []byte {
	hash := md5.New()
	hash.Write([]byte(key))
	hash.Write(valu)
	return hash.Sum(nil)
}

func signMessageBytes(byts []byte, relDataBytes []byte, pk *ecdsa.PrivateKey) (string, string) {
	hash := sha256.New()
	hash.Write(byts)
	hash.Write(relDataBytes)
	bytHash := hash.Sum(nil)
	//bytHash := sha256HashBytes(byts)
	i1, i2 := signBytes(bytHash, pk)
	s1, s2 := signatureToStrings(i1, i2)
	return s1, s2
}

func sha256HashBytes(bs []byte) []byte {
	h := sha256.New()
	h.Write(bs)
	return h.Sum(nil)
}

func GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	pubkeyCurve := elliptic.P256() //see http://golang.org/pkg/crypto/elliptic/#P256

	privatekey := new(ecdsa.PrivateKey)
	privatekey, err := ecdsa.GenerateKey(pubkeyCurve, rand.Reader) // this generates a public & private key pair
	if err != nil {
		return nil, nil
	}
	var pubkey ecdsa.PublicKey
	pubkey = privatekey.PublicKey
	return privatekey, &pubkey
}

func signBytes(bs []byte, privKey *ecdsa.PrivateKey) (*big.Int, *big.Int) {
	r, s, err := ecdsa.Sign(rand.Reader, privKey, bs)
	if err != nil {
		return nil, nil
	}
	return r, s
}

func signatureToStrings(a, b *big.Int) (string, string) {
	x, y := a.String(), b.String()
	return x, y
}

func stringsToSignature(a, b string) (*big.Int, *big.Int) {
	x := new(big.Int)
	x.SetString(a, 10)
	y := new(big.Int)
	y.SetString(b, 10)
	return x, y
}

func validateSignatureBytes(pubKey *ecdsa.PublicKey, stdSigningBytes []byte, specialDataBytes []byte, r, s string) bool {
	ri, si := stringsToSignature(r, s)
	hash := sha256.New()
	hash.Write(stdSigningBytes)
	hash.Write(specialDataBytes)
	data := hash.Sum(nil)
	res := ecdsa.Verify(pubKey, data, ri, si)
	return res
}

func EncodeEcdsaPublicKey(publicKey *ecdsa.PublicKey) ([]byte, error) {
	x509EncodedPub, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return pemEncodedPub, nil
}

func EncodeEcdsaPrivateKey(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	x509Encoded, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return pemEncoded, nil
}

func DecodeEcdsaPublicKey(pemEncodedPub string) (*ecdsa.PublicKey, error) {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return publicKey, nil
}

func DecodeEcdsaPrivateKey(pemEncoded string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	return privateKey, nil
}
