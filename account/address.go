package address

import (
    "log"
    "bytes"
    "crypto/rand"
    "encoding/hex"
    "encoding/base32"
    "encoding/binary"

    "github.com/golang/crypto/blake2b"

    // golangs ed25519 implementation 
    // forked to use blake2b instead of sha3
    "github.com/frankh/crypto/ed25519"
)

// xrb uses a non-standard base32 character set.
const EncodeXrb = "13456789abcdefghijkmnopqrstuwxyz"

// Returns random 32 bit byte array
func GenerateSeed() (seed []byte, err error) {

    seed = make([]byte, 32)
    _, err = rand.Read(seed)
    if err != nil {
        return
    }

    return
}

func KeyPairFromSeed(seed []byte, index uint32) (pubKey ed25519.PublicKey, privKey ed25519.PrivateKey, err error) {

    hash, err := blake2b.New(32, nil)
    if err != nil {
        return
    }

    bs := make([]byte, 4)
    binary.BigEndian.PutUint32(bs, index)

    hash.Write(seed)
    hash.Write(bs)

    seedBytes := hash.Sum(nil)
    pubKey, privKey, err = ed25519.GenerateKey(bytes.NewReader(seedBytes))

    return
}


func GetAddressChecksum(pub ed25519.PublicKey) (checksum []byte, err error) {

    hash, err := blake2b.New(5, nil)
    if err != nil {
        return
    }

    hash.Write(pub)

    checkBytes := hash.Sum(nil)

    for i := len(checkBytes) -1; i >= 0; i-- {
        checksum = append(checksum, checkBytes[i])
    }

    return
}


func AddressFromPubKey(pub ed25519.PublicKey) (xrbAddress string, err error) {

    XrbEncoding := base32.NewEncoding(EncodeXrb)

    padded := append([]byte{0, 0, 0}, pub...)
    address := XrbEncoding.EncodeToString(padded)[4:]

    rawChecksum, err := GetAddressChecksum(pub)
    if err != nil {
        return
    }

    checksum := XrbEncoding.EncodeToString(rawChecksum)

    xrbAddress = "xrb_" + address + checksum

    return
}

/*
func main() {

    seed, err := GenerateSeed()
    if err != nil {
        panic("Failed to generate seed")
    }

    pub, priv, err := KeyPairFromSeed(seed, 0)
    if err != nil {
        panic("Failed to generate key pair")
    }

    address, err := AddressFromPubKey(pub)
    if err != nil {
        panic("Failed to generate address")
    }

    log.Println("Seed:", hex.EncodeToString(seed))
    log.Println("Public Key:", hex.EncodeToString(pub))
    log.Println("Private Key:", hex.EncodeToString(priv))
    log.Println("Address:", address)


}
*/
