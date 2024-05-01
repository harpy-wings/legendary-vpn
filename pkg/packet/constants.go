package packet

const (
	HOP_REQ uint8 = 0x20
	HOP_ACK uint8 = 0xAC
	HOP_DAT uint8 = 0xDA

	HOP_FLG_PSH byte = 0x80 // port knocking and heartbeat
	HOP_FLG_HSH byte = 0x40 // handshaking
	HOP_FLG_FIN byte = 0x20 // finish session
	HOP_FLG_MFR byte = 0x08 // more fragments
	HOP_FLG_ACK byte = 0x04 // acknowledge
	HOP_FLG_DAT byte = 0x00 // acknowledge

	HOP_STAT_INIT      int32 = iota // initing
	HOP_STAT_HANDSHAKE              // handeshaking
	HOP_STAT_WORKING                // working
	HOP_STAT_FIN                    // finishing

	HOP_HDR_LEN int = 16

	HOP_PROTO_VERSION byte = 0x01 // protocol version
)

const (
	packetFormatString = "{%v, Payload: %v, Noise: %v}"
)
