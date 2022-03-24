package uuid

import "github.com/google/uuid"

func CreateUUID() string {
	return uuid.New().String()
}

func CreatUUIDBinary() []byte {
	d,_ := uuid.New().MarshalBinary()
	return d
}