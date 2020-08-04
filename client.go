/*
Package zetabase provides a Go client library for the Zetabase cloud database along
with a command-line tool for managing data and users.

Most users will only require the ZetabaseClient class; other exposed functions are
primarily of use for working at the protocol layer, e.g. for performance fine-tuning.

The client library communicates using a protocol based on gRPC/Protocol Buffers.

This code and the included command-line tool are licensed under the Modified BSD License
included in the attached LICENSE file.
*/
package zetabase

import (
	"context"
	"crypto/ecdsa"
	"crypto/tls"
	"errors"
	"github.com/zetabase/zetabase-client/zbprotocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	ClientVersion = "0.1"
	ClientVersionCode = "Curry"
)

// Type ZetabaseClient represents a long-lived Zetabase connection for a particular identity.
type ZetabaseClient struct {
	userId       string
	serverAddr   string
	insecure     bool
	noCertVerify bool
	parentId     *string
	privKey      *ecdsa.PrivateKey
	pubKey       *ecdsa.PublicKey
	loginId      *string
	password     *string
	nonceMaker   *NonceMaker
	conn         *grpc.ClientConn
	client       zbprotocol.ZetabaseProviderClient
	jwtToken     *string
	debugMode    bool
	ctx          context.Context
	maxItemSize  int64
}

// Creates a new client for a given user ID uid. The user ID should be in UUID form.
func NewZetabaseClient(uid string) *ZetabaseClient {
	return &ZetabaseClient{
		userId:       uid,
		serverAddr:   "api.zetabase.io:443",
		insecure:     false,
		noCertVerify: false,
		debugMode:    false,
		parentId:     nil,
		loginId:      nil,
		privKey:      nil,
		pubKey:       nil,
		password:     nil,
		nonceMaker:   NewNonceMaker(),
		conn:         nil,
		client:       nil,
		jwtToken:     nil,
		ctx:          context.Background(),
		maxItemSize:  int64(1000),
	}
}

func NewZetabaseUserClient(parentId string) *ZetabaseClient {
	return &ZetabaseClient{
		userId:       "",
		serverAddr:   "api.zetabase.io:443",
		insecure:     false,
		noCertVerify: false,
		debugMode:    false,
		parentId:     &parentId,
		loginId:      nil,
		privKey:      nil,
		pubKey:       nil,
		password:     nil,
		nonceMaker:   NewNonceMaker(),
		conn:         nil,
		client:       nil,
		jwtToken:     nil,
		ctx:          context.Background(),
		maxItemSize:  int64(1000),
	}
}

// Checks version compatibility between client and server
func (z *ZetabaseClient) CheckVersion() (bool, *zbprotocol.VersionDetails, error) {
	if !z.checkReady() {
		return false, nil, errors.New("NotReady")
	}
	info, err := z.client.VersionInfo(z.ctx, &zbprotocol.ZbEmpty{})
	if err != nil {
		return false, nil, err
	}
	minClientVersion := info.GetMinClientVersion()
	isEnough := IsSemVerVersionAtLeast(ClientVersion, minClientVersion)
	return isEnough, info, nil
}

// Toggle certificate verification
func (z *ZetabaseClient) SetCertVerify(b bool) {
	z.noCertVerify = b
}

// Toggle insecure (plaintext) connection
func (z *ZetabaseClient) SetInsecure() {
	z.insecure = true
}

// Toggle debug mode
func (z *ZetabaseClient) SetDebugMode() {
	z.debugMode = true
}

// Set parent user ID (when connecting as a subuser) to id. This parent ID should be in UUID form.
func (z *ZetabaseClient) SetParent(id string) {
	z.parentId = &id
}

// Set the private and public keys corresponding to this identity
func (z *ZetabaseClient) SetIdKey(priv *ecdsa.PrivateKey, pub *ecdsa.PublicKey) {
	z.privKey = priv
	z.pubKey = pub
}

// Set a customer server address (for on-premises installations only)
func (z *ZetabaseClient) SetServerAddr(addr string) {
	z.serverAddr = addr
}

// Set a handle and password as authentication credentials (for JWT authentication mode)
func (z *ZetabaseClient) SetIdPassword(loginId, pwd string) {
	z.password = &pwd
	z.loginId = &loginId
}

// Check if client is ready to communicate with server
func (z *ZetabaseClient) checkReady() bool {
	if z.privKey != nil || (z.password != nil && z.loginId != nil) {
		if z.conn != nil {
			if z.loginId != nil && z.jwtToken == nil {
				err := z.authLoginJwt()
				if err != nil {
					return false
				}
			}
			return true
		}
	}
	return false
}

// Function ConfirmUserIdentity confirms a new user ID with the SMS confirmation code they received.
func ConfirmUserIdentity(ctx context.Context, uid string, parentId *string, verifyStr string, client zbprotocol.ZetabaseProviderClient) error {
	parStr := ""
	if parentId != nil {
		parStr = *parentId
	}
	_, err := client.ConfirmNewIdentity(ctx, &zbprotocol.NewIdentityConfirm{
		Id:               uid,
		ParentId:         parStr,
		VerificationCode: verifyStr,
	})
	return err
}

func (z *ZetabaseClient) jwtCredential() *zbprotocol.ProofOfCredential {
	if z.jwtToken != nil {
		return MakeCredentialJwt(*z.jwtToken)
	}
	return nil
}

func (z *ZetabaseClient) ecdsaCredential(nonce int64, extraBytes []byte) *zbprotocol.ProofOfCredential {
	return MakeCredentialEcdsa(nonce, z.userId, extraBytes, z.privKey)
}

func (z *ZetabaseClient) authLoginJwt() error {
	if z.conn == nil {
		return errors.New("NotReady")
	} else if z.password == nil {
		return errors.New("NoPasswordProvided")
	}
	parId := ""
	if z.parentId != nil {
		parId = *z.parentId
	}
	res, err := z.client.LoginUser(z.ctx, &zbprotocol.AuthenticateUser{
		ParentId:   parId,
		Handle:     *z.loginId,
		Password:   *z.password,
		Nonce:      z.nonceMaker.Get(),
		Credential: MakeEmptyCredentials(),
	})
	if err != nil {
		return err
	} else {
		j := res.GetJwtToken()
		if len(j) > 0 {
			z.jwtToken = &j
			z.userId = res.GetId()
		}
	}
	return nil
}

func (z *ZetabaseClient) getCredential(nonce int64, xBytes []byte) *zbprotocol.ProofOfCredential {
	var poc *zbprotocol.ProofOfCredential
	if z.jwtToken != nil {
		poc = z.jwtCredential()
	} else {
		poc = z.ecdsaCredential(nonce, xBytes)
	}
	return poc
}

// Method ListTables lists the tables associated with the ZetabaseClient's account
func (z *ZetabaseClient) ListTables() ([]string, error) {
	if !z.checkReady() {
		return nil, errors.New("NotReady")
	}
	
	tableNames := []string{}

	nonce := z.nonceMaker.Get()
	poc := z.getCredential(nonce, nil)

	res, err := z.client.ListTables(z.ctx, &zbprotocol.ListTablesRequest{
		Id:            z.userId,
		Nonce:         nonce,
		TableOwnerId:  z.userId,
		Credential:    poc,
	})
	
	if err != nil {
		return nil, err
	}
	
	for _, table := range(res.GetTableDefinitions()) {
		tableNames = append(tableNames, table.GetTableId())
	}
	return tableNames, nil
}

// Method ListKeys lists the keys for a given table
func (z *ZetabaseClient) ListKeys(tableOwnerId, tableId string) *PaginationHandler {
	return z.ListKeysWithPattern(tableOwnerId, tableId, "")
}

func unwrapZbError(zbError *zbprotocol.ZbError) error {
	if zbError == nil {
		return nil
	} else if zbError.Code == 0 && len(zbError.Message) > 0 {
		return errors.New(zbError.Message)
	}
	return nil
}

// Method PutMulti puts multiple key-value pairs into a table at once
func (z *ZetabaseClient) putMultiRaw(tableOwnerId, tableId string, keys []string, valus [][]byte, overwrite bool) error {
	if len(valus) != len(keys) {
		return errors.New("ImproperDimensions")
	}
	nonce := z.nonceMaker.Get()
	var dps []*zbprotocol.DataPair
	for i := 0; i < len(keys); i++ {
		dps = append(dps, &zbprotocol.DataPair{
			Key:   keys[i],
			Value: valus[i],
		})
	}
	//xBytes := MultiPutExtraSigningBytes(dps)
	xBytes := MultiPutExtraSigningBytesMd5(dps)
	//log.Printf("Multi extra signing bytes: %x\n", xBytes)
	cred := z.getCredential(nonce, xBytes)
	res, err := z.client.PutDataMulti(z.ctx, &zbprotocol.TablePutMulti{
		Id:           z.userId,
		TableOwnerId: tableOwnerId,
		TableId:      tableId,
		Overwrite:    overwrite,
		Nonce:        nonce,
		Credential:   cred,
		Pairs:        dps,
	})
	if err != nil {
		return err
	} else {
		return unwrapZbError(res)
	}
}

const (
	GrpcMaxBytes = 4000000
)

// Method PutMulti puts multiple key-value pairs into a table at once
func (z *ZetabaseClient) PutMulti(tableOwnerId, tableId string, keys []string, valus [][]byte, overwrite bool) error {
	if len(valus) != len(keys) {
		return errors.New("ImproperDimensions")
	}

	maxBytes := GrpcMaxBytes / 2

	pgs := makePutPages(z, keys, valus, uint64(maxBytes))
	err := pgs.putAll(tableOwnerId, tableId, overwrite)

	if err != nil {
		return err
	}
	return nil
}

func (z *ZetabaseClient) SetMaxItemSize(newSize int64) {
	z.maxItemSize = newSize
}

func (z *ZetabaseClient) Get(tableOwnerId, tableId string, keys []string) *getPages {
	getPages := MakeGetPages(z, keys, z.maxItemSize, tableOwnerId, tableId)
	return getPages
}

// Method Get fetches a given set of keys from a table and returns a PaginationHandler object.
func (z *ZetabaseClient) getPag(tableOwnerId, tableId string, keys []string) *PaginationHandler {
	f := func(idx int64) (map[string][]byte, bool, error) {
		tim, hasNxt, err := z.get(tableOwnerId, tableId, keys, idx)
		if err == nil {
			return tim, hasNxt, nil
		} else {
			return nil, false, err
		}
	}
	return StandardPaginationHandlerFor(f)
}

func (z *ZetabaseClient) get(tableOwnerId, tableId string, keys []string, pageIdx int64) (map[string][]byte, bool, error) {
	if !z.checkReady() {
		return nil, false, errors.New("NotReady")
	}
	nonce := z.nonceMaker.Get()
	poc := z.getCredential(nonce, nil)
	res, err := z.client.GetData(z.ctx, &zbprotocol.TableGet{
		Id:           z.userId,
		TableOwnerId: tableOwnerId,
		TableId:      tableId,
		Nonce:        nonce,
		Credential:   poc,
		PageIndex:    pageIdx,
		Keys:         keys,
	})
	if err != nil {
		return nil, false, err
	}
	m := map[string][]byte{}
	for _, x := range res.GetData() {
		m[x.GetKey()] = x.GetValue()
	}
	return m, res.GetPagination().GetHasNextPage(), nil
}

// Get user ID
func (z *ZetabaseClient) Id() string {
	return z.userId
}

// Confirm new subuser givern subuserId and verification code
func (z *ZetabaseClient) ConfirmNewSubUser(subuserId, verificationCode string) error {
	_, err := z.client.ConfirmNewIdentity(z.ctx, &zbprotocol.NewIdentityConfirm{
		Id:               subuserId,
		ParentId:         z.userId,
		VerificationCode: verificationCode,
	})
	if err != nil {
		return err
	}
	return nil
}

// Create a new subuser with the given attributes
func (z *ZetabaseClient) NewSubUser(handle, email, mobile, password, signupCode, groupId string, pubKey *ecdsa.PublicKey) (string, error) {
	pkBs, err := EncodeEcdsaPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	res, err := z.client.CreateUser(z.ctx, &zbprotocol.NewSubIdentityRequest{
		Id:            z.userId,
		Name:          handle,
		Email:         email,
		Mobile:        mobile,
		LoginPassword: password,
		PubKeyEncoded: string(pkBs),
		SignupCode:    signupCode,
		GroupId:       groupId,
	})
	if err != nil {
		return "", err
	}
	id := res.Id
	return id, nil
}

// Create a new table tblId with the given data format, indexed fields, and permissions
func (z *ZetabaseClient) CreateTable(tblId string, dataType zbprotocol.TableDataFormat, indexedFields []*IndexedField, perms []*PermEntry, allowJwt bool) error {
	if !z.checkReady() {
		return errors.New("NotReady")
	}
	var pEntries []*zbprotocol.PermissionsEntry
	for _, p := range perms {
		pEntries = append(pEntries, p.ToProtocol(z.userId, tblId))
	}
	nonce := z.nonceMaker.Get()
	sigBytes := TableCreateSigningBytes(tblId, pEntries)
	cred := z.getCredential(nonce, sigBytes)
	tc := &zbprotocol.TableCreate{
		Id:             z.userId,
		TableId:        tblId,
		DataFormat:     dataType,
		Indices:        indexedFieldsToProtocol(indexedFields),
		Nonce:          nonce,
		AllowTokenAuth: allowJwt,
		Credential:     cred,
		Permissions:    pEntries,
	}
	res, err := z.client.CreateTable(z.ctx, tc)
	if err != nil {
		return err
	} else if res.Code != 0 && len(res.Message) > 0 {
		return errors.New(res.Message)
	}
	return nil
}

// Method AddPermission adds permission perm to the given table tblId.
func (z *ZetabaseClient) AddPermission(tblOwnerId, tblId string, perm *PermEntry) error {
	if !z.checkReady() {
		return errors.New("NotReady")
	}
	nonce := z.nonceMaker.Get()
	permsEnt := perm.ToProtocol(tblOwnerId, tblId)
	permsEnt.Nonce = nonce
	poc := z.getCredential(nonce, PermissionsEntrySigningBytes(permsEnt))
	permsEnt.Credential = poc

	_, err := z.client.SetPermission(z.ctx, permsEnt)
	if err != nil {
		return err
	}
	return nil
}

// Method ListKeysWithPattern lists keys with a given prefix pattern, where the suffix wildcard operator
// is represented by %.
func (z *ZetabaseClient) ListKeysWithPattern(tableOwnerId, tableId, pattern string) *PaginationHandler {
	f := func(idx int64) (map[string][]byte, bool, error) {
		m := map[string][]byte{}
		tim, hasNxt, err := z.listKeysWithPattern(tableOwnerId, tableId, pattern, idx)
		if err == nil {
			for _, k := range tim {
				m[k] = nil
			}
			return m, hasNxt, nil
		} else {
			return nil, false, err
		}
	}
	return StandardPaginationHandlerFor(f)
}

// Method GetSubIdentities lists subusers of the authenticated user.
func (z *ZetabaseClient) GetSubIdentities() ([]*zbprotocol.NewSubIdentityRequest, error) {
	if !z.checkReady() {
		return nil, errors.New("NotReady")
	}
	nonce := z.nonceMaker.Get()
	poc := z.getCredential(nonce, nil)
	res, err := z.client.ListSubIdentities(z.ctx, &zbprotocol.SimpleRequest{
		Id:         z.userId,
		Nonce:      nonce,
		Credential: poc,
	})
	if err != nil {
		return nil, err
	} else {
		return res.GetSubIdentities(), nil
	}
}

// Method ModifySubIdentity modifies an existing subuser. Non-nil fields will be updated.
func (z *ZetabaseClient) ModifySubIdentity(subUserId string, newHandle *string, newEmail *string, newMobile *string, newPass *string, newPubKey *string) error {
	var email, mobile, pass, pubkey, name string
	if newEmail != nil {
		email = *newEmail
	}
	if newMobile != nil {
		mobile = *newMobile
	}
	if newPass != nil {
		pass = *newPass
	}
	if newPubKey != nil {
		pubkey = *newPubKey
	}
	if newHandle != nil {
		name = *newHandle
	}
	nonce := z.nonceMaker.Get()
	poc := z.getCredential(nonce, nil)
	_, err := z.client.ModifySubIdentity(z.ctx, &zbprotocol.SubIdentityModify{
		Id:                   z.userId,
		SubId:                subUserId,
		NewName:              name,
		NewEmail:             email,
		NewMobile:            mobile,
		NewPassword:          pass,
		NewPubKey:            pubkey,
		Nonce:                nonce,
		Credential:           poc,
	})
	return err
}

func (z *ZetabaseClient) Query(tableOwnerId, tableId string, qry0 SubQueryConvertible) *PaginationHandler {
	qry := qry0.ToSubQuery(tableOwnerId, tableId)
	f := func(idx int64) (map[string][]byte, bool, error) {
		m := map[string][]byte{}
		tim, hasNxt, err := z.query(tableOwnerId, tableId, idx, qry)
		if err == nil {
			for _, k := range tim {
				m[k] = nil
			}
			return m, hasNxt, nil
		} else {
			return nil, false, err
		}
	}
	return StandardPaginationHandlerFor(f)
}

func (z *ZetabaseClient) query(tblOwnerId, tblId string, pgIdx int64, qry *zbprotocol.TableSubQuery) ([]string, bool, error) {
	if !z.checkReady() {
		return nil, false, errors.New("NotReady")
	}
	nonce := z.nonceMaker.Get()
	poc := z.getCredential(nonce, nil)
	res, err := z.client.QueryKeys(z.ctx, &zbprotocol.TableQuery{
		Id:           z.userId,
		TableOwnerId: tblOwnerId,
		TableId:      tblId,
		Query:        qry,
		Nonce:        nonce,
		PageIndex:    pgIdx,
		Credential:   poc,
	})
	if err != nil {
		return nil, false, err
	} else {
		keys := res.GetKeys()
		return keys, res.GetPagination().GetHasNextPage(), nil
	}
}

func (z *ZetabaseClient) QueryData(tbldOwnerId, tblId string, qry SubQueryConvertible) (*getPages, error) {
	res := z.Query(tbldOwnerId, tblId, qry)
	data, err := res.DataAll()

	if err != nil {
		return nil, err
	}

	keys := []string{}
	for k := range(data) {
		keys = append(keys, k)
	}

	tblData := z.Get(tbldOwnerId, tblId, keys)
	return tblData, nil
}

// Put a given key-value pair into a table
func (z *ZetabaseClient) PutData(tableOwnerId, tableId, key string, valu []byte, overwrite bool) error {
	if !z.checkReady() {
		return errors.New("NotReady")
	}
	nonce := z.nonceMaker.Get()
	xBytes := TablePutExtraSigningBytes(key, valu)
	poc := z.getCredential(nonce, xBytes)
	_, err := z.client.PutData(z.ctx, &zbprotocol.TablePut{
		Id:           z.userId,
		TableOwnerId: tableOwnerId,
		TableId:      tableId,
		Key:          key,
		Value:        valu,
		Overwrite:    overwrite,
		Nonce:        nonce,
		Credential:   poc,
	})
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Delete a given key-value pair from a table
func (z *ZetabaseClient) DeleteKey(tableOwnerId, tableId, key string) error {
	if !z.checkReady() {
		return errors.New("NotReady")
	}
	nonce := z.nonceMaker.Get()
	extraBytes := []byte(key)
	poc := z.getCredential(nonce, extraBytes)

	_, err := z.client.DeleteObject(z.ctx, &zbprotocol.DeleteSystemObjectRequest{
		Id:           z.userId,
		ObjectType:   zbprotocol.SystemObjectType_KEY,
		TableOwnerId: tableOwnerId,
		TableId:      tableId,
		ObjectId:     key,
		Nonce:        nonce,
		Credential:   poc,
	})

	if err != nil {
		return err
	} else {
		return nil
	}
}

// Delete a table and all its contents
func (z *ZetabaseClient) DeleteTable(tableOwnerId, tableId string) error {
	if !z.checkReady() {
		return errors.New("NotReady")
	}
	nonce := z.nonceMaker.Get()
	extraBytes := []byte(tableId)
	poc := z.getCredential(nonce, extraBytes)

	_, err := z.client.DeleteObject(z.ctx, &zbprotocol.DeleteSystemObjectRequest{
		Id:           z.userId,
		ObjectType:   zbprotocol.SystemObjectType_TABLE,
		TableOwnerId: tableOwnerId,
		TableId:      tableId,
		ObjectId:     tableId,
		Nonce:        nonce,
		Credential:   poc,
	})

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (z *ZetabaseClient) listKeysWithPattern(tableOwnerId, tableId, pattern string, pgIdx int64) ([]string, bool, error) {
	if !z.checkReady() {
		return nil, false, errors.New("NotReady")
	}
	nonce := z.nonceMaker.Get()
	poc := z.getCredential(nonce, nil)
	res, err := z.client.ListKeys(z.ctx, &zbprotocol.ListKeysRequest{
		Id:           z.userId,
		TableId:      tableId,
		TableOwnerId: tableOwnerId,
		Pattern:      pattern,
		Nonce:        nonce,
		PageIndex:    pgIdx,
		Credential:   poc,
	})
	if err != nil {
		return nil, false, err
	} else {
		rig := res.GetKeys()
		return rig, res.GetPagination().GetHasNextPage(), nil
	}
}

// Return the underlying gRPC client
func (z *ZetabaseClient) GrpcClient() zbprotocol.ZetabaseProviderClient {
	return z.client
}

//func (z *ZetabaseClient) UserPasswordAuthenticate(userId, passw string) error {
//	if z.privKey != nil && z.pubKey != nil {
//		// Just to catch weird cases
//		return errors.New("NonNilKeysAlreadyPassedIn")
//	}
//	res, err := z.client.LoginUser(z.ctx, &AuthenticateUser{
//		Handle: userId,
//		Password: passw,
//		Nonce: int64(0),
//		Credential: MakeEmptyCredentials(),
//	})
//	if err != nil {
//		return err
//	} else {
//		jwt := res.GetJwtToken()
//		uid := res.GetId()
//		if len(uid) > 0 && len(jwt) > 0 {
//			z.userId = uid
//		} else {
//			return errors.New("NoSuchId")
//		}
//		return nil
//	}
//}


const (
	DoSaveCertificates = false
)
// Connect to Zetabase with the provided credentials.
func (z *ZetabaseClient) Connect() error {
	if z.insecure {
		wi := grpc.WithInsecure()
		conn, err := grpc.Dial(z.serverAddr, wi)
		if err != nil {
			return err
		}
		z.conn = conn
	} else {
		if z.noCertVerify {
			configT := &tls.Config{InsecureSkipVerify: true}
			wtc := grpc.WithTransportCredentials(credentials.NewTLS(configT))
			conn, err := grpc.Dial(z.serverAddr, wtc)
			if err != nil {
				return err
			}
			z.conn = conn
		} else {
			configT := &tls.Config{InsecureSkipVerify: false}
			wtc := grpc.WithTransportCredentials(credentials.NewTLS(configT))
			conn, err := grpc.Dial(z.serverAddr, wtc)
			if err != nil {
				return err
			}
			z.conn = conn
		}
	}
	z.client = zbprotocol.NewZetabaseProviderClient(z.conn)
	return nil
}
