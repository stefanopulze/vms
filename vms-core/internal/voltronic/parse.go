package voltronic

import (
	"strconv"
)

func parseBoolean(v []byte) bool {
	return v[0] == '1'
}

func parseFirmwareVersion(data []byte, prefix string) (*Firmware, error) {
	// TODO add length check

	if nak == string(data) {
		return &Firmware{
			Major: -1,
			Minor: -1,
		}, nil
	}

	// (VERFW(X):00046.82
	data = data[len(prefix)+2:]
	majorString := string(data[:len(data)-3])
	minorString := string(data[len(data)-2:])
	major, err := strconv.Atoi(majorString)
	if err != nil {
		return nil, err
	}

	minor, err := strconv.Atoi(minorString)
	if err != nil {
		return nil, err
	}

	return &Firmware{
		Major: major,
		Minor: minor,
	}, err
}
