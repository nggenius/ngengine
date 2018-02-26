package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
)

type StoreArchive struct {
	buf []byte
	pos int
}

func NewStoreArchiver(buf []byte) *StoreArchive {
	if buf == nil || cap(buf) == 0 {
		return nil
	}
	ar := &StoreArchive{}
	ar.buf = buf[:0]
	ar.pos = 0
	return ar
}

func (ar *StoreArchive) Write(p []byte) (n int, err error) {
	l := len(p)
	if ar.pos+l > cap(ar.buf) {
		return 0, io.EOF
	}
	ar.buf = ar.buf[:ar.pos+l]
	copy(ar.buf[ar.pos:], p)
	ar.pos += l
	return l, nil
}

func (ar *StoreArchive) Data() []byte {
	return ar.buf[:ar.pos]
}

func (ar *StoreArchive) Len() int {
	return ar.pos
}

func (ar *StoreArchive) WriteAt(offset int, val interface{}) error {
	if offset >= cap(ar.buf) {
		return fmt.Errorf("offset out of range")
	}

	old := ar.pos
	ar.pos = offset
	var err error
	switch val.(type) {
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64:
		err = binary.Write(ar, binary.LittleEndian, val)
	case int:
		err = binary.Write(ar, binary.LittleEndian, int32(val.(int)))
	default:
		err = fmt.Errorf("unsupport type")
	}

	ar.pos = old
	return err
}

func (ar *StoreArchive) Put(val interface{}) error {
	switch val.(type) {
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64:
		return binary.Write(ar, binary.LittleEndian, val)
	case int:
		return binary.Write(ar, binary.LittleEndian, int64(val.(int)))
	case string:
		return ar.PutString(val.(string))
	case []byte:
		return ar.PutData(val.([]byte))
	default:
		return ar.PutObject(val)
	}
}

func (ar *StoreArchive) PutString(val string) error {
	data := []byte(val)
	size := len(data)
	err := binary.Write(ar, binary.LittleEndian, int16(size))
	if err != nil {
		return err
	}
	_, err = ar.Write(data)
	return err
}

func (ar *StoreArchive) PutObject(obj interface{}) error {
	enc := gob.NewEncoder(ar)
	return enc.Encode(obj)
}

func (ar *StoreArchive) PutData(data []byte) error {
	err := ar.Put(uint16(len(data)))
	if err != nil {
		return err
	}
	_, err = ar.Write(data)
	return err
}

type LoadArchive struct {
	reader *bytes.Reader
}

func NewLoadArchiver(data []byte) *LoadArchive {
	ar := &LoadArchive{}
	ar.reader = bytes.NewReader(data)
	return ar
}

func (ar *LoadArchive) Position() int {
	return int(ar.reader.Size()) - ar.reader.Len()
}

func (ar *LoadArchive) AvailableBytes() int {
	return ar.reader.Len()
}

func (ar *LoadArchive) Size() int {
	return int(ar.reader.Size())
}

func (ar *LoadArchive) Seek(offset int, whence int) (int, error) {
	ret, err := ar.reader.Seek(int64(offset), whence)
	return int(ret), err
}

func (ar *LoadArchive) Read(val interface{}) (err error) {
	switch val.(type) {
	case *int8, *int16, *int32, *int64, *uint8, *uint16, *uint32, *uint64, *float32, *float64:
		return binary.Read(ar.reader, binary.LittleEndian, val)
	case *int:
		var out int64
		err = binary.Read(ar.reader, binary.LittleEndian, &out)
		if err != nil {
			return err
		}
		*(val.(*int)) = int(out)
		return nil
	case *string:
		inst := val.(*string)
		*inst, err = ar.ReadString()
		return err
	case *[]byte:
		inst := val.(*[]byte)
		*inst, err = ar.ReadData()
		return err
	default:
		return ar.ReadObject(val)
	}
}

func (ar *LoadArchive) ReadInt8() (val int8, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadUInt8() (val uint8, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadInt16() (val int16, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadUInt16() (val uint16, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadInt32() (val int32, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadUInt32() (val uint32, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadInt64() (val int64, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadUInt64() (val uint64, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadFloat32() (val float32, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadFloat64() (val float64, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadString() (val string, err error) {
	var size int16
	binary.Read(ar.reader, binary.LittleEndian, &size)
	if size == 0 {
		val = ""
		return
	}
	data := make([]byte, size)
	_, err = ar.reader.Read(data)
	if err != nil {
		return
	}
	val = string(data)
	return
}

func (ar *LoadArchive) ReadObject(val interface{}) error {
	dec := gob.NewDecoder(ar.reader)
	return dec.Decode(val)
}

func (ar *LoadArchive) ReadData() (data []byte, err error) {
	var l uint16
	l, err = ar.ReadUInt16()
	data = make([]byte, int(l))
	_, err = ar.reader.Read(data)
	return data, err
}
