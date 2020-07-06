package zetabase

import (
	"crypto/ecdsa"
	"github.com/zetabase/zetabase-client/zbprotocol"
)

func MakeCredentialEcdsa(nonce int64, uid string, relBytes []byte, pk *ecdsa.PrivateKey) *zbprotocol.ProofOfCredential {
	r, s := MakeZetabaseSignature(uid, nonce, relBytes, pk)
	if r == "" || s == "" {
		return nil
	}
	return &zbprotocol.ProofOfCredential{
		CredType: zbprotocol.CredentialProofType_SIGNATURE,
		Signature: &zbprotocol.EcdsaSignature{
			R: r,
			S: s,
		},
		JwtToken: "",
	}
}

func MakeCredentialJwt(tok string) *zbprotocol.ProofOfCredential {
	return &zbprotocol.ProofOfCredential{
		CredType:  zbprotocol.CredentialProofType_JWT_TOKEN,
		Signature: EmptySignature(),
		JwtToken:  tok,
	}
}

func MakeEmptyCredentials() *zbprotocol.ProofOfCredential {
	return MakeCredentialJwt("debug mode only")
}
