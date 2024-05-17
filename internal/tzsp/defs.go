package tzsp

type Type uint8
type TagType uint8

const (
	TypeReceivedTagList Type = iota
	TypePacketForTransmit
	TypeReserved
	TypeConfiguration
	TypeKeepAlive
	TypePortOpener
)

const (
	TagPadding        TagType = 0x00
	TagEnd            TagType = 0x01
	TagRawRSSI        TagType = 0x0a
	TagSNR            TagType = 0x0b
	TagDataRate       TagType = 0x0c
	TagTimestamp      TagType = 0xd
	TagContentionFree TagType = 0x0f
	TagDecrypted      TagType = 0x10
	TagFCSError       TagType = 0x11
	TagRXChannel      TagType = 0x12
	TagPacketCount    TagType = 0x28
	TagRXFrameLength  TagType = 0x29
	TagWLANRHDRSerial TagType = 0x3c
)
