package utils

const (
	ConstHeader         = "www.01happy.com"
	ConstHeaderLength   = 15
	ConstSaveDataLength = 4
)

func Pack(message []byte) []byte {
	return append(append(S2B(ConstHeader), IntToBytes(len(message))...), message...)
}

func UnPack(buffer []byte, readerChannel chan []byte) []byte {
	length := len(buffer)

	var i int
	for i = 0; i < length; i++ {
		if length < i+ConstHeaderLength+ConstSaveDataLength {
			break
		}
		if B2S(buffer[i:i+ConstHeaderLength]) == ConstHeader {
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstSaveDataLength])
			if length < i+ConstHeaderLength+ConstSaveDataLength+messageLength {
				break
			}
			data := buffer[i+ConstHeaderLength+ConstSaveDataLength : i+ConstHeaderLength+ConstSaveDataLength+messageLength]
			readerChannel <- data
		}
	}
	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}
