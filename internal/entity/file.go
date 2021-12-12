package entity

type HashFunction string

const (
	Md5Function    HashFunction = "MD5"
	Sha128Function HashFunction = "SHA128"
	Sha256Function HashFunction = "SHA256"
	Sha512Function HashFunction = "SHA512"
)

type File struct {
	ID       int64
	Function HashFunction
	ByteSize int64
	FileName string
}
