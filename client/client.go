package client

// CS 161 Project 2

// may break the autograder!

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string


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

	// Declares a Course struct type, creates an instance of it, and prints it.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Course: %v", course)
	userlib.DebugMsg("Course Bytes: %v", courseBytes)

	publicKey, privateKey, err := userlib.PKEKeyGen()
	if err != nil {
		panic(err)
	}

	encryptedCourse, err := userlib.PKEEnc(publicKey, courseBytes)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Encrypted Course: %v", encryptedCourse)

	decryptedCourse, err := userlib.PKEDec(privateKey, encryptedCourse)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Decrypted Course: %v", decryptedCourse)

	var courseFromJSON Course
	err = json.Unmarshal(decryptedCourse, &courseFromJSON)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Course From JSON: %v", courseFromJSON)
}

type UserData struct {
	Username    string                 `json:"username"`
	RootKey     []byte                 `json:"root_key"`
	FileMap     map[string]uuid.UUID   `json:"file_map"`
	PublicKey   userlib.PKEEncKey      `json:"public_key"`
	PrivateKey  userlib.PKEDecKey      `json:"private_key"`
	SignKey     userlib.DSSignKey      `json:"sign_key"`
	VerifyKey   userlib.DSVerifyKey    `json:"verify_key"`
}

type FileMetadata struct {
	FileUUID    uuid.UUID `json:"file_uuid"`
	EncKey      []byte    `json:"enc_key"`
	MacKey      []byte    `json:"mac_key"`
	Owner       string    `json:"owner"`
	AppendHead  uuid.UUID `json:"append_head"`
}

type FileData struct {
	Content []byte `json:"content"`
	HMAC    []byte `json:"hmac"`
}

type AppendNode struct {
	Content []byte    `json:"content"`
	Next    uuid.UUID `json:"next"`
	HMAC    []byte    `json:"hmac"`
}

type Invitation struct {
	FileMetadataUUID uuid.UUID `json:"file_metadata_uuid"`
	EncKey           []byte    `json:"enc_key"`
	MacKey           []byte    `json:"mac_key"`
	Signature        []byte    `json:"signature"`
}

type User struct {
	Username string
	RootKey  []byte
}

func InitUser(username string, password string) (userdataptr *User, err error) {
	if len(username) == 0 {
		return nil, errors.New("username cannot be empty")
	}

	userUUID, err := uuid.FromBytes(userlib.Hash([]byte(username))[:16])
	if err != nil {
		return nil, err
	}

	_, exists := userlib.DatastoreGet(userUUID)
	if exists {
		return nil, errors.New("user already exists")
	}

	salt := userlib.RandomBytes(16)
	rootKey := userlib.Argon2Key([]byte(password), salt, 16)

	publicKey, privateKey, err := userlib.PKEKeyGen()
	if err != nil {
		return nil, err
	}

	signKey, verifyKey, err := userlib.DSKeyGen()
	if err != nil {
		return nil, err
	}

	err = userlib.KeystoreSet(username+"_public", publicKey)
	if err != nil {
		return nil, err
	}

	err = userlib.KeystoreSet(username+"_verify", verifyKey)
	if err != nil {
		return nil, err
	}

	userData := UserData{
		Username:    username,
		RootKey:     rootKey,
		FileMap:     make(map[string]uuid.UUID),
		PublicKey:   publicKey,
		PrivateKey:  privateKey,
		SignKey:     signKey,
		VerifyKey:   verifyKey,
	}

	userDataBytes, err := json.Marshal(userData)
	if err != nil {
		return nil, err
	}

	encKey, err := userlib.HashKDF(rootKey, []byte("user_enc"))
	if err != nil {
		return nil, err
	}
	encKey = encKey[:16]

	macKey, err := userlib.HashKDF(rootKey, []byte("user_mac"))
	if err != nil {
		return nil, err
	}
	macKey = macKey[:16]

	iv := userlib.RandomBytes(16)
	encryptedData := userlib.SymEnc(encKey, iv, userDataBytes)

	hmac, err := userlib.HMACEval(macKey, encryptedData)
	if err != nil {
		return nil, err
	}

	storedData := append(encryptedData, hmac...)
	storedData = append(salt, storedData...)

	userlib.DatastoreSet(userUUID, storedData)

	user := &User{
		Username: username,
		RootKey:  rootKey,
	}

	return user, nil
}

func GetUser(username string, password string) (userdataptr *User, err error) {
	if len(username) == 0 {
		return nil, errors.New("username cannot be empty")
	}

	userUUID, err := uuid.FromBytes(userlib.Hash([]byte(username))[:16])
	if err != nil {
		return nil, err
	}

	storedData, exists := userlib.DatastoreGet(userUUID)
	if !exists {
		return nil, errors.New("user does not exist")
	}

	if len(storedData) < 16 {
		return nil, errors.New("invalid user data")
	}

	salt := storedData[:16]
	encryptedData := storedData[16 : len(storedData)-64]
	storedHMAC := storedData[len(storedData)-64:]

	rootKey := userlib.Argon2Key([]byte(password), salt, 16)

	macKey, err := userlib.HashKDF(rootKey, []byte("user_mac"))
	if err != nil {
		return nil, err
	}
	macKey = macKey[:16]

	computedHMAC, err := userlib.HMACEval(macKey, encryptedData)
	if err != nil {
		return nil, err
	}

	if !userlib.HMACEqual(storedHMAC, computedHMAC) {
		return nil, errors.New("invalid password or data tampering detected")
	}

	encKey, err := userlib.HashKDF(rootKey, []byte("user_enc"))
	if err != nil {
		return nil, err
	}
	encKey = encKey[:16]

	userDataBytes := userlib.SymDec(encKey, encryptedData)

	var userData UserData
	err = json.Unmarshal(userDataBytes, &userData)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username: username,
		RootKey:  rootKey,
	}

	return user, nil
}

func (userdata *User) getUserData() (UserData, error) {
	userUUID, err := uuid.FromBytes(userlib.Hash([]byte(userdata.Username))[:16])
	if err != nil {
		return UserData{}, err
	}

	storedData, exists := userlib.DatastoreGet(userUUID)
	if !exists {
		return UserData{}, errors.New("user data not found")
	}

	if len(storedData) < 16 {
		return UserData{}, errors.New("invalid user data")
	}

	_ = storedData[:16]
	encryptedData := storedData[16 : len(storedData)-64]
	storedHMAC := storedData[len(storedData)-64:]

	macKey, err := userlib.HashKDF(userdata.RootKey, []byte("user_mac"))
	if err != nil {
		return UserData{}, err
	}
	macKey = macKey[:16]

	computedHMAC, err := userlib.HMACEval(macKey, encryptedData)
	if err != nil {
		return UserData{}, err
	}

	if !userlib.HMACEqual(storedHMAC, computedHMAC) {
		return UserData{}, errors.New("user data integrity check failed")
	}

	encKey, err := userlib.HashKDF(userdata.RootKey, []byte("user_enc"))
	if err != nil {
		return UserData{}, err
	}
	encKey = encKey[:16]

	userDataBytes := userlib.SymDec(encKey, encryptedData)

	var userData UserData
	err = json.Unmarshal(userDataBytes, &userData)
	if err != nil {
		return UserData{}, err
	}

	return userData, nil
}

func (userdata *User) saveUserData(userData UserData) error {
	userUUID, err := uuid.FromBytes(userlib.Hash([]byte(userdata.Username))[:16])
	if err != nil {
		return err
	}

	userDataBytes, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	encKey, err := userlib.HashKDF(userdata.RootKey, []byte("user_enc"))
	if err != nil {
		return err
	}
	encKey = encKey[:16]

	macKey, err := userlib.HashKDF(userdata.RootKey, []byte("user_mac"))
	if err != nil {
		return err
	}
	macKey = macKey[:16]

	salt := userlib.RandomBytes(16)
	iv := userlib.RandomBytes(16)
	encryptedData := userlib.SymEnc(encKey, iv, userDataBytes)

	hmac, err := userlib.HMACEval(macKey, encryptedData)
	if err != nil {
		return err
	}

	storedData := append(encryptedData, hmac...)
	storedData = append(salt, storedData...)

	userlib.DatastoreSet(userUUID, storedData)
	return nil
}

func (userdata *User) getFileMetadata(filename string) (FileMetadata, error) {
	userData, err := userdata.getUserData()
	if err != nil {
		return FileMetadata{}, err
	}

	metadataUUID, exists := userData.FileMap[filename]
	if !exists {
		return FileMetadata{}, errors.New("file not found")
	}

	metadataBytes, exists := userlib.DatastoreGet(metadataUUID)
	if !exists {
		return FileMetadata{}, errors.New("file metadata not found")
	}

	var metadata FileMetadata
	err = json.Unmarshal(metadataBytes, &metadata)
	if err != nil {
		return FileMetadata{}, err
	}

	return metadata, nil
}

func (userdata *User) saveFileMetadata(metadataUUID uuid.UUID, metadata FileMetadata) error {
	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	userlib.DatastoreSet(metadataUUID, metadataBytes)
	return nil
}

func (userdata *User) loadAppendData(appendUUID uuid.UUID, encKey []byte, macKey []byte) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (userdata *User) loadAppendDataWithKeys(appendUUID uuid.UUID, encKey []byte, macKey []byte) ([]byte, error) {
	var allContent []byte
	currentUUID := appendUUID

	for currentUUID != uuid.Nil {
		nodeBytes, exists := userlib.DatastoreGet(currentUUID)
		if !exists {
			break
		}

		if len(nodeBytes) < 64 {
			return nil, errors.New("invalid append node data")
		}

		encryptedData := nodeBytes[:len(nodeBytes)-64]
		storedHMAC := nodeBytes[len(nodeBytes)-64:]

		computedHMAC, err := userlib.HMACEval(macKey, encryptedData)
		if err != nil {
			return nil, err
		}

		if !userlib.HMACEqual(storedHMAC, computedHMAC) {
			return nil, errors.New("append node integrity check failed")
		}

		nodeDataBytes := userlib.SymDec(encKey, encryptedData)

		var node AppendNode
		err = json.Unmarshal(nodeDataBytes, &node)
		if err != nil {
			return nil, err
		}

		allContent = append(allContent, node.Content...)
		currentUUID = node.Next
	}

	return allContent, nil
}

func (userdata *User) deleteAppendChain(appendUUID uuid.UUID) error {
	currentUUID := appendUUID

	for currentUUID != uuid.Nil {
		nodeBytes, exists := userlib.DatastoreGet(currentUUID)
		if !exists {
			break
		}

		userlib.DatastoreDelete(currentUUID)

		if len(nodeBytes) >= 64 {
			encryptedData := nodeBytes[:len(nodeBytes)-64]
			var node AppendNode
			json.Unmarshal(encryptedData, &node)
			currentUUID = node.Next
		} else {
			break
		}
	}

	return nil
}

func (userdata *User) updateInvitationForUser(recipientUsername string, metadataUUID uuid.UUID, encKey []byte, macKey []byte) (uuid.UUID, error) {
	recipientPublicKey, ok := userlib.KeystoreGet(recipientUsername + "_public")
	if !ok {
		return uuid.Nil, errors.New("recipient public key not found")
	}

	invitation := Invitation{
		FileMetadataUUID: metadataUUID,
		EncKey:           encKey,
		MacKey:           macKey,
	}

	invitationBytes, err := json.Marshal(invitation)
	if err != nil {
		return uuid.Nil, err
	}

	userData, err := userdata.getUserData()
	if err != nil {
		return uuid.Nil, err
	}

	signature, err := userlib.DSSign(userData.SignKey, invitationBytes)
	if err != nil {
		return uuid.Nil, err
	}

	invitation.Signature = signature

	finalInvitationBytes, err := json.Marshal(invitation)
	if err != nil {
		return uuid.Nil, err
	}

	encryptedInvitation, err := userlib.PKEEnc(recipientPublicKey, finalInvitationBytes)
	if err != nil {
		return uuid.Nil, err
	}

	invitationUUID := uuid.New()
	userlib.DatastoreSet(invitationUUID, encryptedInvitation)

	return invitationUUID, nil
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {
	userData, err := userdata.getUserData()
	if err != nil {
		return err
	}

	fileEncKey, err := userlib.HashKDF(userdata.RootKey, []byte("file_enc_"+filename))
	if err != nil {
		return err
	}
	fileEncKey = fileEncKey[:16]

	fileMacKey, err := userlib.HashKDF(userdata.RootKey, []byte("file_mac_"+filename))
	if err != nil {
		return err
	}
	fileMacKey = fileMacKey[:16]

	fileUUID := uuid.New()

	fileData := FileData{
		Content: content,
	}

	fileDataBytes, err := json.Marshal(fileData)
	if err != nil {
		return err
	}

	iv := userlib.RandomBytes(16)
	encryptedFileData := userlib.SymEnc(fileEncKey, iv, fileDataBytes)

	hmac, err := userlib.HMACEval(fileMacKey, encryptedFileData)
	if err != nil {
		return err
	}

	finalFileData := append(encryptedFileData, hmac...)
	userlib.DatastoreSet(fileUUID, finalFileData)

	metadata := FileMetadata{
		FileUUID:   fileUUID,
		EncKey:     fileEncKey,
		MacKey:     fileMacKey,
		Owner:      userdata.Username,
		AppendHead: uuid.Nil,
	}

	metadataUUID := uuid.New()
	err = userdata.saveFileMetadata(metadataUUID, metadata)
	if err != nil {
		return err
	}

	userData.FileMap[filename] = metadataUUID
	return userdata.saveUserData(userData)
}

func (userdata *User) AppendToFile(filename string, content []byte) error {
	metadata, err := userdata.getFileMetadata(filename)
	if err != nil {
		return err
	}

	appendUUID := uuid.New()

	appendNode := AppendNode{
		Content: content,
		Next:    metadata.AppendHead,
	}

	appendNodeBytes, err := json.Marshal(appendNode)
	if err != nil {
		return err
	}

	iv := userlib.RandomBytes(16)
	encryptedAppendData := userlib.SymEnc(metadata.EncKey, iv, appendNodeBytes)

	hmac, err := userlib.HMACEval(metadata.MacKey, encryptedAppendData)
	if err != nil {
		return err
	}

	finalAppendData := append(encryptedAppendData, hmac...)
	userlib.DatastoreSet(appendUUID, finalAppendData)

	metadata.AppendHead = appendUUID

	userData, err := userdata.getUserData()
	if err != nil {
		return err
	}

	metadataUUID := userData.FileMap[filename]
	return userdata.saveFileMetadata(metadataUUID, metadata)
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	metadata, err := userdata.getFileMetadata(filename)
	if err != nil {
		return nil, err
	}

	fileBytes, exists := userlib.DatastoreGet(metadata.FileUUID)
	if !exists {
		return nil, errors.New("file data not found")
	}

	if len(fileBytes) < 64 {
		return nil, errors.New("invalid file data")
	}

	encryptedData := fileBytes[:len(fileBytes)-64]
	storedHMAC := fileBytes[len(fileBytes)-64:]

	computedHMAC, err := userlib.HMACEval(metadata.MacKey, encryptedData)
	if err != nil {
		return nil, err
	}

	if !userlib.HMACEqual(storedHMAC, computedHMAC) {
		return nil, errors.New("file integrity check failed")
	}

	fileDataBytes := userlib.SymDec(metadata.EncKey, encryptedData)

	var fileData FileData
	err = json.Unmarshal(fileDataBytes, &fileData)
	if err != nil {
		return nil, err
	}

	content = fileData.Content

	// if metadata.AppendHead != uuid.Nil {
	//	appendContent, err := userdata.loadAppendDataWithKeys(metadata.AppendHead, metadata.EncKey, metadata.MacKey)
	//	if err != nil {
	//		return nil, err
	//	}
	//	content = append(content, appendContent...)
	// }

	return content, nil
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (invitationPtr uuid.UUID, err error) {
	if recipientUsername == userdata.Username {
		return uuid.Nil, errors.New("cannot invite yourself")
	}

	_, ok := userlib.KeystoreGet(recipientUsername + "_public")
	if !ok {
		return uuid.Nil, errors.New("recipient user does not exist")
	}

	userData, err := userdata.getUserData()
	if err != nil {
		return uuid.Nil, err
	}

	metadataUUID, exists := userData.FileMap[filename]
	if !exists {
		return uuid.Nil, errors.New("file not found")
	}

	metadata, err := userdata.getFileMetadata(filename)
	if err != nil {
		return uuid.Nil, err
	}

	return userdata.updateInvitationForUser(recipientUsername, metadataUUID, metadata.EncKey, metadata.MacKey)
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	userData, err := userdata.getUserData()
	if err != nil {
		return err
	}

	_, exists := userData.FileMap[filename]
	if exists {
		return errors.New("filename already exists")
	}

	encryptedInvitation, exists := userlib.DatastoreGet(invitationPtr)
	if !exists {
		return errors.New("invitation not found")
	}

	invitationBytes, err := userlib.PKEDec(userData.PrivateKey, encryptedInvitation)
	if err != nil {
		return errors.New("failed to decrypt invitation")
	}

	var invitation Invitation
	err = json.Unmarshal(invitationBytes, &invitation)
	if err != nil {
		return err
	}

	senderVerifyKey, ok := userlib.KeystoreGet(senderUsername + "_verify")
	if !ok {
		return errors.New("sender verify key not found")
	}

	invitationForVerification := Invitation{
		FileMetadataUUID: invitation.FileMetadataUUID,
		EncKey:           invitation.EncKey,
		MacKey:           invitation.MacKey,
	}

	invitationBytesForVerification, err := json.Marshal(invitationForVerification)
	if err != nil {
		return err
	}

	err = userlib.DSVerify(senderVerifyKey, invitationBytesForVerification, invitation.Signature)
	if err != nil {
		return errors.New("invitation signature verification failed")
	}

	userData.FileMap[filename] = invitation.FileMetadataUUID
	return userdata.saveUserData(userData)
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {
	userData, err := userdata.getUserData()
	if err != nil {
		return err
	}

	metadataUUID, exists := userData.FileMap[filename]
	if !exists {
		return errors.New("file not found")
	}

	metadata, err := userdata.getFileMetadata(filename)
	if err != nil {
		return err
	}

	if metadata.Owner != userdata.Username {
		return errors.New("only file owner can revoke access")
	}

	content, err := userdata.LoadFile(filename)
	if err != nil {
		return err
	}

	if metadata.AppendHead != uuid.Nil {
		err = userdata.deleteAppendChain(metadata.AppendHead)
		if err != nil {
			return err
		}
	}

	userlib.DatastoreDelete(metadata.FileUUID)

	newFileEncKey := userlib.RandomBytes(16)
	newFileMacKey := userlib.RandomBytes(64)

	newFileUUID := uuid.New()

	fileData := FileData{
		Content: content,
	}

	fileDataBytes, err := json.Marshal(fileData)
	if err != nil {
		return err
	}

	iv := userlib.RandomBytes(16)
	encryptedFileData := userlib.SymEnc(newFileEncKey, iv, fileDataBytes)

	hmac, err := userlib.HMACEval(newFileMacKey, encryptedFileData)
	if err != nil {
		return err
	}

	finalFileData := append(encryptedFileData, hmac...)
	userlib.DatastoreSet(newFileUUID, finalFileData)

	newMetadata := FileMetadata{
		FileUUID:   newFileUUID,
		EncKey:     newFileEncKey,
		MacKey:     newFileMacKey,
		Owner:      userdata.Username,
		AppendHead: uuid.Nil,
	}

	err = userdata.saveFileMetadata(metadataUUID, newMetadata)
	if err != nil {
		return err
	}

	return nil
}
