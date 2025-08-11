package client

// CS 161 Project 2

// may break the autograder!

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string


	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	//

	_ = fmt.Sprintf("%s_%d", "file", 1)
}

type User struct {
	Username string
	MasterKey []byte
	RSAPrivateKey userlib.PKEDecKey
	DigitalSigningKey userlib.DSSignKey
}

type FileInfo struct {
	Owner bool
	FileHeaderUUID uuid.UUID
	AccessPointUUID uuid.UUID
	EncKey []byte
	MACKey []byte
	OwnerAccessInfoMapUUID uuid.UUID
}

type FileHeader struct {
	FirstUUID uuid.UUID
	NewUUID uuid.UUID
}

type File struct {
	Content []byte
	NextUUID uuid.UUID
}

type Invitation struct {
	AccessPointUUID uuid.UUID
	APEncKey []byte
	APMACKey []byte
}

type AccessPoint struct {
	EncKey []byte
	MACKey []byte
	FileHeaderUUID uuid.UUID
	Revoked bool
}

type OwnerAccessInfoMap struct {
	InfoMap map[string]OwnerAccessInfo
}

type OwnerAccessInfo struct {
	AccessPointUUID uuid.UUID
	APEncKey []byte
	APMACKey []byte
}


func InitUser(username string, password string) (userdataptr *User, err error) {
	var userdata User



	if (username == "") {
		return nil, errors.New("empty username")
	}
	userdata.Username = username

	LoginUUID, err := uuid.FromBytes(userlib.Hash([]byte(username))[:16])
	if err != nil {
		return nil, err
	}

	_, ok := userlib.DatastoreGet(LoginUUID)
	
	if ok {
		return nil, errors.New("error user already exists")
	}

	IV := userlib.RandomBytes(16)
	MasterKey := userlib.Argon2Key([]byte(password), IV, 16)
	UserKey, err := userlib.HashKDF(MasterKey, []byte("UserKey"))
	if err != nil {
		return nil, err
	}

	userdata.MasterKey = MasterKey

	pub, priv, err := userlib.PKEKeyGen()
	if err != nil {
		return nil, err
	}
	userdata.RSAPrivateKey = priv
	userlib.KeystoreSet(username + "RSA", pub)

	sign, verify, err := userlib.DSKeyGen()
	if err != nil {
		return nil, err
	}
	userdata.DigitalSigningKey = sign
	userlib.KeystoreSet(username + "DigSig", verify)

	userBytes, err := json.Marshal(userdata)
	if err != nil {
		return nil, err
	}
	
	ciphertext := userlib.SymEnc(UserKey[:16], IV, userBytes)
	mac, err := userlib.HMACEval(UserKey[16:32], ciphertext)
	if err != nil {
		return nil, err
	}

	partial := append(ciphertext, mac...)
	final := append(partial, IV...)

	userlib.DatastoreSet(LoginUUID, final)

	return &userdata, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {


	LoginUUID, err := uuid.FromBytes(userlib.Hash([]byte(username))[:16])
	if err != nil {
		return nil, err
	}

	data, ok := userlib.DatastoreGet(LoginUUID)
	if !ok {
		return nil, errors.New("GetUser: error getting Data")
	}

	if len(data) < 80 {
		return nil, errors.New("GetUser: data not long enough")
	}

	IV := data[len(data)-16:]
	ciphertext := data[:len(data)-80]

	MasterKey := userlib.Argon2Key([]byte(password), IV, 16)
	UserKey, err := userlib.HashKDF(MasterKey, []byte("UserKey"))

	mac, err := userlib.HMACEval(UserKey[16:32], ciphertext)
	if err != nil {
		return nil, err
	}

	if !userlib.HMACEqual(mac, data[len(data)-80:len(data)-16]) {
		return nil, errors.New("GetUser: HMAC not matching")
	}

	marshaleduser := userlib.SymDec(UserKey[:16], ciphertext)
	var userdata User
	err = json.Unmarshal(marshaleduser, &userdata)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return &userdata, nil
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {

	
	FileHeaderUUID, err := uuid.FromBytes(userlib.RandomBytes(16))
	if err != nil {
		return err
	}

	FileNameKey, _ := userlib.HashKDF(userdata.MasterKey, []byte(filename))

	UserKey, err := userlib.HashKDF(userdata.MasterKey, []byte("UserKey"))

	FileInfoUUID, err := uuid.FromBytes(FileNameKey[:16])
	if err != nil {
		return err
	}

	_, ok := userlib.DatastoreGet(FileInfoUUID)

	if ok {
		var fileinfo FileInfo
		err = getCiphAndCheckMAC(FileInfoUUID, UserKey[:16], UserKey[16:32], &fileinfo)
		if err != nil {
			return err
		}


		var fileheaderuuid uuid.UUID
		var fileenckey []byte
		var filemackey []byte
		
		if fileinfo.Owner {
			fileheaderuuid = fileinfo.FileHeaderUUID
			fileenckey = fileinfo.EncKey
			filemackey = fileinfo.MACKey
		} else {
			var ap AccessPoint
			err = getCiphAndCheckMAC(fileinfo.AccessPointUUID, fileinfo.EncKey, fileinfo.MACKey, &ap)
			if err != nil {
				return err
			}
			if ap.Revoked {
				return errors.New("access revoked")
			}
			fileheaderuuid = ap.FileHeaderUUID
			fileenckey = ap.EncKey
			filemackey = ap.MACKey
		}

		var fileheader FileHeader
		err = getCiphAndCheckMAC(fileheaderuuid, fileenckey, filemackey, &fileheader)
		if err != nil {
			return err
		}

		FirstUUID, err := uuid.FromBytes(userlib.RandomBytes(16))
		if err != nil {
			return err
		}
		NewUUID, err := uuid.FromBytes(userlib.RandomBytes(16))
		if err != nil {
			return err
		}

		fileheader.FirstUUID = FirstUUID
		fileheader.NewUUID = NewUUID

		var file File
		file.Content = content
		file.NextUUID = NewUUID


		err = setFinalCiph(FirstUUID, fileenckey, filemackey, file)
		if err != nil {
			return err
		}


		err = setFinalCiph(fileheaderuuid, fileenckey, filemackey, fileheader)
		if err != nil {
			return err
		}

		return nil
	}

	var fileinfo FileInfo
	fileinfo.Owner = true
	fileinfo.FileHeaderUUID = FileHeaderUUID
	fileinfo.EncKey = userlib.RandomBytes(16)
	fileinfo.MACKey = userlib.RandomBytes(16)
	fileinfo.OwnerAccessInfoMapUUID, _ = uuid.FromBytes(userlib.RandomBytes(16))

	
	FirstUUID, err := uuid.FromBytes(userlib.RandomBytes(16))
	if err != nil {
		return err
	}
	NewUUID, err := uuid.FromBytes(userlib.RandomBytes(16))
	if err != nil {
		return err
	}

	var fileheader FileHeader
	fileheader.FirstUUID = FirstUUID
	fileheader.NewUUID = NewUUID

	var file File
	file.Content = content
	file.NextUUID = NewUUID


	err = setFinalCiph(FirstUUID, fileinfo.EncKey, fileinfo.MACKey, file)
	if err != nil {
		return err
	}


	err = setFinalCiph(FileHeaderUUID, fileinfo.EncKey, fileinfo.MACKey, fileheader)
	if err != nil {
		return err
	}

	err = setFinalCiph(FileInfoUUID, UserKey[:16], UserKey[16:32], fileinfo)
	if err != nil {
		return err
	}

	var accmap OwnerAccessInfoMap
	accmap.InfoMap = make(map[string]OwnerAccessInfo)
	err = setFinalCiph(fileinfo.OwnerAccessInfoMapUUID, UserKey[:16], UserKey[16:32], accmap)
	if err != nil {
		return err
	}

	return
}

func (userdata *User) AppendToFile(filename string, content []byte) error {

	FileNameKey, _ := userlib.HashKDF(userdata.MasterKey, []byte(filename))
	UserKey, err := userlib.HashKDF(userdata.MasterKey, []byte("UserKey"))

	FileInfoUUID, err := uuid.FromBytes(FileNameKey[:16])
	if err != nil {
		return err
	}

	var fileinfo FileInfo
	err = getCiphAndCheckMAC(FileInfoUUID, UserKey[:16], UserKey[16:32], &fileinfo)
	if err != nil {
		return err
	}


	var fileheaderuuid uuid.UUID
	var fileenckey []byte
	var filemackey []byte


	if fileinfo.Owner {
		fileheaderuuid = fileinfo.FileHeaderUUID
		fileenckey = fileinfo.EncKey
		filemackey = fileinfo.MACKey
	} else {
		var ap AccessPoint
		err = getCiphAndCheckMAC(fileinfo.AccessPointUUID, fileinfo.EncKey, fileinfo.MACKey, &ap)
		if err != nil {
			return err
		}
		if ap.Revoked {
			return errors.New("access revoked")
		}
		fileheaderuuid = ap.FileHeaderUUID
		fileenckey = ap.EncKey
		filemackey = ap.MACKey
	}

	var fileheader FileHeader
	err = getCiphAndCheckMAC(fileheaderuuid, fileenckey, filemackey, &fileheader)
	if err != nil {
		return err
	}
	
	nextuuid, err := uuid.FromBytes(userlib.RandomBytes(16))
	if err != nil {
		return err
	}

	var file File
	file.Content = content
	file.NextUUID = nextuuid

	err = setFinalCiph(fileheader.NewUUID, fileenckey, filemackey, file)
	if err != nil {
		return err
	}

	fileheader.NewUUID = nextuuid
	err = setFinalCiph(fileheaderuuid, fileenckey, filemackey, fileheader)
	if err != nil {
		return err
	}

	return nil
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {

	FileNameKey, _ := userlib.HashKDF(userdata.MasterKey, []byte(filename))
	UserKey, err := userlib.HashKDF(userdata.MasterKey, []byte("UserKey"))

	FileInfoUUID, err := uuid.FromBytes(FileNameKey[:16])
	if err != nil {
		return nil, err
	}

	var fileinfo FileInfo
	err = getCiphAndCheckMAC(FileInfoUUID, UserKey[:16], UserKey[16:32], &fileinfo)
	if err != nil {
		return nil, err
	}

	var fileheaderuuid uuid.UUID
	var fileenckey []byte
	var filemackey []byte


	if fileinfo.Owner {
		fileheaderuuid = fileinfo.FileHeaderUUID
		fileenckey = fileinfo.EncKey
		filemackey = fileinfo.MACKey
		
	} else {
		var ap AccessPoint
		err = getCiphAndCheckMAC(fileinfo.AccessPointUUID, fileinfo.EncKey, fileinfo.MACKey, &ap)
		if err != nil {
			return nil, err
		}
		if ap.Revoked {
			return nil, errors.New("access revoked")
		}
		fileheaderuuid = ap.FileHeaderUUID
		fileenckey = ap.EncKey
		filemackey = ap.MACKey
	}

	var fileheader FileHeader
	err = getCiphAndCheckMAC(fileheaderuuid, fileenckey, filemackey, &fileheader)
	if err != nil {
		return nil, err
	}


	nextuuid := fileheader.FirstUUID


	for nextuuid != fileheader.NewUUID {
		var file File
		err = getCiphAndCheckMAC(nextuuid, fileenckey, filemackey, &file)
		if err != nil {
			return nil, err
		}
		content = append(content, file.Content...)
		nextuuid = file.NextUUID
	}


	return content, nil
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (
	invitationPtr uuid.UUID, err error) {

	FileNameKey, _ := userlib.HashKDF(userdata.MasterKey, []byte(filename))
	UserKey, err := userlib.HashKDF(userdata.MasterKey, []byte("UserKey"))

	FileInfoUUID, err := uuid.FromBytes(FileNameKey[:16])
	if err != nil {
		return uuid.Nil, err
	}

	var fileinfo FileInfo
	err = getCiphAndCheckMAC(FileInfoUUID, UserKey[:16], UserKey[16:32], &fileinfo)
	if err != nil {
		return uuid.Nil, err
	}

	if fileinfo.Owner {

		invitationuuid, err := uuid.FromBytes(userlib.RandomBytes(16))
		if err != nil {
			return uuid.Nil, err
		}
		accesspointuuid, err := uuid.FromBytes(userlib.RandomBytes(16))
		if err != nil {
			return uuid.Nil, err
		}
		apenckey := userlib.RandomBytes(16)
		apmackey := userlib.RandomBytes(16)

		var inv Invitation
		inv.AccessPointUUID = accesspointuuid
		inv.APEncKey = apenckey
		inv.APMACKey = apmackey

		err = setFinalCiphRSA(invitationuuid, recipientUsername, userdata.DigitalSigningKey, inv)
		if err != nil {
			return uuid.Nil, err
		}

		var ap AccessPoint
		ap.FileHeaderUUID = fileinfo.FileHeaderUUID
		ap.EncKey = fileinfo.EncKey
		ap.MACKey = fileinfo.MACKey

		err = setFinalCiph(accesspointuuid, apenckey, apmackey, ap)
		if err != nil {
			return uuid.Nil, err
		}

		var accmap OwnerAccessInfoMap
		err = getCiphAndCheckMAC(fileinfo.OwnerAccessInfoMapUUID, UserKey[:16], UserKey[16:32], &accmap)
		if err != nil {
			return uuid.Nil, err
		}

		var info OwnerAccessInfo
		info.AccessPointUUID = accesspointuuid
		info.APEncKey = apenckey
		info.APMACKey = apmackey
		accmap.InfoMap[recipientUsername] = info
		err = setFinalCiph(fileinfo.OwnerAccessInfoMapUUID, UserKey[:16], UserKey[16:32], accmap)
		if err != nil {
			return uuid.Nil, err
		}

		return invitationuuid, nil


	} else {

		invitationuuid, err := uuid.FromBytes(userlib.RandomBytes(16))
		if err != nil {
			return uuid.Nil, err
		}

		var inv Invitation
		inv.AccessPointUUID = fileinfo.AccessPointUUID
		inv.APEncKey = fileinfo.EncKey
		inv.APMACKey = fileinfo.MACKey

		err = setFinalCiphRSA(invitationuuid, recipientUsername, userdata.DigitalSigningKey, inv)
		if err != nil {
			return uuid.Nil, err
		}

		return invitationuuid, nil
	}

	 
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {

	var inv Invitation
	err := getCiphAndCheckMACRSA(invitationPtr, senderUsername, userdata.RSAPrivateKey, &inv)
	if err != nil {
		return err
	}

	var ap AccessPoint
	err = getCiphAndCheckMAC(inv.AccessPointUUID, inv.APEncKey, inv.APMACKey, &ap)
	if err != nil {
		return err
	}
	if ap.Revoked {
		return errors.New("access revoked")
	}

	FileNameKey, _ := userlib.HashKDF(userdata.MasterKey, []byte(filename))
	UserKey, err := userlib.HashKDF(userdata.MasterKey, []byte("UserKey"))

	FileInfoUUID, err := uuid.FromBytes(FileNameKey[:16])
	if err != nil {
		return err
	}

	_, ok := userlib.DatastoreGet(FileInfoUUID)
	if ok {
		return errors.New("cannot accept invite with a name that already exists")
	}

	var fileinfo FileInfo
	fileinfo.Owner = false
	fileinfo.AccessPointUUID = inv.AccessPointUUID
	fileinfo.EncKey = inv.APEncKey
	fileinfo.MACKey = inv.APMACKey
	fileinfo.OwnerAccessInfoMapUUID = uuid.Nil

	err = setFinalCiph(FileInfoUUID, UserKey[:16], UserKey[16:32], fileinfo)
	if err != nil {
		return err
	}

	return nil
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {

	FileNameKey, _ := userlib.HashKDF(userdata.MasterKey, []byte(filename))
	UserKey, err := userlib.HashKDF(userdata.MasterKey, []byte("UserKey"))

	FileInfoUUID, err := uuid.FromBytes(FileNameKey[:16])
	if err != nil {
		return err
	}

	var fileinfo FileInfo
	err = getCiphAndCheckMAC(FileInfoUUID, UserKey[:16], UserKey[16:32], &fileinfo)
	if err != nil {
		return err
	}

	var accmap OwnerAccessInfoMap
	err = getCiphAndCheckMAC(fileinfo.OwnerAccessInfoMapUUID, UserKey[:16], UserKey[16:32], &accmap)
	if err != nil {
		return err
	}

	info := accmap.InfoMap[recipientUsername]
	
	var ap AccessPoint
	err = getCiphAndCheckMAC(info.AccessPointUUID, info.APEncKey, info.APMACKey, &ap)
	if err != nil {
		return err
	}

	ap.Revoked = true
	err = setFinalCiph(info.AccessPointUUID, info.APEncKey, info.APMACKey, ap)
	if err != nil {
		return err
	}

	delete(accmap.InfoMap, recipientUsername)

	err = setFinalCiph(fileinfo.OwnerAccessInfoMapUUID, UserKey[:16], UserKey[16:32], accmap)
	if err != nil {
		return err
	}

	newFileHeaderUUID, err := uuid.FromBytes(userlib.RandomBytes(16))
	if err != nil {
		return err
	}
	newFirstUUID, err := uuid.FromBytes(userlib.RandomBytes(16))
	if err != nil {
		return err
	}
	newNewUUID, err := uuid.FromBytes(userlib.RandomBytes(16))
	if err != nil {
		return err
	}
	newEncKey := userlib.RandomBytes(16)
	newMACKey := userlib.RandomBytes(16)


	data, err := userdata.LoadFile(filename)

	var fileheader FileHeader
	err = getCiphAndCheckMAC(fileinfo.FileHeaderUUID, fileinfo.EncKey, fileinfo.MACKey, &fileheader)
	if err != nil {
		return err
	}

	fileheader.FirstUUID = newFirstUUID
	fileheader.NewUUID = newNewUUID

	err = setFinalCiph(newFileHeaderUUID, newEncKey, newMACKey, fileheader)
	if err != nil {
		return err
	}

	var file File
	file.Content = data
	file.NextUUID = newNewUUID

	err = setFinalCiph(newFirstUUID, newEncKey, newMACKey, file)
	if err != nil {
		return err
	}

	fileinfo.FileHeaderUUID = newFileHeaderUUID
	fileinfo.EncKey = newEncKey
	fileinfo.MACKey = newMACKey
	err = setFinalCiph(FileInfoUUID, UserKey[:16], UserKey[16:32], &fileinfo)
	if err != nil {
		return err
	}

	for _, v := range accmap.InfoMap {

		var ap AccessPoint
		err = getCiphAndCheckMAC(v.AccessPointUUID, v.APEncKey, v.APMACKey, &ap)
		if err != nil {
			return err
		}
		ap.FileHeaderUUID = newFileHeaderUUID
		ap.EncKey = newEncKey
		ap.MACKey = newMACKey

		err = setFinalCiph(v.AccessPointUUID, v.APEncKey, v.APMACKey, ap)
		if err != nil {
			return err
		}
	}

	return nil
}


func getCiphAndCheckMAC (uuidLocation uuid.UUID, EncKey []byte, MACKey []byte, object interface{}) (error) {
	
	data, ok := userlib.DatastoreGet(uuidLocation)
	if !ok {
		return errors.New("getCiphAndCheckMAC: error getting Data")
	}
	
	if len(data) < 64 {
		return errors.New("getCiphAndCheckMAC: data not long enough")
	}

	ciphertext := data[:len(data)-64]

	mac, err := userlib.HMACEval(MACKey, ciphertext)
	if err != nil {
		return err
	}

	if !userlib.HMACEqual(mac, data[len(data)-64:]) {
		return errors.New("getCiphAndCheckMAC: HMAC not matching")
	}

	decripted := userlib.SymDec(EncKey, ciphertext)

	err = json.Unmarshal(decripted, object)
	if err != nil {
		return err
	}

	return nil
}


func setFinalCiph (uuidLocation uuid.UUID, EncKey []byte, MACKey []byte, object interface{}) (error) {
	
	bytes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	
	ciphertext := userlib.SymEnc(EncKey, userlib.RandomBytes(16), bytes)
	mac, err := userlib.HMACEval(MACKey, ciphertext)
	if err != nil {
		return err
	}

	userlib.DatastoreSet(uuidLocation, append(ciphertext, mac...))

	return nil
}
func setFinalCiphRSA (uuidLocation uuid.UUID, recieverUsername string, senderSigKey userlib.PrivateKeyType, object interface{}) (error) {
	
	bytes, err := json.Marshal(object)
	if err != nil {
		return err
	}

	enckey, ok := userlib.KeystoreGet(recieverUsername + "RSA")
	if !ok {
		return errors.New("key not there")
	}

	symKey := userlib.RandomBytes(16)
	
	encriptedKey, err := userlib.PKEEnc(enckey, symKey)
	if err != nil {
		return err
	}

	ciphertext := userlib.SymEnc(symKey, userlib.RandomBytes(16), bytes)
	if err != nil {
		return err
	}

	ciphertext = append(ciphertext, encriptedKey...)

	digsig, err := userlib.DSSign(senderSigKey, ciphertext)
	
	if err != nil {
		return err
	}

	userlib.DatastoreSet(uuidLocation, append(ciphertext, digsig...))

	return nil
}
func getCiphAndCheckMACRSA (uuidLocation uuid.UUID, senderUsername string, recieverPrivateKey userlib.PrivateKeyType, object interface{}) (error) {
	
	data, ok := userlib.DatastoreGet(uuidLocation)
	if !ok {
		return errors.New("getCiphAndCheckMACRSA: error getting Data")
	}
	
	if len(data) < 512 {
		return errors.New("getCiphAndCheckMACRSA: data not long enough")
	}

	ciphertext := data[:len(data)-512]

	verificationKey, ok := userlib.KeystoreGet(senderUsername + "DigSig")
	if !ok {
		return errors.New("key not there")
	}
	
	err := userlib.DSVerify(verificationKey, data[:len(data)-256], data[len(data)-256:])
	if err != nil {
		return err
	}

	SymKey, err := userlib.PKEDec(recieverPrivateKey, data[len(data) - 512 :len(data) - 256])
	if err != nil {
		return err
	}

	decripted := userlib.SymDec(SymKey, ciphertext)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decripted, object)
	if err != nil {
		return err
	}

	return nil
}
