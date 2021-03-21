package function

// CS 161 Project 2 Spring 2020
// You MUST NOT change what you import.  If you add ANY additional
// imports it will break the autograder. We will be very upset.

import (
	// You neet to add with
	// go get github.com/cs161-staff/userlib
	"github.com/cs161-staff/userlib"

	// Life is much easier with json:  You are
	// going to want to use this so you can easily
	// turn complex structures into strings etc...
	"encoding/json"

	// Likewise useful for debugging, etc...
	"encoding/hex"

	// UUIDs are generated right based on the cryptographic PRNG
	// so lets make life easier and use those too...
	//
	// You need to add with "go get github.com/google/uuid"
	"github.com/google/uuid"

	// Useful for debug messages, or string manipulation for datastore keys.
	"strings"

	// Want to import errors.
	"errors"

	// Optional. You can remove the "_" there, but please do not touch
	// anything else within the import bracket.
	_ "strconv"
	// if you are looking for fmt, we don't give you fmt, but you can use userlib.DebugMsg.
	// see someUsefulThings() below:
)

// This serves two purposes:
// a) It shows you some useful primitives, and
// b) it suppresses warnings for items not being imported.
// Of course, this function can be deleted.
func someUsefulThings() {
	// Creates a random UUID
	f := uuid.New()
	userlib.DebugMsg("UUID as string:%v", f.String())

	// Example of writing over a byte of f
	f[0] = 10
	userlib.DebugMsg("UUID as string:%v", f.String())

	// takes a sequence of bytes and renders as hex
	h := hex.EncodeToString([]byte("fubar"))
	userlib.DebugMsg("The hex: %v", h)

	// Marshals data into a JSON representation
	// Will actually work with go structures as well
	d, _ := json.Marshal(f)
	userlib.DebugMsg("The json data: %v", string(d))
	var g uuid.UUID
	json.Unmarshal(d, &g)
	userlib.DebugMsg("Unmashaled data %v", g.String())

	// This creates an error type
	userlib.DebugMsg("Creation of error %v", errors.New(strings.ToTitle("This is an error")))

	// And a random RSA key.  In this case, ignoring the error
	// return value
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("Key is %v, %v", pk, sk)
}

// Helper function: Takes the first 16 bytes and
// converts it into the UUID type
func bytesToUUID(data []byte) (ret uuid.UUID) {
	for x := range ret {
		ret[x] = data[x]
	}
	return
}

// The structure definition for a user record
type User struct {
	Username string
	Key      []byte
	Pk       userlib.PKEDecKey
	Dss      userlib.DSSignKey
	// You can add other fields here if you want...
	// Note for JSON to marshal/unmarshal, the fields need to
	// be public (start with a capital letter)
}
type UserId struct {
	Salt    []byte
	HmacVal []byte
	Hashed  [64]byte
}
type EncMac struct {
	Enc []byte
	Mac []byte
}
type LogRec struct {
	FileKey  []byte
	NrBlocks int
}
type PseudoRec struct {
	EncKey []byte
	MacKey []byte
	Log    uuid.UUID
}
type TokenRec struct {
	EncKey []byte
	MacKey []byte
	Pseudo uuid.UUID
}

func EncDsData(data []byte, PkeEncKey userlib.PKEEncKey, DsKey userlib.DSSignKey) (result []byte) {
	encData, _ := userlib.PKEEnc(PkeEncKey, data)
	dsEncData, _ := userlib.DSSign(DsKey, encData)
	encMacPair := EncMac{encData, dsEncData}
	result, _ = json.Marshal(encMacPair)
	return result
}
func DecDsData(data []byte, PkeDecKey userlib.PKEDecKey, VKey userlib.DSVerifyKey) (result []byte, ok error) {
	var encMacPair EncMac
	json.Unmarshal(data, &encMacPair)
	err := userlib.DSVerify(VKey, encMacPair.Enc, encMacPair.Mac)
	if err != nil {
		userlib.DebugMsg("Error cannot verify signiture")
		return result, err
	}
	result, err1 := userlib.PKEDec(PkeDecKey, encMacPair.Enc)
	if err1 != nil {
		userlib.DebugMsg("Error cannot decrypt the data")
		return result, err1
	}
	return result, ok
}

func EncMacData(data []byte, EncKey []byte, MacKey []byte) (result []byte) {
	iv := userlib.RandomBytes(16)
	encData := userlib.SymEnc(EncKey, iv, data)
	macEncData, _ := userlib.HMACEval(MacKey, encData)
	encMacPair := EncMac{encData, macEncData}
	result, _ = json.Marshal(encMacPair)
	return result
}
func DecMacData(data []byte, EncKey []byte, MacKey []byte) (result []byte, ok bool) {
	//ok is flase if the data has been modified
	var encMacPair EncMac
	json.Unmarshal(data, &encMacPair)
	macVal, _ := userlib.HMACEval(MacKey, encMacPair.Enc)
	ok = true
	if userlib.HMACEqual(macVal, encMacPair.Mac) == false {
		ok = false
		return result, ok
	}
	result = userlib.SymDec(EncKey, encMacPair.Enc)
	return result, ok
}
func FileNameUUIDSaver(fileKey []byte, i int, data []byte) {
	blockID, _ := userlib.HashKDF(fileKey, []byte(string(i)))
	blockUUID := bytesToUUID(blockID[:16])
	/*if i == 0 {
		userlib.DebugMsg("Saver FileKey %v", fileKey)
		userlib.DebugMsg("Saver UUID %v", blockUUID)
	}*/
	userlib.DatastoreSet(blockUUID, data)
}
func FileNameUUIDLoader(fileKey []byte, i int) (data []byte, ok bool) {
	blockID, _ := userlib.HashKDF(fileKey, []byte(string(i)))
	blockUUID := bytesToUUID(blockID[:16])
	/*if i == 0 {
		userlib.DebugMsg("Loader FileKey %v", fileKey)
		userlib.DebugMsg("Loader UUID %v", blockUUID)
	}*/
	data, ok = userlib.DatastoreGet(blockUUID)
	return data, ok
}
func KeyGen(key []byte) (macKey, encriptionKey []byte) {
	hmacKey, _ := userlib.HashKDF(key, []byte("mac"))
	hmacKey = hmacKey[:16]
	encryptionKey, _ := userlib.HashKDF(key, []byte("encryption"))
	encryptionKey = encryptionKey[:16]
	return hmacKey, encryptionKey
}
func TokenFinder(username string, filename string, key []byte, encryptionKey []byte, macKey []byte) (result TokenRec, flag bool) {
	tokenKey, _ := userlib.HashKDF(key, []byte("token"+filename+username))
	tokenRecUUID := bytesToUUID(tokenKey[:16])
	tokenRecEnc, ok := userlib.DatastoreGet(tokenRecUUID)
	flag = true
	if !ok {
		flag = false
		return result, flag
	}
	//token exists decrypt the data in token and reset ur enc and mac keys
	var tokenRec TokenRec
	tokenData, ok1 := DecMacData(tokenRecEnc, encryptionKey, macKey)
	if !ok1 {
		flag = false
		return result, flag
	}
	json.Unmarshal(tokenData, &tokenRec)
	result = tokenRec
	return result, flag
}
func PseudoFinder(pseudoRecUUID uuid.UUID, encryptionKey []byte, macKey []byte) (result PseudoRec, flag bool) {
	pseudoRecEnc, ok2 := userlib.DatastoreGet(pseudoRecUUID)
	flag = true
	if !ok2 {
		flag = false
		return result, flag
	}
	pseudoData, ok3 := DecMacData(pseudoRecEnc, encryptionKey, macKey)
	if !ok3 {
		flag = false
		return result, flag
	}

	json.Unmarshal(pseudoData, &result)
	return result, flag
}
func LogFinder(logRecUUID uuid.UUID, encryptionKey []byte, macKey []byte) (result LogRec, flag bool) {
	logRecEnc, ok4 := userlib.DatastoreGet(logRecUUID)
	flag = true
	if !ok4 {
		flag = false
		return result, flag
	}
	logData, ok5 := DecMacData(logRecEnc, encryptionKey, macKey)
	if !ok5 {
		flag = false
		return result, flag
	}
	var logRec LogRec
	json.Unmarshal(logData, &logRec)
	result = logRec
	return result, flag
}

// This creates a user.  It will only be called once for a user
// (unless the keystore and datastore are cleared during testing purposes)

// It should store a copy of the userdata, suitably encrypted, in the
// datastore and should store the user's public key in the keystore.

// The datastore may corrupt or completely erase the stored
// information, but nobody outside should be able to get at the stored

// You are not allowed to use any global storage other than the
// keystore and the datastore functions in the userlib library.

// You can assume the password has strong entropy, EXCEPT
// the attackers may possess a precomputed tables containing
// hashes of common passwords downloaded from the internet.
func InitUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	userdataptr = &userdata

	//TODO: This is a toy implementation.
	// dont have the key come only from the pass
	if len(username) == 0 {
		return nil, errors.New(strings.ToTitle("Username needs to be non empty"))
	}
	userdata.Username = username
	clientName, _ := json.Marshal(username)
	clientNameHash := userlib.Hash(clientName)
	nameJson, _ := json.Marshal(clientNameHash)
	nameJson = nameJson[:16]
	userkey, _ := userlib.HashKDF(nameJson, []byte("userkey"))
	userkey = userkey[:16]
	userDataUUID := bytesToUUID(userkey)
	_, ok := userlib.DatastoreGet(userDataUUID)
	if ok {
		return nil, errors.New(strings.ToTitle("User with same username already exists"))
	}

	client, _ := json.Marshal(password + username)
	clientHash := userlib.Hash(client)
	passJson, _ := json.Marshal(clientHash)
	passJson = passJson[:16]
	//usernameJson, _ := json.Marshal(username)
	passUsernameJson, _ := json.Marshal(password + username)
	//generate public key and store it under username
	sk, pk, _ := userlib.PKEKeyGen()
	userdata.Pk = pk
	dss, dsv, _ := userlib.DSKeyGen()
	userdata.Dss = dss
	userlib.KeystoreSet(username, sk)
	userlib.KeystoreSet(username+"ds", dsv)
	//creation of keys and salt, hashing of password
	salt := userlib.RandomBytes(16)
	storeData := passUsernameJson
	hashed := userlib.Hash(storeData)
	hashedJson, _ := json.Marshal(hashed)
	key := userlib.Argon2Key(passJson, salt, 16)
	userdata.Key = key
	//generate mac and enc keys
	macKey, encryptionKey := KeyGen(key)
	hmacVal, _ := userlib.HMACEval(macKey, hashedJson)
	userID := UserId{salt, hmacVal, hashed}
	byteUserID, _ := json.Marshal(userID)
	passkey, _ := userlib.HashKDF(passJson, []byte("passkey"))
	passkey = passkey[:16]
	userUUID := bytesToUUID(passkey)
	userlib.DatastoreSet(userUUID, byteUserID)
	//where the data is enc and mac
	byteUserData, _ := json.Marshal(userdata)
	encUserData := userlib.SymEnc(encryptionKey, salt, byteUserData)
	macEncData, _ := userlib.HMACEval(macKey, encUserData)
	encMacPair := EncMac{encUserData, macEncData}
	byteEncMacPair, _ := json.Marshal(encMacPair)
	//line u need to change is below

	userlib.DatastoreSet(userDataUUID, byteEncMacPair)

	//End of toy implementation

	return &userdata, nil
}

// This fetches the user information from the Datastore.  It should
// fail with an error if the user/password is invalid, or if the user
// data was corrupted, or if the user can't be found.
func GetUser(username string, password string) (userdataptr *User, err error) {
	var userdata User
	userdataptr = &userdata
	//identity verification
	if len(username) == 0 {
		return nil, errors.New(strings.ToTitle("Username needs to be non empty"))
	}
	clientName, _ := json.Marshal(username)
	clientNameHash := userlib.Hash(clientName)
	nameJson, _ := json.Marshal(clientNameHash)
	nameJson = nameJson[:16]
	client, _ := json.Marshal(password + username)
	clientHash := userlib.Hash(client)
	passJson, _ := json.Marshal(clientHash)
	passJson = passJson[:16]
	//usernameJson, _ := json.Marshal(username)
	passUsernameJson, _ := json.Marshal(password + username)
	passkey, err := userlib.HashKDF(passJson, []byte("passkey"))
	passkey = passkey[:16]
	userUUID := bytesToUUID(passkey)
	userID, ok := userlib.DatastoreGet(userUUID)
	//datastore gives error user not found
	if ok == false {
		err = errors.New(strings.ToTitle("User not found/ Or wrong credentials"))
		userlib.DebugMsg("Creation of error %v", err)
		return nil, err
	}
	var tempUser UserId
	json.Unmarshal(userID, &tempUser)
	key := userlib.Argon2Key(passJson, tempUser.Salt, 16)
	storeData := passUsernameJson
	hashed := userlib.Hash(storeData)
	hashedJson, _ := json.Marshal(hashed)
	macKey, encryptionKey := KeyGen(key)
	hmacValNew, _ := userlib.HMACEval(macKey, hashedJson)
	//hashing and mac the storeData to compare macs. Error if macs !=
	if userlib.HMACEqual(tempUser.HmacVal, hmacValNew) == false {
		err = errors.New(strings.ToTitle("Data has been modified"))
		userlib.DebugMsg("Creation of error %v", err)
		return nil, err
	}
	//Compares the hashes
	if hashed != tempUser.Hashed {
		err = errors.New(strings.ToTitle("Wrong user/password"))
		userlib.DebugMsg("Creation of error %v", err)
		return nil, err
	}
	//Collects the userdata
	userkey, err := userlib.HashKDF(nameJson, []byte("userkey"))
	userkey = userkey[:16]
	userDataUUID := bytesToUUID(userkey)
	byteEncMacPair, ok := userlib.DatastoreGet(userDataUUID)
	var encMacPair EncMac
	json.Unmarshal(byteEncMacPair, &encMacPair)
	hmacValNew2, _ := userlib.HMACEval(macKey, encMacPair.Enc)
	if userlib.HMACEqual(encMacPair.Mac, hmacValNew2) == false {
		err = errors.New(strings.ToTitle("Data has been modified"))
		userlib.DebugMsg("Creation of error %v", err)
		return nil, err
	}
	byteUserData := userlib.SymDec(encryptionKey, encMacPair.Enc)
	json.Unmarshal(byteUserData, userdataptr)
	return userdataptr, nil
}

// This stores a file in the datastore.
//
// The plaintext of the filename + the plaintext and length of the filename
// should NOT be revealed to the datastore!
func (userdata *User) StoreFile(filename string, data []byte) {

	//check if token record exist first
	key := userdata.Key
	tokenKey, _ := userlib.HashKDF(key, []byte("token"+filename+userdata.Username))
	tokenRecUUID := bytesToUUID(tokenKey[:16])
	//if userdata.Username == "welp"{
	//	userlib.DebugMsg("The token saved %v", tokenRecUUID)
	//}
	//add pseudo
	_, ok := userlib.DatastoreGet(tokenRecUUID)
	//token record does not exist first time file is being saved
	if ok == false {

		//generating mac/enc key and filename key which will be used for uuid generator
		fileKey, _ := userlib.HashKDF(key, []byte(filename+userdata.Username))
		fileKey = fileKey[:16]
		macKey, encryptionKey := KeyGen(key)
		var i int
		//split in blocks enc each block and save them to file
		for i = 0; len(data) >= 16777216; i++ {
			block := data[:16777216]
			data = data[16777216:]
			blockEncMacByte := EncMacData(block, encryptionKey, macKey)
			FileNameUUIDSaver(fileKey, i, blockEncMacByte)
		}
		//the last block with less than 16 char more than 0
		if len(data) > 0 {
			block := data
			blockEncMacByte := EncMacData(block, encryptionKey, macKey)
			FileNameUUIDSaver(fileKey, i, blockEncMacByte)
			i = i + 1
		}
		//setting up log record
		var logRec LogRec
		logRecKey, _ := userlib.HashKDF(key, []byte("logRec"+filename+userdata.Username))
		logRecKey = logRecKey[:16]
		logRecUUID := bytesToUUID(logRecKey)
		logRec.NrBlocks = i
		logRec.FileKey = fileKey
		logRecByte, _ := json.Marshal(logRec)
		logRecEnc := EncMacData(logRecByte, encryptionKey, macKey)
		userlib.DatastoreSet(logRecUUID, logRecEnc)
		//finished setting up log record
		//setting up pseudo record
		var pseudoRec PseudoRec
		pseudoRec.EncKey = encryptionKey
		pseudoRec.MacKey = macKey
		pseudoRec.Log = logRecUUID
		pseudoData, _ := json.Marshal(pseudoRec) //logrecuuid
		pseudoRecEnc := EncMacData(pseudoData, encryptionKey, macKey)
		pseudoRecKey, _ := userlib.HashKDF(key, []byte("pseudo"+filename+userdata.Username))
		pseudoRecKey = pseudoRecKey[:16]
		pseudoRecUUID := bytesToUUID(pseudoRecKey)
		userlib.DatastoreSet(pseudoRecUUID, pseudoRecEnc)
		//finished setting up pseudo record
		//setting up token record
		var tokenRec TokenRec
		tokenRec.EncKey = encryptionKey
		tokenRec.MacKey = macKey
		tokenRec.Pseudo = pseudoRecUUID
		tokenData, _ := json.Marshal(tokenRec)
		tokenRecEnc := EncMacData(tokenData, encryptionKey, macKey)
		userlib.DatastoreSet(tokenRecUUID, tokenRecEnc)
		//if userdata.Username == "mone"{
		//	userlib.DebugMsg("we dont want 2 be here 2  %v", fileKey)
		//}
		//finished setting up token rec
	} else {
		key := userdata.Key
		macKey, encryptionKey := KeyGen(key)
		//keys generated
		//check if token exists throw error if it doesnt
		tokenRec, ok := TokenFinder(userdata.Username, filename, key, encryptionKey, macKey)
		if !ok {
			userlib.DebugMsg("Token not found or Hmac not matching")
		}
		//resetting the keys
		macKey = tokenRec.MacKey
		encryptionKey = tokenRec.EncKey
		pseudoRecUUID := tokenRec.Pseudo
		//finished with getting the token
		//start getting pseudo
		pseudoRec, ok1 := PseudoFinder(pseudoRecUUID, encryptionKey, macKey)
		if !ok1 {
			userlib.DebugMsg("Pseudo not found or Hmac not matching")
		}
		//got pseudo record and log record uuid
		//retrieve log rec
		//reset the keys
		macKey = pseudoRec.MacKey
		encryptionKey = pseudoRec.EncKey
		logRecUUID := pseudoRec.Log
		logRec, ok2 := LogFinder(logRecUUID, encryptionKey, macKey)
		if !ok2 {
			userlib.DebugMsg("Log not found or Hmac not matching")
		}
		//found log rec
		//start writing data from begining
		fileKey := logRec.FileKey
		var i int
		for i = 0; len(data) >= 16777216; i++ {
			block := data[:16777216]
			data = data[16777216:]
			blockEncMacByte := EncMacData(block, encryptionKey, macKey)
			FileNameUUIDSaver(fileKey, i, blockEncMacByte)
		}
		//the last block with less than 16 char more than 0
		if len(data) > 0 {
			block := data
			blockEncMacByte := EncMacData(block, encryptionKey, macKey)
			FileNameUUIDSaver(fileKey, i, blockEncMacByte)
			i = i + 1
		}
		//update log rec
		logRec.NrBlocks = i
		logRecByte, _ := json.Marshal(logRec)
		logRecEnc := EncMacData(logRecByte, encryptionKey, macKey)
		userlib.DatastoreSet(logRecUUID, logRecEnc)
		//if userdata.Username == "mone"{
		//	userlib.DebugMsg("here we are   %v", fileKey)
		//}
	}
	
	return
}

// This adds on to an existing file.
//
// Append should be efficient, you shouldn't rewrite or reencrypt the
// existing file, but only whatever additional information and
// metadata you need.
func (userdata *User) AppendFile(filename string, data []byte) (err error) {
	key := userdata.Key
	macKey, encryptionKey := KeyGen(key)
	//keys generated
	//check if token exists throw error if it doesnt
	tokenRec, ok := TokenFinder(userdata.Username, filename, key, encryptionKey, macKey)
	if !ok {
		return errors.New(strings.ToTitle("Token not found or Hmac not matching"))
	}
	//resetting the keys
	macKey = tokenRec.MacKey
	encryptionKey = tokenRec.EncKey
	pseudoRecUUID := tokenRec.Pseudo
	//finished with getting the token
	//start getting pseudo
	pseudoRec, ok1 := PseudoFinder(pseudoRecUUID, encryptionKey, macKey)
	if !ok1 {
		return errors.New(strings.ToTitle("Pseudo not found or Hmac not matching"))
	}
	//got pseudo record and log record uuid
	//retrieve log rec
	//reset the keys
	macKey = pseudoRec.MacKey
	encryptionKey = pseudoRec.EncKey
	logRecUUID := pseudoRec.Log
	logRec, ok2 := LogFinder(logRecUUID, encryptionKey, macKey)
	if !ok2 {
		return errors.New(strings.ToTitle("Log not found or Hmac not matching"))
	}
	var i int
	//encrypting the new data
	//split in blocks enc each block and save them to file
	for i = logRec.NrBlocks; len(data) >= 16777216; i++ {
		block := data[:16777216]
		data = data[16777216:]
		blockEncMacByte := EncMacData(block, encryptionKey, macKey)
		FileNameUUIDSaver(logRec.FileKey, i, blockEncMacByte)
	}
	//the last block with less than 16 char more than 0
	if len(data) > 0 {
		block := data
		blockEncMacByte := EncMacData(block, encryptionKey, macKey)
		FileNameUUIDSaver(logRec.FileKey, i, blockEncMacByte)
		i = i + 1
	}
	//updating the logRec
	logRec.NrBlocks = i
	logRecByte, _ := json.Marshal(logRec)
	logRecEnc := EncMacData(logRecByte, encryptionKey, macKey)
	userlib.DatastoreSet(logRecUUID, logRecEnc)
	return
}

// This loads a file from the Datastore.
//
// It should give an error if the file is corrupted in any way.
func (userdata *User) LoadFile(filename string) (data []byte, err error) {
	//TODO: This is a toy implementation.
	//generate keys to decrypt token
	key := userdata.Key
	macKey, encryptionKey := KeyGen(key)
	//keys generated
	//check if token exists throw error if it doesnt
	tokenRec, ok := TokenFinder(userdata.Username, filename, key, encryptionKey, macKey)
	if !ok {
		return nil, errors.New(strings.ToTitle("Token not found or Hmac not matching"))
	}
	//resetting the keys
	macKey = tokenRec.MacKey
	encryptionKey = tokenRec.EncKey
	pseudoRecUUID := tokenRec.Pseudo
	//finished with getting the token
	//start getting pseudo
	pseudoRec, ok1 := PseudoFinder(pseudoRecUUID, encryptionKey, macKey)
	if !ok1 {
		return nil, errors.New(strings.ToTitle("Pseudo not found or Hmac not matching"))
	}
	//got pseudo record and log record uuid
	//retrieve log rec
	//reset the keys
	macKey = pseudoRec.MacKey
	encryptionKey = pseudoRec.EncKey
	logRecUUID := pseudoRec.Log
	logRec, ok2 := LogFinder(logRecUUID, encryptionKey, macKey)
	if !ok2 {
		return nil, errors.New(strings.ToTitle("Log not found or Hmac not matching"))
	}
	//log rec data found
	//loading data
	var dataResult []byte
	//its 0 index so we use <
	for i := 0; i < logRec.NrBlocks; i++ {
		//first loading the data
		blockEncMacByte, ok6 := FileNameUUIDLoader(logRec.FileKey, i)
		if !ok6 {
			userlib.DebugMsg("Block %d not found", i)
			return nil, errors.New(strings.ToTitle("Block not found"))
		}
		//dec and checking mac of data
		blockData, ok7 := DecMacData(blockEncMacByte, encryptionKey, macKey)
		if !ok7 {
			userlib.DebugMsg("Block %d HMac not matching", i)
			return nil, errors.New(strings.ToTitle("Block Hmac not matching"))
		}
		dataResult = append(dataResult, blockData...)

	}
	//data found
	return dataResult, nil
	//End of toy implementation
	//return
}

// This creates a sharing record, which is a key pointing to something
// in the datastore to share with the recipient.

// This enables the recipient to access the encrypted file as well
// for reading/appending.

// Note that neither the recipient NOR the datastore should gain any
// information about what the sender calls the file.  Only the
// recipient can access the sharing record, and only the recipient
// should be able to know the sender.
func (userdata *User) ShareFile(filename string, recipient string) (
	magic_string string, err error) {

	key := userdata.Key
	macKey, encryptionKey := KeyGen(key)
	//add recipient to list of people file is shared with
	listKey, _ := userlib.HashKDF(key, []byte("list"+filename+userdata.Username))
	listUUID := bytesToUUID(listKey[:16])
	listEnc, ok := userlib.DatastoreGet(listUUID)
	if !ok {
		//list doesnt exist have to make a new one
		m := make(map[string]int)
		m[recipient] = 1
		listByte, _ := json.Marshal(m)
		listEncMac := EncMacData(listByte, encryptionKey, macKey)
		userlib.DatastoreSet(listUUID, listEncMac)
	} else {
		//list exist we need to load it add value and restore it the data store
		m := make(map[string]int)
		listData, ok1 := DecMacData(listEnc, encryptionKey, macKey)
		if !ok1 {
			return magic_string, errors.New(strings.ToTitle("List Hmac not matching"))
		}
		json.Unmarshal(listData, &m)
		m[recipient] = 1
		listByte, _ := json.Marshal(m)
		listEncMac := EncMacData(listByte, encryptionKey, macKey)
		userlib.DatastoreSet(listUUID, listEncMac)
	}
	//keys generated
	//check if token exists throw error if it doesnt
	tokenRec, ok2 := TokenFinder(userdata.Username, filename, key, encryptionKey, macKey)
	if !ok2 {
		return magic_string, errors.New(strings.ToTitle("Token not found or Hmac not matching"))
	}
	//resetting the keys
	macKey = tokenRec.MacKey
	encryptionKey = tokenRec.EncKey
	pseudoRecUUID := tokenRec.Pseudo
	//getting the pseudo rec data

	pseudoRec, ok3 := PseudoFinder(pseudoRecUUID, encryptionKey, macKey)
	if !ok3 {
		return magic_string, errors.New(strings.ToTitle("Pseudo not found or Hmac not matching"))
	}
	//creating pseudo rec for recipient
	recipientEncryptionKey, _ := userlib.HashKDF(key, []byte("list"+filename+recipient+"enc"))
	recipientEncryptionKey = recipientEncryptionKey[:16]
	recipientMacKey, _ := userlib.HashKDF(key, []byte("list"+filename+recipient+"mac"))
	recipientMacKey = recipientMacKey[:16]
	recipientPseudoKey, _ := userlib.HashKDF(key, []byte("list"+filename+recipient))
	recipientPseudoUUID := bytesToUUID(recipientPseudoKey[:16])
	recipientPseudoByte, _ := json.Marshal(pseudoRec)
	recipientPseudocEnc := EncMacData(recipientPseudoByte, recipientEncryptionKey, recipientMacKey)
	userlib.DatastoreSet(recipientPseudoUUID, recipientPseudocEnc)
	//creating token for rec
	//using his encryption key to encrypt and sign with our key
	tokenRec.Pseudo = recipientPseudoUUID
	tokenRec.EncKey = recipientEncryptionKey
	tokenRec.MacKey = recipientMacKey
	recipientKey, ok4 := userlib.KeystoreGet(recipient)
	if !ok4 {
		return magic_string, errors.New(strings.ToTitle("User has not set up account DS key not found"))
	}
	recipientTokenByte, _ := json.Marshal(tokenRec)
	magic_string = string(EncDsData(recipientTokenByte, recipientKey, userdata.Dss))

	return magic_string, nil
}

// Note recipient's filename can be different from the sender's filename.
// The recipient should not be able to discover the sender's view on
// what the filename even is!  However, the recipient must ensure that
// it is authentically from the sender.
func (userdata *User) ReceiveFile(filename string, sender string, magic_string string) error {
	senderKey, ok := userlib.KeystoreGet(sender + "ds")
	if !ok {
		return errors.New(strings.ToTitle("User has not set up account DS key not found"))
	}
	tokenByte, ok1 := DecDsData([]byte(magic_string), userdata.Pk, senderKey)
	if ok1 != nil {
		//return ok1
		return errors.New(strings.ToTitle("Error with signiture of Dec of token"))
	}
	key := userdata.Key
	macKey, encryptionKey := KeyGen(key)
	tokenKey, _ := userlib.HashKDF(key, []byte("token"+filename+userdata.Username))
	tokenRecUUID := bytesToUUID(tokenKey[:16])
	_, ok = userlib.DatastoreGet(tokenRecUUID)
	//check if user already has a file with that filename
	if ok {
		return errors.New(strings.ToTitle("Error with filename. Already existing filename for user"))
	}
	tokenRecEnc := EncMacData(tokenByte, encryptionKey, macKey)
	userlib.DatastoreSet(tokenRecUUID, tokenRecEnc)
	return nil
}

// Removes target user's access.
func (userdata *User) RevokeFile(filename string, target_username string) (err error) {
	key := userdata.Key
	macKey, encryptionKey := KeyGen(key)
	//add recipient to list of people file is shared with
	listKey, _ := userlib.HashKDF(key, []byte("list"+filename+userdata.Username))
	listUUID := bytesToUUID(listKey[:16])
	listEncMac, ok := userlib.DatastoreGet(listUUID)
	if !ok {
		//no files shared nothing to revoke
		return errors.New(strings.ToTitle("User share list not found user hasn't shared"))
	}
	pseudoRecKey, _ := userlib.HashKDF(key, []byte("pseudo"+filename+userdata.Username))
	pseudoRecKey = pseudoRecKey[:16]
	pseudoRecUUID := bytesToUUID(pseudoRecKey)
	pseudoRecEnc, ok := userlib.DatastoreGet(pseudoRecUUID)
	if !ok {
		//no revoking privileges
		return errors.New(strings.ToTitle("User does not have revoking privileges"))
	}

	listData, ok := DecMacData(listEncMac, encryptionKey, macKey)
	if !ok {
		//Enc/Mac failed
		return errors.New(strings.ToTitle("User list has been modiefied or cannot be encrypted"))
	}
	m := make(map[string]int)
	json.Unmarshal(listData, &m)
	_, present := m[target_username]
	if !present {
		//Target not in the list
		return errors.New(strings.ToTitle("Target not in the list"))
	}
	delete(m, target_username)
	//userlib.DebugMsg("User %v gets access", m)

	listByte, _ := json.Marshal(m)
	listEncMac = EncMacData(listByte, encryptionKey, macKey)
	userlib.DatastoreSet(listUUID, listEncMac)
	//list has been updated
	//data needs to be loaded
	data, err := userdata.LoadFile(filename)
	if err != nil {
		return err
	}
	//get token
	tokenRec, ok := TokenFinder(userdata.Username, filename, key, encryptionKey, macKey)
	if !ok {
		return errors.New(strings.ToTitle("Token not found or modified"))
	}
	macKey = tokenRec.MacKey
	encryptionKey = tokenRec.EncKey
	//finished with getting the token
	//start getting pseudo
	pseudoRec, ok1 := PseudoFinder(pseudoRecUUID, encryptionKey, macKey)
	oldlogUUID := pseudoRec.Log
	userlib.DatastoreDelete(oldlogUUID)
	if !ok1 {
		return errors.New(strings.ToTitle("Pseudo not found or Hmac not matching"))
	}
	//got pseudo
	//create new key for generating keys
	tempKey, _ := userlib.HashKDF(key, []byte(filename+userdata.Username+target_username))
	tempKey = tempKey[:16]
	//key generated
	//enc and mac data with the new pair
	fileKey, _ := userlib.HashKDF(tempKey, []byte(filename+userdata.Username+target_username))
	fileKey = fileKey[:16]
	macNewKey, encryptionNewKey := KeyGen(tempKey)
	var i int
	//split in blocks enc each block and save them to file
	for i = 0; len(data) >= 16777216; i++ {
		block := data[:16777216]
		data = data[16777216:]
		blockEncMacByte := EncMacData(block, encryptionNewKey, macNewKey)
		FileNameUUIDSaver(fileKey, i, blockEncMacByte)
	}
	//the last block with less than 16 char more than 0
	if len(data) > 0 {
		block := data
		blockEncMacByte := EncMacData(block, encryptionNewKey, macNewKey)
		FileNameUUIDSaver(fileKey, i, blockEncMacByte)
		i = i + 1
	}
	//data enc and mac
	//save new log rec in data store enc and mac
	var logRec LogRec
	logRecKey, _ := userlib.HashKDF(tempKey, []byte("logRec"+filename+userdata.Username))
	logRecKey = logRecKey[:16]
	logRecUUID := bytesToUUID(logRecKey)
	logRec.NrBlocks = i
	logRec.FileKey = fileKey
	logRecByte, _ := json.Marshal(logRec)
	logRecEnc := EncMacData(logRecByte, encryptionNewKey, macNewKey)
	userlib.DatastoreSet(logRecUUID, logRecEnc)
	//update pseudo rec
	pseudoRec.EncKey = encryptionNewKey
	pseudoRec.MacKey = macNewKey
	pseudoRec.Log = logRecUUID
	pseudoData, _ := json.Marshal(pseudoRec)
	pseudoRecEnc = EncMacData(pseudoData, encryptionKey, macKey)
	userlib.DatastoreSet(pseudoRecUUID, pseudoRecEnc)
	//owner pseudorec has been updated
	for recipient, _ := range m {
		userlib.DebugMsg("User %v gets access", recipient)
		recipientEncryptionKey, _ := userlib.HashKDF(key, []byte("list"+filename+recipient+"enc"))
		recipientEncryptionKey = recipientEncryptionKey[:16]
		recipientMacKey, _ := userlib.HashKDF(key, []byte("list"+filename+recipient+"mac"))
		recipientMacKey = recipientMacKey[:16]
		recipientPseudoKey, _ := userlib.HashKDF(key, []byte("list"+filename+recipient))
		recipientPseudoUUID := bytesToUUID(recipientPseudoKey[:16])
		recipientPseudoByte, _ := json.Marshal(pseudoRec)
		recipientPseudocEnc := EncMacData(recipientPseudoByte, recipientEncryptionKey, recipientMacKey)
		userlib.DatastoreSet(recipientPseudoUUID, recipientPseudocEnc)
	}

	return
}
