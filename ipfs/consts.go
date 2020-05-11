package ipfs

// IpfsFileLength defines the length of the ipfs filename
var IpfsFileLength int

// SetConsts sets values for consts related to ipfs
func SetConsts(fileLength int) {
	IpfsFileLength = fileLength
}
