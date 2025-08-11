# Project 2 Starter Code

This repository contains the starter code for Project 2!

For comprehensive documentation, see the Project 2 Spec (https://cs161.org/proj2/). also pasted below:
 Secure File Sharing System
Codabot and Evanbot in a forest.

Table of contents
Story
Project 2 Policies
Getting Started
Design Overview
Skip to main content
CS161 Summer 2025
Exam Logistics
Calendar
Policies

Resources
Staff

Project 1

Project 2
Story
Project 2 Policies
Getting Started
Design Overview
Library Functions
Users And User Authentication
File Operations
Sharing and Revocation
Debugging and Errors
Advice and Tips
Appendix

Project 3
This site uses Just the Docs, a documentation theme for Jekyll.
Search CS161 Summer 2025
Dark Mode
Textbook
Ed
OH Queue
Extensions
Anonymous Feedback
Project 2	Design Overview
Functionality Overview
Threat Model
The User Class
Atomic Operations
Stateless Design
Keystore
Datastore
Threat Model: Datastore Adversary
Error Handling
Terminology Note: Pointers
Design Overview
In this section, we give a high-level overview of what you’ll be designing. We also describe general requirements that apply to the entire design. (Requirements specific to individual functions are described later.)

Functionality Overview
In this project, you will be designing a system that allows users to securely store and share files in the presence of attackers. In particular, you will be implementing the following 8 functions:

InitUser: Given a new username and password, create a new user.
GetUser: Given a username and password, let the user log in if the password is correct.
User.StoreFile: For a logged-in user, given a filename and file contents, create a new file or overwrite an existing file.
User.LoadFile: For a logged-in user, given a filename, fetch the corresponding file contents.
User.AppendToFile: For a logged-in user, given a filename and additional file contents, append the additional file contents at the end of the existing file contents, while following some efficiency requirements.
User.CreateInvitation: For a logged-in user, given a filename and target user, generate an invitation UUID that the target user can use to gain access to the file.
User.AcceptInvitation: For a logged-in user, given an invitation UUID, obtain access to a file shared by a different user. Allow the recipient user to access the file using a (possibly different) filename of their own choosing.
User.RevokeAccess: For a logged-in user, given a filename and target user, revoke the target user’s access so that they are no longer able to access a shared file.
More details about these functions later in the spec. You don’t need to implement a frontend for this project; all users will interact with your system by running a copy of your code and calling these 8 functions.

Threat Model
You could say that the high-level goal is to design a system that users can use to store their data on an untrusted server. Other users should not be able to access their data, and even if the server is malicious, it should not be able to access their data. You can assume that legitimate users will use your code, but there are malicious actors in our threat model that can use any client code they wanted, and the system will have to remain secure despite that fact. The details of these attackers are detailed later, in the sections Threat Model: Datastore Adversary and Design Requirements: Revoked User Adversary.

The User Class
You will implement the 8 functions above as part of the User class. The User class has:

Two constructors, InitUser and GetUser, which create and return new User objects. A User object is represented as a reference to a User struct containing instance variables. A constructor is called every time a user logs in.
Six instance methods (the other functions listed above), which can be called on a User object. The user calls these functions to perform file and sharing operations.
Instance variables, specific to each User object, which can be used to store information about each user. The instance variables are stored in the User struct returned by a constructor, and the instance variables are accessible to all the instance methods.
If a user calls a constructor multiple times, there will be multiple User objects that all represent the same user. You will need to make sure that these User objects do not use outdated data (more details later).

Atomic Operations
You do not need to worry about parallel function calls or any concurrency issues for this project. You can assume that at any given time, only one function is being executed. Another function call will only begin after the current function call has completed.

You can also assume that we will not quit the program in the middle of a function call. As long as your code doesn’t crash, every function call will run to completion.

You can also assume that malicious action will only happen in between function calls, and will not happen in the middle of a function call.

Stateless Design
Multiple users may use the system by calling your functions. Each user may have multiple devices (e.g. Alice may have a laptop and a phone). Every device runs a separate identical copy of your code.

As a consequence of this, your code cannot have global variables (except for basic constants). This is because these global variables will not be synced across devices.

If a user’s copy of the code crashes after running a function, or if the user quits running the code, they will lose all data stored in the code’s local memory (e.g. global variables, instance variables in User objects, etc). Therefore, you cannot rely on storing persistent information in the code’s local memory.

All devices running your code are able to send and receive data from two shared remote databases called Datastore and Keystore (more details later). All persistent data must be stored on Datastore or Keystore.

Keystore
Keystore is a remote database where you can store public keys.

Keystore is organized as a set of name-value pairs, similar to a dictionary in Python or a HashMap in Java. (Note: Sometimes these are called key-value pairs, but we will call them name-value pairs to avoid confusion with cryptographic keys.)

The name in each name-value pair must be string. The value in each name-value pair must be a public key (either a public encryption key of type PKEEncKey, or a public verification key of type DSVerifyKey).

Go’s type-checking will enforce that all values stored are public keys. You cannot store salts, hashes, structs, files, or any other data that is not a public key on Keystore. If you want to store something that is not a public key, you must store it in Datastore, not Keystore.

Each user may only store a small, constant number of public keys on Keystore. In other words, the total number of keys on Keystore should only scale with the number of users. In other words, you cannot create a new public key per file, per share, etc.

Once a name-value pair is written to Keystore, it cannot be modified or deleted. Everybody (including attackers) can read all values, but cannot modify any values, on Keystore.

Keystore is already implemented for you. You can write new name-value pairs using KeystoreSet, and you can read the value corresponding to a name using KeystoreGet (more details later).

Datastore
Datastore is an insecure remote database where you can store data. The Datastore Adversary is an attacker who can read and modify any data on Datastore (more details later). Therefore, you must protect the confidentiality and integrity of any sensitive data you store on Datastore.

Datastore is organized as a set of name-value pairs, just like Keystore. The name in each name-value pair must be a UUID, a unique 16-byte string (more details later). The value in each name-value pair can be any byte array of data.

Given a specific name (UUID), there is one and only one corresponding value, which can be read and modified by anybody who knows the name (UUID).

Datastore is already implemented for you. You can write new name-value pairs using DatastoreSet, you can read the value corresponding to a name using DatastoreGet, and you can delete a name-value pair using DatastoreDelete (more details later).

Data structures: What data structures are you going to use to organize data on Datastore?

List any struct definitions that you plan on using, along with the fields that each struct will contain.

We recommend starting with a few core data structures (e.g. struct user, struct file, struct invitation, etc.) and adding additional fields and structs as you need them.

For every subsequent design question, think about what data structures are used, and what is being stored and where. How will you generate the UUID that you are storing the information at? How will you ensure confidentiality and integrity? Which keys will you use, and how are you generating and accessing them?

Threat Model: Datastore Adversary
The Datastore Adversary is an attacker who can read and modify all name-value pairs, and add new name-value pairs, on Datastore. They can modify Datastore at any time (but not in the middle of another function executing).

The Datastore Adversary has a global view of Datastore; in other words, they can list out all name-value pairs that currently exist on Datastore.

The Datastore Adversary can take snapshots of Datastore at any time. For example, they could write down all existing name-value pairs before a user calls StoreFile. Then, they could write down all existing name-value pairs after a user calls StoreFile and compare the difference to see which name-value pairs changed as a result of the function call.

The Datastore Adversary can see when a user calls a function (e.g. if a user calls StoreFile, they know which user called it and when).

The Datastore Adversary can view and record the content and metadata of all requests to the Datastore API. This means that they will know what the inputs and outputs to the functions are.

The Datastore Adversary is not a user in the system, and will not collude with other users. However, the Datastore Adversary has a copy of your source code (Kerckhoff’s principle), and they could execute lines of your code on their own in order to modify Datastore in a way that mimics your code’s behavior.

The Datastore adversary will not perform any rollback attacks: Given a specific UUID, they will not read the value at that UUID, and then later replace the value at that UUID with the older value they read. (Deleting a value at a UUID is not a rollback attack.)

They will also not perform any rollback attacks on multiple UUIDs. For example, they will not revert the entire contents of Datastore to some previous snapshot of Datastore they took.

There is one other adversary besides the Datastore Adversary, called the Revoked User Adversary. The two adversaries do not collude. This additional adversary will be described later.

Datastore Adversary: How should you store information in the Datastore?

You might be tempted to follow a design where you hash the password, store the hashed password on the Datastore, and then GetUser will use that to check whether the user’s password is correct, and if so, allow the user to access their data. This approach isn’t going to work. Exercise: What’s wrong with this approach? Where will you store the user’s data? If it is encrypted, where will the key be stored? Can you ensure that if the Datastore Adversary tampers with the stored values, that you can detect thi

Error Handling
All 8 functions have err as one of their return values. If the function successfully executes without an error, you should return nil (the null value in Go) for the err return value.

If a function is unable to execute correctly, all you need to do is return an error that is not nil. The function could fail due to functionality issues (e.g. a user supplies an invalid argument), or security issues (e.g. an attacker has tampered with data that prevents your function from executing correctly). The error message can be anything you want (as long as the error is not nil), though we recommend using informative error messages for easier debugging.

You only need to detect when errors occur in this project; you do not need to recover from errors. For example, suppose an attacker has tampered with a file stored in Datastore, and the user calls LoadFile to try and read the file contents. Your code only needs to detect that tampering has occured and return any non-nil error. You do not need to recover from the error (i.e. you don’t need to try and recover the original file contents).

After an adversary performs malicious action, your function must either return an error, or execute correctly as if the adversary had not performed malicious action.

As soon as a non-nil error is returned, all subsequent function calls can have undefined behavior. Undefined behavior means that your code can do anything (execute without an error, return an error, crash, etc.), as long as you do not violate any security requirements.

A user stores two files, FileA.txt and FileB.txt. Your code stores the contents of these files in Datastore. Then, an attacker modifies some values on Datastore, including the contents of FileB.txt.

If the user tries to load FileB.txt, you should return a non-nil error value to detect that tampering has occurred. You do not need to try and recover the original unmodified contents of FileB.txt.

If the user instead tries to load FileA.txt, you have two options. If your code can successfully load the original unmodified contents of this file, you can return a nil error. Alternatively, if your code is unable to load the original unmodified contents of this file, you can return a non-nil error.

If the user’s call to the load function returned a non-nil error, then all subsequent function calls have undefined behavior, as long as you do not violate any security requirements (e.g. subsequent function calls still cannot leak the contents of a confidential file).

Your code should never panic (the Go term for crashing the program). You should always return a non-nil error that can be safely processed by other code (e.g. the autograder). If your code panics, then the autograder might crash and be unable to give you a score.

Terminology Note: Pointers
In this project, the term “pointer” can refer to two different concepts. In this spec, we will always clarify which concept we’re referring to.

A Go memory pointer is a variable that contains a memory address of some other object in the Go runtime’s local memory. Go memory pointers are similar to pointers in C. In this project, because you cannot store persistent data in the program’s local memory, you should rarely need to use Go memory pointers.

The only scenario where Go memory pointers are required for this project are in the InitUser and GetUser constructors. The constructors create a new User object by creating a new User struct in local memory, and returning a Go memory pointer that contains the address of this User struct in local memory.

Here’s a code snippet showing how the User class works.

var alice *client.User
var err error
alice, err = client.InitUser("Alice", "password")
alice.StoreFile("FileA.txt", []byte("Some text."))

The variable alice is a Go memory pointer, containing the address of a User struct that was created by the InitUser constructor.

StoreFile is an instance method, so you cannot call it by itself. Instead, you need to call it on an existing User object (alice in this example). This function call has access to all the instance variables in the User struct that alice points to.

Recall that in Datastore, you can store any data you want at a given UUID. If the data you choose to store is another UUID, then you’ve created a Datastore pointer. If you fetch the data at the given UUID, you will receive another UUID which references another location in Datastore.


Library Functions

Skip to main content
CS161 Summer 2025
Exam Logistics
Calendar
Policies

Resources
Staff

Project 1

Project 2
Story
Project 2 Policies
Getting Started
Design Overview
Library Functions
Users And User Authentication
File Operations
Sharing and Revocation
Debugging and Errors
Advice and Tips
Appendix

Project 3
This site uses Just the Docs, a documentation theme for Jekyll.
Search CS161 Summer 2025
Dark Mode
Textbook
Ed
OH Queue
Extensions
Anonymous Feedback
Project 2	Library Functions
Keystore
Datastore
UUID
JSON Marshal and Unmarshal
Random Byte Generator
Cryptographic Hash
Symmetric-Key Encryption
HMAC
Public-Key Encryption
Digital Signatures
Password-Based Key Derivation Function
Hash-Based Key Derivation Function
Library Functions
In this section, we provide documentation for some cryptographic functions and some utility helper functions that you can use at any time in your design. These functions have already been implemented for you in the project2-userlib library, which will be imported for you in the starter code.

Please carefully read through the provided functions while coming up with your design so that you are aware of what is possible to actually implement in code.

You cannot import any libraries besides what we’ve already imported in the starter code. You should not need any external libraries for this project.

You should not write your own cryptographic functions for this project. For example, you shouldn’t write code to implement AES-CTR yourself. Instead, you should call the existing SymEnc function that we’ve provided.

As discussed in class, you should avoid any unsafe cryptographic design patterns, such as reusing the same keys in different algorithms (see the tips section for more details), or using MAC-then-encrypt.

Helper functions: As you come up with a design, think about any helper functions you might write in addition to the cryptographic functions included here.

Having helper functions can simplify your code. Consider authenticated encryption, hybrid encryption, etc.

Keystore
userlib.KeystoreSet(name string, value PKEEncKey/DSVerifyKey) (err error)

Stores a name and value as a name-value pair into Keystore. The name can be any unique string, and the value must be a public key. You cannot store any data that is not a public key in Keystore.

Keystore is immutable: A name-value pair cannot be modified or deleted after being stored in Keystore. Any attempt to modify an existing name-value pair will return an error.

userlib.KeystoreGet(name string) (value PKEEncKey/DSVerifyKey, ok bool)

Looks up the provided name and returns the corresponding value.

If a corresponding value exists, then ok will be true; otherwise, ok will be false.

Datastore
userlib.DatastoreSet(name UUID, value []byte)

Stores name and value as a name-value pair into Datastore.

Datastore is mutable: If name already maps to an existing name-value pair, then the existing value will be overwritten with the provided value.

userlib.DatastoreGet(name UUID) (value []byte, ok bool)

Looks up the provided name and returns the corresponding value.

If a corresponding value exists, then ok will be true; otherwise, ok will be false.

userlib.DatastoreDelete(key UUID)

Looks up the provided name and deletes the corresponding value, if it exists.

UUID
Recall that in the name-value pairs of Datastore, the name should be a UUID. UUID stands for Universal Unique Identifier, and is a unique 16-byte (128-bit) value.

There are two ways to create UUIDs. You can randomly generate a new UUID from scratch. Alternatively, you can take an existing 16-byte string, and deterministically cast it into a UUID.

The uuid library also provides uuid.Nil, a UUID consisting of all zeros to represent a nil value.

uuid.New() (uuid.UUID)

Returns a randomly generated UUID.

Note: If you’re concerned about two randomly-generated UUIDs being the same, think about the probability that two randomly-generated 128-bit values are identical. In this project, you don’t have to worry about events that are astronomically unlikely to occur.

uuid.FromBytes(b []byte) (uuid UUID, err error)

Creates a new UUID by copying the 16 bytes in b into a new UUID.

Returns an error if the byte slice b does not have a length of 16.

Note: This function does not apply any additional security to the inputted byte slice. You can think of this function as casting a 16-byte value into a UUID. Anybody who reads the UUID will be able to determine what 16-byte value you used to generate the UUID, so you should not pass sensitive information into this function.

JSON Marshal and Unmarshal
Recall that in the name-value pairs of Datastore, the value should be a byte array.

If you want to store other types of data (e.g. structs) in Datastore, you will need to convert that data into a byte array before storing it. Then, you will need to convert the byte array back into the original data structure when retrieving the data.

We’ve provided the json.Marshal serialization function, which takes any arbitrary data and converts it into a byte array.

We’ve also provided the json.Unmarshal deserialization function, which takes a byte array outputted by json.Marshal, and converts it back into the original data.

json.Marshal(v interface{}) (bytes []byte, err error)

Converts an arbitrary Go value, v, into a byte slice containing the JSON representation of the struct.

If the value is a struct, only fields that start with a capital letter are converted. Fields starting with a lowercase letter are not marshaled into the output.

This function will automatically follow Go memory pointers (including nested Go memory pointers) when marshalling.

// Serialize a User struct into JSON.
type User struct {
     Username string
     Password string
     lostdata int
}

alice := &User{
 "alice",
 "password",
 42,
}

aliceBytes, err := json.Marshal(alice)
userlib.DebugMsg("%s\n", string(aliceBytes))
// {"Username":"alice","Password":"password"}

json.Unmarshal(v []byte, obj interface{}) (err)

Converts a byte slice v, generated by json.Marshal, back into a Go struct. Assigns obj to the converted Go struct.

Only struct fields that start with a capital letter will have their values restored. Struct fields that start with a lowercase letter will be initialized to their default value.

This function automatically generates nested Go memory pointers where needed to generate a valid struct.

This function will return an error if there is a type mismatch between the JSON and the struct (e.g. storing a string into a number field in a struct).

// Serialize a User struct into JSON.
// The lostdata field will NOT be included in the byte array output.
type User struct {
     Username string
     Password string
     lostdata int
}

aliceBytes := []byte("{\"Username\":\"alice\",\"Password\":\"password\"}")
 var alice User
 err = json.Unmarshal(aliceBytes, &alice)
 if err != nil { return }

userlib.DebugMsg("%v\n", alice)
// {alice password 0}

Random Byte Generator
RandomBytes(bytes int) (data []byte)

Given a length bytes, return that number of randomly generated bytes.

The random bytes returned could be used as an IV, symmetric key, or anything else you’d like.

You don’t need to worry about the underlying implementation (e.g. you don’t have to think about reseeding any PRNG). You can assume the returned bytes are indistinguishable from truly random bytes.

Cryptographic Hash
Hash(data []byte) (sum []byte)

Takes in arbitrary-length data, and outputs sum, a 64-byte SHA-512 hash of the data.

Note: you should use HMACEqual to determine hash equality. This function runs in constant time and avoids timing side-channel attacks.

Symmetric-Key Encryption
SymEnc(key []byte, iv []byte, plaintext []byte) (ciphertext []byte)

Encrypts the plaintext using AES-CTR mode with the provided 16-byte key and 16-byte iv.

Returns the ciphertext, which will contain the IV (you do not need to store the IV separately).

This function is capable of encrypting variable-length plaintext, regardless of size. You do not need to pad your plaintext to any specific block size.

SymDec(key []byte, ciphertext []byte) (plaintext []byte)

Decrypts the ciphertext using the 16-byte key.

The IV should be included in the ciphertext (see SymEnc).

If the provided ciphertext is less than the length of one cipher block, then SymDec will panic (remember, your code should always return errors, and not panic).

Notice that the SymDec method does not return an error. In other words, if some ciphertext has been mutated, SymDec will return non-useful plaintext (e.g. garbage), since AES-CTR mode does not provide integrity.

HMAC
HMACEval(key []byte, msg []byte) (sum []byte, err error)

Takes in an arbitrary-length msg, and a 16-byte key. Computes a 64-byte HMAC-SHA-512 on the message.

HMACEqual(a []byte, b []byte) (equal bool)

Compare whether two HMACs (or hashes) a and b are the same, in constant time.

If a and b are the same HMAC/hash, then equals will be true; otherwise, equals will be false.

Public-Key Encryption
PKEEncKey: A data type for RSA public (encryption) keys.

PKEDecKey: A data type for RSA private (decryption) keys.

PKEKeyGen() (PKEEncKey, PKEDecKey, err error)

Generates a 256-byte RSA key pair for public-key encryption.

PKEEnc(ek PKEEncKey, plaintext []byte) (ciphertext []byte, err error)

Uses the RSA public key ek to encrypt the plaintext, using RSA-OAEP.

PKEDec(dk PKEDecKey, ciphertext []byte) (plaintext []byte, err error)

Use the RSA private key dk to decrypt the ciphertext.

Note: RSA encryption does not support very long plaintext. If you need to use a public key to encrypt long plaintext, consider writing a helper function that implements hybrid encryption.

Recall the hybrid encryption process: Use the given public key to encrypt a random symmetric key. Then, use the symmetric key to encrypt the actual data. Return the symmetric key (encrypted with the public key) and the data (encrypted with the symmetric key).

Recall the decryption process for hybrid encryption schemes: Use the given private key to decrypt the symmetric key. Then, use the symmetric key to decrypt the data.

Digital Signatures
DSSignKey: A data type for RSA private (signing) keys.

DSVerifyKey: A data type for RSA public (verification) keys.

DSKeyGen() (DSSignKey, DSVerifyKey, err error)

Generates an RSA key pair for digital signatures.

DSSign(sk DSSignKey, msg []byte) (sig []byte, err error)

Given an RSA private (signing) key sk and a msg, outputs a 256-byte RSA signature sig.

DSVerify(vk DSVerifyKey, msg []byte, sig []byte) (err error)

Uses the RSA public (verification) key vk to verify that the signature sig on the message msg is valid. If the signature is valid, err is nil; otherwise, err is not nil.

Password-Based Key Derivation Function
Argon2Key is a slow hash function, designed specifically for hashing passwords.

Argon2Key is called a Password-Based Key Derivation Function (PBKDF) because the output (i.e. the hashed password) can be used as a symmetric key. An attacker cannot brute-force passwords to learn the key because the hash function is too slow. Also, the hash function makes the hashed password look unpredictably random, so it can be used as a symmetric key.

You can assume that the user’s chosen password has sufficient entropy for the PBKDF output to be used as a symmetric key.

The salt argument is used to ensure that two users with the same password don’t have the same password hash. If you choose to use the hash as a key, then the salt also ensures that the two users don’t use the same key.

Argon2Key(password []byte, salt []byte, keyLen uint32) (result []byte)

Applies a slow hash to the given password and salt. The outputted hash is keyLen bytes long, and can be used as a symmetric key.

Hash-Based Key Derivation Function
You can use the HashKDF to deterministically derive multiple keys from a single root key. This can simplify your key management schemes.

HashKDF is a fast hash function, similar to HMAC, that essentially hashes the source key and the purpose together. Changing either the source key, or the purpose, or both, will cause the output of HashKDF to be unpredictably different.

One way you can use HashKDF is by calling it multiple times with the same source key but different, hard-coded purposes every time. This will generate multiple keys, one per call to HashKDF. Anybody who knows the source key and the purposes can re-generate the keys by calling HashKDF again (without needing to store the derived keys). Anybody who doesn’t know the source key will be unable to generate the keys.

If the source key is insecure (e.g. an attacker knows its value), and the purpose is insecure (e.g. it’s a hard-coded string and the attacker has a copy of your code), then the derived keys outputted by HashKDF will also be insecure.

HashKDF(sourceKey []byte, purpose []byte) (derivedKey []byte, err error)

Hashes together a 16-byte sourceKey and some arbitrary-length byte array purpose to deterministically derive a new 64-byte derivedKey.

If you don’t need all 64 bytes of the output, you can slice to obtain a key of the desired length.

Here’s a code snippet showing how you could use HashKDF to take one source key and derive two keys, one for encryption and one for MACing.

sourceKey := userlib.RandomBytes(16)

encKey, err := userlib.HashKDF(sourceKey, []byte("encryption"))
if err != nil { return }

macKey, err := userlib.HashKDF(sourceKey, []byte("mac"))
if err != nil { return }


Users And User Authentication

Skip to main content
CS161 Summer 2025
Exam Logistics
Calendar
Policies

Resources
Staff

Project 1

Project 2
Story
Project 2 Policies
Getting Started
Design Overview
Library Functions
Users And User Authentication
File Operations
Sharing and Revocation
Debugging and Errors
Advice and Tips
Appendix

Project 3
This site uses Just the Docs, a documentation theme for Jekyll.
Search CS161 Summer 2025
Dark Mode
Textbook
Ed
OH Queue
Extensions
Anonymous Feedback
Project 2	Users And User Authentication
InitUser
GetUser
Design Requirements: Usernames and Passwords
Design Requirements: Multiple Devices
Users And User Authentication
In this section, you’ll design two constructors to support creating new users and letting users log in to the system.

This example scenario illustrates how to create new users and let existing users log in.

EvanBot calls InitUser("evanbot", "password123").
This creates a new user with username “evanbot” and password “password123”. If the username “evanbot” already exists, the function would return an error.
This constructor function creates and returns a User object with instance variables. EvanBot can call the instance methods of this object to perform file operations.
There is no log out operation. If EvanBot is done running file operations, they can simply quit the program, which will destroy the User object (and its instance variables). This should not cause any data to be lost.
Later, EvanBot runs your code again and calls GetUser("evanbot", "password123").
This constructor function should create and return a User object corresponding to the existing EvanBot user. As before, the object can have instance variables, and EvanBot can call the instance methods to perform file operations.
If the password is incorrect, the function would return an error.
CodaBot calls InitUser("codabot", "password123").
The function should create and return a User object corresponding to the new CodaBot user.
Note that different users could choose the same password.
InitUser
InitUser(username string, password string) (userdataptr *User, err error)

Constructor of the User class, used when a user logs in for the first time.

Creates a new User object and returns a Go memory pointer to the new User struct. The User struct can include any instance variables for the corresponding User object.

Recall that when the program quits, the User struct and all its data will be lost. Only data on Datastore and Keystore persists after the program quits.

Returns an error if:

A user with the same username exists.
An empty username is provided.
GetUser
GetUser(username string, password string) (userdataptr *User, err error)

Constructor of the User class, used when an existing user logs in.

Creates a User object for a user who has already been initialized with InitUser, and returns a Go memory pointer to it.

If you stored the data of the User struct to Datastore in InitUser, then you could download that data from Datastore and create a new User object in local memory using that data.

Returns an error if:

There is no initialized user for the given username.
The user credentials are invalid.
The User struct cannot be obtained due to malicious action, or the integrity of the user struct has been compromised.
User Authentication: How will you authenticate users?

How does EvanBot create a new user? When they try to log in, how will you check if the username and password are correct? How do you ensure that a user cannot login with an incorrect password?

Design Requirements: Usernames and Passwords
Usernames:

Each user has a unique username.
Usernames are case-sensitive: Bob and bob are different users.
Usernames can be any string with 1 or more characters (not necessarily alphanumeric).
Passwords:

Different users might choose to use the same password.
The passwords provided by users have sufficient entropy for the PBKDF slow hash function to output an unpredictable string that an attacker cannot guess by brute force.
The passwords provided by users do not have sufficient entropy to resist brute-force attacks on any of the other fast hash functions (Hash, HashKDF, or HMAC).
Passwords can be any string with 0 or more characters (not necessarily alphanumeric, and could be the empty string).
Design Requirements: Multiple Devices
Users must be able to create multiple User instances on different devices. In other words, a user should be able to call GetUser multiple times, with the same username and password, to obtain multiple different copies of the User struct on multiple different devices.

All changes to files made from one device must be reflected on all other devices immediately (i.e. a user should not have to call GetUser again to see the changes).

This example scenario illustrates how users should be able to create multiple User instances on multiple devices:

EvanBot has a copy of the system’s code running on their laptop. EvanBot has another, duplicate copy of the system’s code running on their phone.
On the laptop, EvanBot calls GetUser("evanbot", "password").
The system creates a User object in the laptop’s local memory. We’ll denote this object as evanbot-laptop.
Without terminating the code running on the laptop, EvanBot calls GetUser("evanbot", "password") on their phone.
The system creates another User object in the phone’s local memory. We’ll denote this object as evanbot-phone.
evanbot-laptop and evanbot-phone are two different User structs. They exist on two different devices, and they both correspond to the same user (EvanBot).
On the laptop, EvanBot calls evanbot-laptop.StoreFile("toppings.txt", "syrup").
On the phone, EvanBot calls evanbot-phone.LoadFile("toppings.txt") and sees “syrup”.
Note that duplicate user objects, running on separate devices, should be able to see the latest updates to files.
On the phone, EvanBot calls evanbot-phone.AppendToFile("toppings.txt", " and butter").
On the laptop, EvanBot calls evanbot-laptop.LoadFile("toppings.txt") and sees “syrup and butter”.
It would be incorrect behavior if the system returned “syrup”, because this means the append from the other device was not properly synced.
Multiple Devices: How will you ensure that multiple User objects for the same user always see the latest changes reflected?

EvanBot logs in on their laptop and phone, creating two User objects. If EvanBot makes a change on their laptop (e.g. storing a file), how do you ensure that EvanBot will see the change reflected on their phone?


File Operations

Skip to main content
CS161 Summer 2025
Exam Logistics
Calendar
Policies

Resources
Staff

Project 1

Project 2
Story
Project 2 Policies
Getting Started
Design Overview
Library Functions
Users And User Authentication
File Operations
Sharing and Revocation
Debugging and Errors
Advice and Tips
Appendix

Project 3
This site uses Just the Docs, a documentation theme for Jekyll.
Search CS161 Summer 2025
Dark Mode
Textbook
Ed
OH Queue
Extensions
Anonymous Feedback
Project 2	File Operations
Design Requirements: Namespacing
Design Requirements: Files
User.StoreFile
User.LoadFile
User.AppendToFile
Design Requirements: Bandwidth & Append Efficiency
File Operations
In this section, you’ll design three instance methods to support creating new files or overwriting the contents of existing files, reading file contents, and appending content to the end of existing files.

Design Requirements: Namespacing
Note that different users can have files with the same name. A user’s namespace is defined as all of the filenames they are using. One user’s namespace could contain a filename that another user is also using. In that other user’s namespace, that same filename could refer to a different file (or the same file, if it was shared - more details about sharing later).

This example scenario illustrates how file storage and namespacing works.

EvanBot calls StoreFile("foods.txt", "pancakes").
Assuming that EvanBot has never stored to foods.txt before, this creates a new file called foods.txt in EvanBot’s personal namespace.
EvanBot calls LoadFile("foods.txt") and sees “pancakes”.
EvanBot calls StoreFile("foods.txt", "cookies").
Because foods.txt is an existing file, this call should overwrite the entire file with the new contents.
EvanBot calls LoadFile("foods.txt") and sees “cookies”.
EvanBot calls LoadFile("drinks.txt") and sees an error, because there is no file named drinks.txt in EvanBot’s personal namespace.
EvanBot calls AppendToFile("foods.txt", " and pancakes").
Instead of overwriting the entire file, this should append additional contents to the end of an existing file.
EvanBot calls LoadFile("foods.txt") and sees “cookies and pancakes”.
EvanBot calls AppendToFile("foods.txt", " and hash browns").
EvanBot calls LoadFile("foods.txt") and sees “cookies and pancakes and hash browns”.
EvanBot calls StoreFile("foods.txt", "pancakes").
This overwrites the entire file (including appends) with the new contents.
EvanBot calls LoadFile("foods.txt") and sees “pancakes”.
EvanBot calls AppendToFile("drinks.txt", "and cookies") and sees an error, because there is no file named drinks.txt in EvanBot’s personal namespace.
CodaBot calls StoreFile("foods.txt", "waffles").
Note that this creates a new file in CodaBot’s personal namespace named foods.txt. This should not interfere with the foods.txt file in EvanBot’s namespace, which is a different file.
CodaBot calls LoadFile("foods.txt") and sees “waffles”.
EvanBot calls LoadFile("foods.txt") and sees “pancakes”.
Design Requirements: Files
Confidentiality of data:

You must ensure that no information is leaked about these 3 pieces of data:
File contents for all files.
Filenames for all files.
The length of the filenames for all files.
You must also ensure that no information is leaked that could be directly or indirectly used to learn these 3 pieces of data.
For example, if you have a secret key that you’re using to encrypt some file contents, you’ll need to ensure that secret key is not leaked either.
You may leak information about any other values besides the ones listed above.
For example: It’s okay if an adversary learns usernames, length of a file, how many files a user has, etc.
Integrity of data:

You must be able to detect when an attacker has tampered with the contents of a file.
Filenames:

Filenames can be any string with 0 or more characters (not necessarily alphanumeric, and could be the empty string).
Different users can have files with the same filename, but they could refer to different files.
File Storage and Retrieval: How does a user store and retrieve files?

EvanBot is logged in and stores a file. How does the file get stored in Datastore? What key(s) and UUIDs do you use? How do you access the file at a later time?

User.StoreFile
This instance method is used to both create a file for the first time, or to overwrite an existing file entirely with new contents. To use this method, the user passes in the filename to identify the file, as well as the contents that they wish to store.

User.StoreFile(filename string, content []byte) (err error)

Given a filename in the personal namespace of the caller, this function persistently stores the given content for future retrieval using the same filename.

If the given filename already exists in the personal namespace of the caller, then the content of the corresponding file is overwritten.

The client must allow content to be any arbitrary sequence of bytes, including the empty sequence.

Note that calling StoreFile after malicious tampering has occurred is undefined behavior, and will not be tested.

Note that calling StoreFile on a file whose access has been revoked is undefined behavior, and will not be tested.

User.LoadFile
User.LoadFile(filename string) (content []byte, err error)

Given a filename in the personal namespace of the caller, this function downloads and returns the content of the corresponding file.

Note that, in the case of sharing files, the corresponding file may or may not be owned by the caller.

Returns an error if:

The given filename does not exist in the personal file namespace of the caller.
The integrity of the downloaded content cannot be verified (indicating there have been unauthorized modifications to the file).
Loading the file cannot succeed due to any other malicious action.
User.AppendToFile
User.AppendToFile(filename string, content []byte) (err error)

Given a filename in the personal namespace of the caller, this function appends the given content to the end of the corresponding file.

content can be any arbitrary sequence of 0 or more bytes.

Note that, in the case of sharing files, the corresponding file may or may not be owned by the caller.

You are not required to check the integrity of the existing file before appending the new content (integrity verification is allowed, but not required).

Returns an error if:

The given filename does not exist in the personal file namespace of the caller.
Appending the file cannot succeed due to any other malicious action.
Design Requirements: Bandwidth & Append Efficiency
All functions except for AppendToFile have no efficiency requirements, as long as they don’t time out the autograder. You can submit your code to Gradescope to check that you aren’t timing out the autograder (~20 minutes).

The efficiency requirement for appending is measured in terms of bandwidth, not in terms of time complexity or space complexity, which you may have seen in other classes. This means that your append can use unlimited amounts of local compute (e.g. you can encrypt and decrypt as much data as you’d like).

Recall that DataStore and KeyStore are remote databases. This means that when you call DataStoreGet, you are downloading all data at the specified UUID from DataStore to the local device running your code. Similarly, when you call DataStoreSet, you are uploading all the specified data from your local device running your code to DataStore. The only efficiency requirement forAppendToFile is that the total amount of data uploaded with calls to DataStoreSet and downloaded with calls to DataStoreGet must be efficient.

The bandwidth used by a call to AppendToFile is defined as the total size of all data in calls to DataStoreSet and DataStoreGet. All calls that are not DataStoreSet or DataStoreGet do not affect the total bandwidth.

The total bandwidth should only scale with the size of the append (i.e. the number of bytes in the content argument to AppendToFile). In other words, if you are appending n bytes to the file, it’s okay (and unavoidable) that you’ll need to upload at least n bytes of data to the Datastore.

Your total append bandwidth can additionally include some small constant factor. We cannot reveal the exact number, but an example of a reasonable constant would be 3,000 bytes on every call to append.

The total bandwidth should not scale with (including but not limited to):

Total file size
Number of files
Length of the filename
Number of appends
Size of previous append
Length of username
Length of password
Number of users the file is shared with
Here is one way to consider whether your design scales with the number of appends. Suppose we call AppendToFile on a file 10,000 times, appending 1 byte every time. The 1,000th and 10,000th call to AppendToFile should use the same total bandwidth as the 1st append operation. This should also be true for the case of appending an arbitrary number of bytes repetitively, 10,000 times (even 0!). Specifically, when a scheme scales with the number of appends, this means that the amount of bandwidth used is directly proportional to the number of appends that have occurred, not just the fact that appends have happened. This means if someone appended 0 bytes to a file 10,000 times (called AppendToFile with an empty string passed in), the 1,000th and 10,000th call to AppendToFile should use the same total bandwidth as the 1st. If the bandwidth increases even though the string appended is always the empty string (and hence, nothing actually changed about the file), then this design would scale with the number of appends.

Here is one way to consider whether your design scales with the size of the previous append. Suppose we call AppendToFile to append 1 terabyte of data to a file. Then, we call AppendToFile again on the same file to append another 100 bytes. The total bandwidth of the second call to append should not include the 1 terabyte of bandwidth from the previous (first) append.

In general, one way to check for efficiency is to imagine a graph where the x-axis is the potential scaling factor (e.g. file size), and the y-axis is the total bandwidth. The plot of scaling factor vs. total bandwidth should be a flat line, not an upwards sloping line.

As an analogy, imagine that the users of this system have a limited phone data plan. We want to avoid excessive charges to their data plan, so we want to avoid downloading or uploading unnecessary data when appending.

For example, a naive implementation would involve:

The user calls DataStoreGet to download the entire contents of the file.
The user decrypts the file locally.
The user appends the contents locally.
The user encrypts the entire file contents.
The user calls DataStoreSet to upload the entire file to DataStore.
Note for steps 2 & 4: These parts do not count against bandwidth efficiency. Recall, only DataStoreGet and DataStoreSet count for bandwidth calculation, and local computations do not count against efficiency requirements.

This implementation is inefficient because in step 1, the call to DataStoreGet downloads the entire file. This implementation is additionally inefficient due to step 5, where we call DataStoreSet and upload the entire file contents.

For example, if we had a 10 terabyte file, and we wanted to append 100 bytes to the file, the implementation above would have a total bandwidth of 20 terabytes + 100 bytes. An efficient implementation would use 100 bytes of bandwidth (possibly plus some constant).

Efficient Append: What is the total bandwidth used in a call to append?

List out every piece of data that you need to upload (DatastoreSet) or download (DatastoreGet) from Datastore in a call to append, and the size of each piece of data. Is the total a constant, or does it scale?


Sharing and Revocation

Skip to main content
CS161 Summer 2025
Exam Logistics
Calendar
Policies

Resources
Staff

Project 1

Project 2
Story
Project 2 Policies
Getting Started
Design Overview
Library Functions
Users And User Authentication
File Operations
Sharing and Revocation
Debugging and Errors
Advice and Tips
Appendix

Project 3
This site uses Just the Docs, a documentation theme for Jekyll.
Search CS161 Summer 2025
Dark Mode
Textbook
Ed
OH Queue
Extensions
Anonymous Feedback
Project 2	Sharing and Revocation
Design Requirements: Sharing and Revoking
User.CreateInvitation
User.AcceptInvitation
User.RevokeAccess
Design Requirements: Revoked User Adversary
Sharing and Revocation
In this section, you’ll design three instance methods to support sharing files with other users and revoking file access from other users.

This example scenario illustrates how file sharing occurs.

EvanBot calls StoreFile("foods.txt", "eggs").
Assuming that foods.txt did not previously exist in EvanBot’s file namespace, this creates a new file named foods.txt in EvanBot’s namespace.
Because EvanBot created the new file with a call to StoreFile, EvanBot is the owner of this file.
EvanBot calls CreateInvitation("foods.txt", "codabot").
This function returns a UUID, which we’ll call an invitation Datastore pointer.
The invitation UUID can be any UUID you like. For example, you could collect/compute any values that you want to send to the recipient user for them to access the file. Then, you could securely store these values on Datastore at some UUID, and return that UUID.
EvanBot uses a secure communication channel (outside of your system) to deliver the invitation UUID to CodaBot. Using this secure channel, CodaBot receives the identity of the sender (EvanBot) and the invitation UUID generated by EvanBot.
CodaBot calls AcceptInvitation("evanbot", invitationPtr, "snacks.txt").
CodaBot passes in the identity of the sender and the invitation UUID generated by EvanBot.
CodaBot also passes in a filename (snacks.txt here). Note that CodaBot (the recipient user) can choose to give the file a different name while accepting the invitation.
CodaBot calls LoadFile("snacks.txt") and sees “eggs”.
Note that CodaBot refers to the file using the name they specified when they accepted the invitation.
EvanBot calls LoadFile("foods.txt") and sees “eggs”.
Note that different users can refer to the same file using different filenames.
EvanBot calls AppendToFile("foods.txt", "and bacon").
CodaBot calls LoadFile("snacks.txt") and sees “eggs and bacon”.
Note that all users should be able to see modifications to the file.
Design Requirements: Sharing and Revoking
File access

The owner of a file is the user who initially created the file (i.e. with the first call to StoreFile).
The owner must always be able to access the file. All users who have accepted an invitation to access the file (and who have not been revoked) must also be able to access the file. These users must be able to:
Read the file contents with LoadFile.
Overwrite the file contents with StoreFile.
Append to the file with AppendToFile.
Share the file with CreateInvitation.
If a user changes the file contents, all users with access must immediately see the changes. The next time they try to access the file, all users with access should see the latest version.
All users should be reading and modifying the same copy of the file. You may not create copies of the file.
User.CreateInvitation
User.CreateInvitation(filename string, recipientUsername string) (invitationPtr UUID, err error)

Generates an invitation UUID invitationPtr, which can be used by the target user recipientUsername to gain access to the file filename.

The invitation UUID invitationPtr can be any UUID value you like. For example, you could collect/compute any values that you want to send to the recipient user for them to access the file. Then, you could securely store these values on Datastore at some UUID, and return that UUID.

The recipient user will not be able to access the file (e.g. load, store) until they call AcceptInvitation, where they will choose their own (possibly different) filename for the file.

If the target user already has access to the file, or if the target user has already had their access to the file revoked, then this function has undefined behavior and will not be tested.

Returns an error if:

The given filename does not exist in the personal file namespace of the caller.
The given recipientUsername does not exist.
Sharing cannot be completed due to any malicious action.
User.AcceptInvitation
User.AcceptInvitation(senderUsername string, invitationPtr UUID, filename string) (err error)

Accepts an invitation by inputting the username of the sender (senderUsername), and the invitation UUID (invitationPtr) that the sender previously generated with a call to CreateInvitation.

You can assume that after the sender generates an invitationPtr UUID, the sender uses a secure communication channel (outside of your system) to deliver that UUID to the recipient. Using this secure channel, the recipient receives the UUID and the sender’s username, and can input them into AcceptInvitation.

Allows the recipient user to choose their own filename for the shared file in their own file namespace. The recipient could choose to give the file a different name than what the sender named the file.

After calling this function, the recipient user should be able to perform all operations (load, store, append, create invitation) on the shared file, using their own chosen filename.

This function has undefined behavior and will not be tested if the user passes in an invitation UUID that has already been accepted.

Returns an error if:

The user already has a file with the chosen filename in their personal file namespace.
Something about the invitationPtr is wrong (e.g. the value at that UUID on Datastore is corrupt or missing, or the user cannot verify that invitationPtr was provided by senderUsername).
The invitation is no longer valid due to revocation.
File Sharing: What gets created on CreateInvitation, and what changes on AcceptInvitation?

EvanBot (the file owner) wants to share the file with CodaBot. What is stored in Datastore when creating the invitation, and what is the UUID returned? What values on Datastore are changed when CodaBot accepts the invitation? How does CodaBot access the file in the future?

CodaBot (not the file owner) wants to share the file with PintoBot. What is the sharing process like when a non-owner shares? (Same questions as above; your answers might be the same or different depending on your design.)

User.RevokeAccess
User.RevokeAccess(filename string, recipientUsername string) (err error)

Revokes access to filename from the target user recipientUsername, and all the users that recipientUsername shared the file with (either directly or indirectly).

The owner of the file must be able to call this function. This function has undefined behavior and will not be tested if the user calling it is not the owner of the file.

The owner can only call this function on a user the owner directly shared the file with. This function has undefined behavior and will not be tested if the target user is not someone the owner directly shared the file with.

Note: This function could be called either before or after the target user calls AcceptInvitation. Your code should be able to revoke access either way.

After revocation, the revoked users should not be able to access the file. They should get an error if they try to access the file using your system (e.g. by calling LoadFile, AppendToFile, etc.). Note that calling StoreFile on a revoked file is undefined behavior and will not be tested.

The revoked users should also be unable to regain access to the file, even if they try to maliciously bypass your system and directly access Datastore.

Non-revoked users should be able to continue accessing the file (e.g. by calling LoadFile, StoreFile, etc.), without needing to re-accept any invitations.

Returns an error if:

The given filename does not exist in the caller’s personal file namespace.
The given filename is not currently shared with recipientUsername.
Revocation cannot be completed due to malicious action.
This example scenario illustrates revocation behavior.

Consider the following sharing tree structure. (An edge from A to B indicates that A shared the user with B.)

A tree with nodes ABCDEFG.

A calls RevokeAccess(file, B).
This call is defined because A is the owner. Any other user calling revoke is undefined behavior and will not be tested.
This call is defined because A directly shared the file with B. A can only call revoke on B and C; revoking on any other user is undefined behavior and will not be tested.
Users B, D, E, and F should all lose access to the file.
If any of these users try to access the file (load, append, create invitation), the function should error.
These users can now be malicious: they can use values they’ve previously written down, and access Datastore (without listing out all UUIDs). However, they should still be unable to read the file, modify the file, or deduce when future updates are happening.
Users C and G can continue accessing the file, without re-accepting any invitation.
For example, C.LoadFile(file) should work without re-accepting an invitation.
Design Requirements: Revoked User Adversary
Once a user has their access revoked, they become a malicious user, who we’ll call the Revoked User Adversary. The Revoked User Adversary will not collude with any other users, and they will not collude with the Datastore Adversary.

The Revoked User Adversary’s goal is to re-obtain access to the file. The revoked user will not perform malicious actions on other files that they still have access to. Their only goal is to re-obtain access to the file that they lost access to.

The Revoked User Adversary might attempt to re-obtain access by calling functions with different arguments (e.g. calling AcceptInvitation again).

The Revoked User Adversary may also try to re-obtain access by calling DatastoreGet and DatastoreSet and maliciously affecting Datastore. However, unlike the Datastore Adversary, they do not have a global view of Datastore (i.e. they cannot list all UUIDs that have been in use).

The Revoked User Adversary will not perform any rollback attacks: Given a specific UUID (that they previously had access to), they will not read the value at that UUID, and then later replace the value at that UUID with the older value they read. They will also not perform any rollback attacks on multiple UUIDs.

Prior to having their access revoked, the Revoked User Adversary could have written down any values that they have previously seen. The Revoked User Adversary has a copy of your code running on their local computer, so they could inspect the code and learn the values of any variables that you computed.

Your code should ensure that the Revoked User Adversary is unable to learn anything about any future writes or appends to the file (learning about the file before they got revoked is okay). For example, they cannot know what the latest contents of the file are, and they should be unable to make modifications to the file without being detected. Also, they cannot know when future updates are happening (e.g. they should not be able to deduce how many times the file has been updated in the past day).

File Revocation: What values need to be updated when revoking?

Using the diagram above as reference, suppose A revokes B’s access. What values in Datastore are updated? How do you ensure C and G still have access to the file? How do you ensure that B, D, E, and F lose access to the file?

How do you ensure that a Revoked User Adversary cannot read or modify the file without being detected, even if they can directly access Datastore and remember values computed earlier? How do you ensure that a Revoked User Adversary cannot learn about when future updates are happening?


Debugging and Errors

Skip to main content
CS161 Summer 2025
Exam Logistics
Calendar
Policies

Resources
Staff

Project 1

Project 2
Story
Project 2 Policies
Getting Started
Design Overview
Library Functions
Users And User Authentication
File Operations
Sharing and Revocation
Debugging and Errors
Advice and Tips
Appendix

Project 3
This site uses Just the Docs, a documentation theme for Jekyll.
Search CS161 Summer 2025
Dark Mode
Textbook
Ed
OH Queue
Extensions
Anonymous Feedback
Project 2	Debugging and Errors
Debugging in VSCode
Debugging a Single Test Case
Errors While Debugging
Common Gradescope Errors
Debugging and Errors
This section does not contain any design requirements (i.e. you could complete the whole project without reading this section). However, we’ve compiled our general guidelines for design, development, and testing.

Debugging in VSCode
Once the extensions are installed we can now follow the debugging flow:

Set breakpoints. To set a breakpoint in your code, Navigate to the line where you want your breakpoint to be, and click on the left side of the line number. You should see a red dot appearing next to the line number, which indicates that a breakpoint has been set.
Setting a breakpoint.

Run a test. To debug a test case and peek around the breakpoint, navigate to the client_test/client_test.go file, and click on the debug test button above the test case that you want to run. Immediately after, the debugger will start. If your debugger doesn’t pause at breakpoints you’ve set, this means that your code flow never went through any lines of code that you’ve set a breakpoint on.
Debug test.

Navigate the debugger. The golang debugger has a really convenient and powerful interface. There are a couple sections you should especially be aware of:
You can step through your code using the menu bar at the center top of the screen just above the code editor (outlined in blue in the image below). Hover your mouse over each function to see the keyboard shortcut for each of these. This is a very important feature, so get familiar with each and every button.
You can use the local variables at the top left quadrant of your screen (dash outlined in red in the image below), which displays the variable name and their values. For nested structures, you can click the expand button to view variables inside the struct.
You can use the watch section at the middle left quandrant of your screen (circled in green in the image below) to constantly evaluate golang expressions. Some of the thing you can use it for is to constant check the length of an array, which you can do with call len([variable]). To constantly evaluate functions, you need to append the function call with call for watchpoints.
You can check the call stack in the call stack section, and edit breakpoints in the breakpoints section.
Using the debugger.

Debugging a Single Test Case
To learn how to debug a single test case, refer to our section on how to run a single test case.

Errors While Debugging
Here are solutions for a couple errors that you may run into while running the debugger.

Couldn't find dlv at the Go tools path or couldn't start dlv dap: Error cannot find delve debbuger pops up when you click the Debug test button, but running go test -v works fine. This seems to be a problem due to recent updates to the Go extension on VSCode (and its compatability with delve). There are a couple possible solutions:
To solve the problem solely in VSCode, run Go: Install/Update Tools from the Command Palette. You can access the Command Palette by running Cmd + Shift + P in Mac or Ctrl + Shift + P in Linux/Windows. Then mark dlv and dlv-dap from the menu, and hit okay. This should start the update. You may need to restart VSCode after doing this.
For macOS users, you can also just run brew install delve.
A general solution for other problems is making sure that you’re running Go >=v1.20. This is the suggested version of Go (you can check your version of Go by running go version in your terminal). For some people, you’ll need to download the most up-to-date Go version.
Common Gradescope Errors
Tests failed to compile! Grading cannot continue. but everything works when you run go test locally: This means that tests in your client_test.go file utilize struct attributes or helper functions (not the core API functionality). However, since these tests are run against the staff implementation, you cannot use struct attributes in your client_test.go (otherwise, you’d be assuming what structs in the staff implementation look like).
Fix: If you want tests that check against struct attributes or helper function functionality, put these in the unit test file! he separation between unit and client tests is that unit tests test for your implementation correctness and can be implementation-specific (so you can access any extra client helper functions or struct attributes) while client/integration tests are implementation-blind, just checking that the overall functionality is correct.
The distinction between the client_test.go integration tests and any tests you write in client_unittest.go as unit tests is that integration tests should pass on anyone’s implementation whereas unit tests may be specific to your design. For integration tests, we want you to write tests to ensure the correct functionality of the client API that would hold under your implementation, the staff implementation, Evanbot’s implementation, etc. These are functionality and security tests (so you can write tampering tests too!), and shouldn’t be dependent on what structs the staff solution has.
Warning: Please remove all FSpecify statements when you submit to the autograder (otherwise, the autograder will only run tests labelled FSpecify) or potentially break the autograder. We cannot promise we will re-run the autograder for you if you forgot to remove them.


Advice and Tips
Skip to main content
CS161 Summer 2025
Exam Logistics
Calendar
Policies

Resources
Staff

Project 1

Project 2
Story
Project 2 Policies
Getting Started
Skip to main content
CS161 Summer 2025
Exam Logistics
Calendar
Policies

Resources
Staff

Project 1

Project 2
Story
Project 2 Policies
Getting Started
Design Overview
Library Functions
Users And User Authentication
File Operations
Sharing and Revocation
Debugging and Errors
Advice and Tips
Appendix

Project 3
This site uses Just the Docs, a documentation theme for Jekyll.
Search CS161 Summer 2025
Dark Mode
Textbook
Ed
OH Queue
Extensions
Anonymous Feedback
Project 2	Getting Started
Design Workflow
Coding Workflow
Development Environment
Visual Studio Code (VSCode)
Getting Started Coding
Testing with Ginkgo
Basic Usage
Asserting Expected Behavior
Organizing Tests
Running a Single Test Case
Optional: Measure Local Test Coverage
Getting Started
This section does not contain any design requirements (i.e. you could complete the whole project without reading this section). However, we’ve compiled our general guidelines for design, development, and testing.

Design Workflow
This project has a lot of moving parts, and it’s normal to feel overwhelmed by the amount of requirements you need to satisfy in your design. Here is one suggestion for how you can break down this project into more manageable steps:

Read through the entire spec. It’s easy to miss a design requirement, which might cause you trouble later when you have to redo your design to meet the requirement you missed. We suggest reading through the entire spec front-to-back at least twice, just to make sure that you have internalized all the design requirements.
Design each section in order. Start with user authentication: How do you ensure that users can log in? Focus on getting InitUser and GetUser properly designed, and don’t worry about file storage or sharing yet. Then, after you’re satisfied with your login functionality, move on to the file storage functions. Don’t worry about sharing yet, and just make sure that a single user is able to LoadFile and StoreFile properly. Then, you can move on to AppendToFile. Finally, once you’re satisfied with your single-user storage design, you can move on to sharing and revoking.
Don’t be afraid to redesign. It’s normal to change your design as you go. In particular, if you follow the order of functions in the spec, then AppendToFile might result in changes to LoadFile and StoreFile. Also, RevokeAccess might result in changes to CreateInvitation and AcceptInvitation. It’s easier to change your design while you’re in the design phase; by contrast, it’s harder to change your design after you’ve already implemented it in code.
Coding Workflow
Stay organized with helper functions. If you fit all your code in 8 functions, it’s easy for the functions to get bloated and hard to debug. By contrast, if you organize your code into helper functions, you can reuse code without needing to copy-paste code blocks, and you can also write unit tests to check that each helper function is working as intended.
Test as you go. Don’t write a huge chunk of code and then test it at the end. This usually results in a failed test, and now you have no idea which part of the giant code block is broken. Instead, write small pieces of code incrementally, and write tests as you go to check that your code is doing what it’s supposed to.
Don’t split the coding between partners. Sometimes, a 2-person project group will try to have each group member independently write half of the functions. As a design-oriented project, the code in different functions will often be connected in subtle ways, and it is difficult (if not impossible) to write code without understanding all the code that has been written so far. A better approach is to work together to figure out the high-level organization of your code. Ideally, you’d use a technique like pair programming to ensure that both partners understand the code being written. The only scenario where writing code individually might be useful is for isolated helper functions, where the behavior is clearly documented and the function can be tested and debugged in isolation. Staff are not responsible for helping you understand code that your partner wrote.
Development Environment
Visual Studio Code (VSCode)
VSCode is a very commonly used IDE, and provides a powerful set of code and debugging environments that can be exploited for Golang Projects. To setup VSCode for this project, follow these steps:

Install Golang. Make sure you have Golang installed before starting this guide.
Install the GoLang extension. In the “Extensions” tab (Use Ctrl+Shift+X to navigate you can’t find it), search up the Go extension that is created by Google.
Install the debugging environment Once the extension is installed, the lower right corner might and most likely will pop up a warning stating that analysis tools might be missing. If so, click on the install, and wait for the analysis tools to install. If you missed this the first time, press (Ctrl+Shift+P) and search up “Extensions: Install Missing Dependencies,” and follow the instructions.
Getting Started Coding
After your design review, you’re ready to start implementing your design in code. Follow these steps to get started:

Install Golang. Go v1.20 is recommended.
Complete the online Golang Tutorial.
The tutorial can take quite a bit of time to complete, so plan accordingly.
The tutorial is a helpful tool that you may end up referencing frequently, especially while learning Go for the first time.
Accept the Project 2 GitHub Classroom Invite Link.
At this step, you may receive an email asking you to join the cs161-students organization.
Enter a team name. If you’re working with a partner, only one partner should create a team - the other partner should join the team through the list of teams.
Clone your repository using the git clone command.
Feel free to review 61B’s git resources (commands, commands continued, guide, common issues) for a refresher!
In README.md, make sure to include the student IDs and emails of both you and your partner, as well as your team GitHub repository link.
In the client_test directory of the checked out repository, run go test.
Go will automatically retrieve the dependencies defined in go.mod and run the tests defined in client_test.go and client_unittest.go.
It is expected that some tests will fail because you have not yet implemented all of the required functions.
Optionally, we’ve provided a unit test framework that you can access in the client directory, where you can create unit tests for your implementation-specific functions inside client_unittest.go.
If you would like to only run unit tests, please rename the unit test file to client_unit_test.go (and an _ between unit and test) and run go test in the client directory. Reminder: If you later want to run unit tests with client tests together, make sure to change the name of the file back to client_unittest.go and run go test from the client_test directory.
If the starter code is buggy and you need to pull updates from the starter code repo, you can do so with these steps:

Run only once:
If you use HTTP: git remote add starter https://github.com/cs161-staff/project2-starter-code.git
If you use SSH: git remote add starter git@github.com:cs161-staff/project2-starter-code
Run each time you need to pull updates: git pull starter main
Please refer to our Troubleshooting Tips for an in-depth debugging guide!

Testing with Ginkgo
This section provides some basic documentation for Ginkgo, which is the framework you’ll be using to write your own test cases.

First, we recommend reading through the basic tests in client_test.go, especially the first few ones, since those are well-documented in-line. Then, come back to this documentation.

Basic Usage
You should be able to write most of your tests using some combination of calls to –

Initialization methods (e.g. client.InitUser, client.GetUser)
User-specific methods (e.g. alice.StoreFile, bob.LoadFile)
Declarations of expected behavior (e.g. Expect(err).To(BeNil()))
Asserting Expected Behavior
To assert expected behavior, you may want to check (a) that an error did or didn’t occur, and/or (b) that some data was what you expected it to be. For example:

// Check that an error didn't occur
alice, err := client.InitUser("alice", "password")
Expect(err).To(BeNil())

// Check that an error didn't occur
err = alice.StoreFile("alice.txt", []byte("hello world"))
Expect(err).To(BeNil())

// Check that an error didn't occur AND that the data is what we expect
data, err := alice.LoadFile("alice.txt")
Expect(err).To(BeNil())
Expect(data).To(Equal([]byte("hello world")))

// Check that an error DID occur
data, err := alice.LoadFile("rubbish.txt")
Expect(err).ToNot(BeNil())

Organizing Tests
You can organize tests using some combination of Describe(...) containers, with tests contained within Specify(...) blocks. The more organization you have, the better! Read more about how to organize your tests.

Running a Single Test Case
If you would like to run only one single test (and not both the client_test and client_unittest test suites), you can change the Specify of that test case to a FSpecify. For instance, if you have a test

Specify("Basic Test: Load and Store", func() {...})

you can rename it to

FSpecify("Basic Test: Load and Store", func() {...})

Then, if you click the Run/Debug button, you will be focusing on that specific test. If you have multiple FSpecify test cases, then all those will be run.

Warning: Please remove all FSpecify statements when you submit to the autograder (otherwise, the autograder will only run tests labelled FSpecify) or potentially break the autograder. We cannot promise we will re-run the autograder for you if you forgot to remove them.

Optional: Measure Local Test Coverage
To measure and visualize local test coverage (e.g. how many lines of your implementation your test file hits), you can run these commands in the root folder of your repository:

go test -v -coverpkg ./... ./... -coverprofile cover.out
go tool cover -html=cover.out

Coverage over your own implementation may serve as an indicator for how well your code will perform (with regards to coverage flags) when compared to the staff implementation! It should also help you write better unit testing to catch edge cases.


Design Overview
Library Functions
Users And User Authentication
File Operations
Sharing and Revocation
Debugging and Errors
Advice and Tips
Appendix

Project 3
This site uses Just the Docs, a documentation theme for Jekyll.
Search CS161 Summer 2025
Dark Mode
Textbook
Ed
OH Queue
Extensions
Anonymous Feedback
Project 2	Advice and Tips
Tips on Database Architecture
Minimizing Complexity
Authenticated Encryption
Checking for Errors
Tips For Writing Test Cases
Coverage Flags Tips
Writing Advanced Tests
Writing Efficiency Tests
Expect(err).ToNot(BeNil())
Notes on Key Reuse
Notes on Key Management
More Notes on Key Reuse
Strings and Byte Arrays
Advice and Tips
Tips on Database Architecture
Whenever you’re thinking about storing lists or maps of things in Datastore, it may help to think about how you can “flatten” your data structure - this might help with security, efficiency, and/or code complexity! As an example: let’s say we have a list of restaurants, where each restaurant has a list of menu items. Let’s say I want to figure out whether a particular restaurant has a specific menu item, and we’re using a datastore-like backend (a key-value store).

Consider these two approaches to representing this information in some key-value store:

// Approach A
{
   "movies": {
      "mcu": ["iron man", "the incredible hulk", "iron man 2", "thor", "captain america", "the avengers"],
      "dceu": ["batman"],
      "pixar": ["turning red"]
   }
}

// Approach B
{
   "movies/mcu": ["iron man", "the incredible hulk", "iron man 2", "thor", "captain america", "the avengers"],
   "movies/dceu": ["batman"],
   "movies/pixar": ["turning red"]
} 

Both of these represent the same data, except instead of storing all of this information in one place, we’re flattening this data out across multiple entries in our key-value store. In approach B, all I need to do is to index into the datastore at the string value I construct (e.g. in this case, "movies/pixar"), and I can retrieve all of the movies that are affiliated with Pixar without having to retrieve a bunch of other data that was previously stored at the same location!

If efficiency is measured in read/write bandwidth (like it is in this project), then it’s much more efficient to represent our data in this flattened structure.

Minimizing Complexity
If you don’t need to store a certain piece of data, then don’t store it. Just recompute it (or re-fetch it, if it’s coming from datastore) “on-the-fly” when you need it. It’ll make your life easier!

Authenticated Encryption
In order to build in support for authenticated encryption, you may find it helpful to create a struct with two properties: (a) some arbitrary ciphertext, and (b) some tag (e.g. MAC, signature) over the ciphertext. Then, you could marshal that struct into a sequence of bytes, and store that in datastore. When loading it from datastore, you could (a) unmarshal it and check the tag for authenticity/integrity, and then decrypt the ciphertext and pass the plaintext downstream.

Checking for Errors
Throughout each API method, you’ll probably have several dozen error checks. Take a look at the following code for an example of good and bad practice when it comes to handling errors.

// This is bad practice: the error is discarded!
value, _ = SomeAPICall(...)

// This is good practice: the error is checked!
value, err = SomeAPICall(...)
if (err != nil) {
    return err;
}

When an error is detected (e.g. malicious tampering occurred), your client just needs to immediately return. We don’t care about what happens afterwards (e.g. program state could be messed up, could have file inconsistency, etc. - none of this matters). In other words, we don’t care about recovery – we solely care about detection.

You should almost never discard the output of an error: always, always check to see if an error occurred!

Tips For Writing Test Cases
Here are a few different ways to think about creative tests:

Functionality Tests
Consider basic functionality for single + multiple users
Consider different sequences of events of Client API methods
Consider interleaving file sharing/revocation/loading/storing actions
Edge Case Tests
Consider edge cases outlined in the Design Overview and other parts of this specification.
Security Tests
Consider what happens when an attacker tampers with different pieces of data
In all cases, you should ensure that your Client API is following all requirements listed in the Design Overview.

Coverage Flags Tips
It’s okay if your tests don’t get all 20 test coverage points! Some flags are very tricky to get, and we don’t expect everyone to get all the flags. If you’re missing a few flags, a few points won’t make or break your project. That said, higher test coverage does mean you’re testing more edge cases, which means you can also be more confident in your own implementation if it’s passing all your tests.

Writing Advanced Tests
Some students have reported success with fuzz testing, which uses lots of different random tampering attacks. Sometimes this can help you catch an edge case you weren’t expecting.
Remember that your tests must work on any project implementation, not just your own implementation. This means you cannot assume anything about the design except the API calls that everyone implements (InitUser, GetUser, LoadFile, StoreFile, etc). For example, you cannot reference a field in a user struct in your tests, because other implementations might not have that field in their user struct.
The userlib has a nifty DatastoreGetMap function which returns the underlying go map structure of the Datastore. This can be used to modify the Datastore directly to simulate attacker action.
DatastoreGetMap can also be used to learn about how the Datastore is changed as a result of an API call. For example, you can scan the Datastore, perform an API call (e.g. StoreFile), then scan the Datastore again to see what has been updated. This can help you write more sophisticated tests that leverage information about what encrypted Datastore entries correspond to what data.
Writing Efficiency Tests
Here’s a helper function you can use to measure efficiency.

// Helper function to measure bandwidth of a particular operation
measureBandwidth := func(probe func()) (bandwidth int) {
   before := userlib.DatastoreGetBandwidth()
   probe()
   after := userlib.DatastoreGetBandwidth()
   return after - before
}

// Example usage
bw = measureBandwidth(func() {
   alice.StoreFile(...)
})

Expect(err).ToNot(BeNil())
You should check for errors after every Client API call. This makes sure errors are caught as soon as they happen, instead of them propagating downstream. Don’t make any assumptions over which methods will throw and which errors won’t: use an Expect after each API method.

Remember the core principle around testing: the more lines of the staff solution you hit, the more flags you’re likely to acquire! Think about all of the if err != nil cases that may be present in the staff code, and see if you’re able to write tests to enter those cases.

value, err = user.InitUser(...);
Expect(err).ToNot(BeNil());

Notes on Key Reuse
In general, avoid reusing a key for multiple purposes. Why? Here’s a snippet from a post that Professor David Wagner wrote:

Suppose I use a key k to encrypt and MAC a username (say), and elsewhere in the code I use the same key k to encrypt and MAC a filename. This might open up some attacks. For instance, an attacker might be able to take the encrypted username and copy-paste it over the encrypted filename. When the client decrypts the filename, everything decrypts fine, and the MAC is valid (so no error is detected), but now the filename has been replaced with a username.

I’ve seen many protocols that have had subtle flaws like this.

There are two ways to protect yourself from this. One way is to think really, really hard about these kind of copy-paste attacks whenever you reuse a key for multiple purposes. The other way is to never reuse a key for multiple purposes. I don’t like to think hard. If I have to think hard, I often make mistakes. Or, as Brian Kernighan once wrote: “Everyone knows that debugging is twice as hard as writing a program in the first place. So if you’re as clever as you can be when you write it, how will you ever debug it?” So, I recommend way #2: never reuse a key for multiple purposes.

Hopefully this rationale will help you recognize the motivation behind that advice, which might help you recognize when the advice does or doesn’t apply.

Notes on Key Management
In order to simplify your key management scheme, it may be useful to store a small number of root keys, and re-derive many keys for specific purposes “on-the-fly” (a.k.a. when you need them) using the HashKDF function.

More Notes on Key Reuse
The term “key reuse” can be ambiguous.

As defined in our lecture, “key reuse” refers to the practice of using the same key in multiple algorithms. For example, if you use the same key to encrypt some data and HMAC that same data, this is a case of “key reuse” as defined in lecture. As discussed in lecture, you should always use one key per algorithm.

Note that using the same key to encrypt different pieces of data, or using the same key to MAC/sign different pieces of data, is not considered key reuse (as defined in lecture). Recall that the whole point of discarding one-time pads and introducing block ciphers was to allow for using the same key to encrypt different pieces of data.

Even though using one key per algorithm solves the “key reuse” problem from lecture, it doesn’t necessarily mean that other issues with key reuse don’t exist.

As an example, suppose you use the same, hard-coded key to encrypt every single value in your design. You use a different key (also hard-coded) to MAC every single value in your design. This is not a case of “key reuse” from lecture, because you did use one key per algorithm. However, this is still an insecure design because an attacker could read the hard-coded keys in your code and decrypt and modify all your stored data.

In summary: Reusing the same key in different algorithms is the “key reuse” problem from lecture, and you should always avoid this. Other cases of key reuse may exist in your design; it’s your job to figure out when these cases are problematic.

Strings and Byte Arrays
Be very careful about casting byte arrays to strings. Byte arrays can contain any byte between 0 and 255, including some bytes that correspond to unprintable characters. If you cast a byte array containing unprintable characters to a string, the resulting string could have unexpected behavior. A safer way to convert byte arrays to strings (and back) is to use the encoding/hex library, which we’ve imported for you.

If you instead need to convert strings to byte arrays (and back), you can use json.Marshal. Recall that marshalling works on any data structure, including strings.


Appendix




A friendly request: please do not make your solution public!

Write your implementation in `client/client.go` and your integration tests in `client_test/client_test.go`. Optionally, you can also use `client/client_unittest.go` to write unit tests (e.g: to test your helper functions).

To test your implementation, run `go test -v` inside of the `client_test` directory. This will run all tests in both `client/client_unittest.go` and `client_test/client_test.go`.
<img width="881" height="394" alt="Screenshot 2025-08-12 at 2 41 23 AM" src="https://github.com/user-attachments/assets/1f7be1f7-49b6-42c6-bf04-feb2a790dcb1" />

My Design:
We will use three core structs to manage user authentication, file storage, and sharing
in the Datastore:
1. User Struct
• Fields: Username (user ID), PasswordHash (hashed password), PrivateKey (PKE
decryption), SigningKey (digital signatures).
• Storage: UUID from hash(username), encrypted with password-derived key
(HashKDF), MAC for integrity.
2. File Struct
• Fields: Owner (username), Content (file data), SharedWith (map of usernames to
invitation UUIDs).
• Storage: UUID from hash(filename + username), encrypted with file-specific
symmetric key, MAC for integrity.
3. Invitation Struct
• Fields: FileUUID (file location), Sender/Recipient (usernames), EncryptionKey
(file key, encrypted with recipient’s public key).
• Storage: Random UUID (uuid.New()), encrypted with recipient’s public key, MAC
and signature for integrity.
• UUID Generation: User struct uses hash(username), File struct uses
hash(filename + username), Invitation struct uses random UUID.
• Confidentiality: Encrypt User with password-derived key, File with symmetric key,
Invitation with recipient’s public key.
• Integrity: MACs for all structs, digital signature for Invitation.
• Keys: User keys from password (HashKDF), file key (random, stored in Invitation),
public keys in Keystore.
Q2: Design Question: Datastore Adversary
Storing a hashed password in the Datastore is insecure because the adversary can
read or tamper with it, enabling oSline attacks or data corruption. It also lacks a
mechanism to securely store or retrieve cryptographic keys needed for file access
and sharing, and it doesn’t ensure data integrity.
Storage Strategy:
• User Data: Stored in User struct at hash(username) UUID, encrypted with
password-derived key, MAC for integrity.
• File Data: Stored in File struct at hash(filename + username) UUID, encrypted
with symmetric key, MAC for integrity.
• Invitation Data: Stored in Invitation struct at random UUID, encrypted with
recipient’s public key, MAC and signature for integrity.
Key Management:
• User struct: Encryption key derived from password via userlib.HashKDF,
recomputed in GetUser, not stored.
• File struct: Random symmetric key (userlib.RandomBytes), stored encrypted in
Invitation or User struct for owner access.
• Invitation struct: File key encrypted with recipient’s public key from Keystore,
ensuring only the recipient can access it.
• Keys are never stored unencrypted, and the Keystore provides public keys for
encryption and verification.
Detecting Tampering: MACs (via userlib.HMACEval) are stored with each struct’s
encrypted data. On retrieval, userlib.HMACEqual verifies the MAC; a mismatch
signals tampering. Invitations use digital signatures (via userlib.DSSign) to ensure
sender authenticity, checked with userlib.DSVerify. This ensures any adversary
modifications are detected, protecting all operations like LoadFile and
AcceptInvitation.
Diagram
Diagram Explanation:
• Datastore (blue): Contains User (cyan, at hash(username)), File (green, at
hash(filename + username)), and Invitation (orange, at random UUID) structs,
each encrypted and MAC-protected.
• Keystore (purple): Stores public keys (PKEEncKey, DSVerifyKey) by username.
• Arrows: Show Shared with linking files to invitations, decryption with private
keys, and encryption with public keys.
• Design: Wide boxes (4.5cm), 1.5cm spacing, and raised labels (0.7cm) ensure
no text overlap, maintaining cla
