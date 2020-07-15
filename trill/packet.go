package trill

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/Ox1O5/trill/utils"
)

type IPacket interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}

const headLen = 8

type packet struct{}

func NewPacket() *packet {
	return &packet{}
}

func (p packet) GetHeadLen() uint32 {
	return headLen
}

func (p *packet) Pack(msg IMessage) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *packet) UnPack(binaryData []byte) (IMessage, error) {
	buf := bytes.NewReader(binaryData)
	msg := &message{}
	if err := binary.Read(buf, binary.LittleEndian, &msg.dataLen); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &msg.id); err != nil {
		return nil, err
	}
	if utils.GlobalObject.MaxPacketSize < msg.dataLen {
		return nil, errors.New("Too large msg data received\n")
	}
	return msg, nil
}
