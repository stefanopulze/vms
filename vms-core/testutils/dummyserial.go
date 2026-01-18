package testutils

import (
	"fmt"
)

type Command struct {
	Request  []byte
	Response []byte
}

type DummySerial struct {
	commands    map[string][]byte
	lastCommand Command
}

func (ds *DummySerial) Start() {
}

func (ds *DummySerial) Close() error {
	return nil
}

func NewDummySerial() *DummySerial {
	return &DummySerial{
		commands: make(map[string][]byte),
	}
}

func (ds *DummySerial) Write(data []byte) ([]byte, error) {
	reply, ok := ds.commands[string(data)]
	if !ok {
		return nil, fmt.Errorf("command not found: %s", data)
	}

	ds.lastCommand = Command{Request: data, Response: reply}

	return reply, nil
}

func (ds *DummySerial) MockCommand(cmd string, req []byte, resp []byte) *DummySerial {
	ds.commands[string(req)] = resp
	return ds
}

func (ds *DummySerial) LastCommand() Command {
	return ds.lastCommand
}

func MockStandardCommands(ds *DummySerial) {
	ds.MockCommand(
		"QPIGS",
		FromHex("5150494753B7A90D"),
		FromHex("283234342E322035302E30203233302E302035302E302030333638203033343020303034203339352035302E3130203031352030353620303032362030332E38203332302E332030302E30302030303030302030303031303131302030302030302030313233302030313061AC0D"),
	)

	ds.MockCommand(
		"QPIRI",
		FromHex("5150495249F8540D"),
		FromHex("283233302E302033342E37203233302E302035302E302033342E37203830303020383030302034382E302034382E302034352E302035332E322035332E32203320303032203132302031203220332039203031203020302034382E3520302031203438302030203030309C920D"),
	)

	ds.MockCommand(
		"QFLAG",
		FromHex("51464C414798740D"),
		FromHex("2845616B78797A44626A75763B790D"),
	)

	ds.MockCommand(
		"QMOD",
		FromHex("514D4F4449C10D"),
		FromHex("2850D5BA0D"),
	)

	ds.MockCommand(
		"QT",
		FromHex("515427FF0D"),
		FromHex("2832303235313232373135303832302CF00D"),
	)
}
