# Project 2 Starter Code

This repository contains the starter code for Project 2!

For comprehensive documentation, see the Project 2 Spec (https://cs161.org/proj2/).

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
