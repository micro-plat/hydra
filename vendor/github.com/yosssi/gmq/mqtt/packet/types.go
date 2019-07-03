package packet

// MQTT Control Packet types
const (
	TypeCONNECT     byte = 0x01
	TypeCONNACK     byte = 0x02
	TypePUBLISH     byte = 0x03
	TypePUBACK      byte = 0x04
	TypePUBREC      byte = 0x05
	TypePUBREL      byte = 0x06
	TypePUBCOMP     byte = 0x07
	TypeSUBSCRIBE   byte = 0x08
	TypeSUBACK      byte = 0x09
	TypeUNSUBSCRIBE byte = 0x0A
	TypeUNSUBACK    byte = 0x0B
	TypePINGREQ     byte = 0x0C
	TypePINGRESP    byte = 0x0D
	TypeDISCONNECT  byte = 0x0E
)
