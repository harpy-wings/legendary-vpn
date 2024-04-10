package cipher

type Method uint8

const (
	MethodUnknown = iota
	MethodAES256
	// todo add more methods
)
const (
	BlockSizeAES256 = 16
)

var (
	blockSizes = map[Method]int{
		MethodAES256: BlockSizeAES256,
	}
)
