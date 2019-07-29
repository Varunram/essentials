package ipfs

// IpfsFileLength defines the length of the ipfs filename
var IpfsFileLength int

func SetConsts(fileLength int) {
	IpfsFileLength = fileLength
}
