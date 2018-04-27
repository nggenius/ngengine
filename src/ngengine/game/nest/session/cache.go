package session

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

type accinfo struct {
	token      string
	expireTime time.Time
}

func (a *accinfo) Expired() bool {
	return time.Now().Sub(a.expireTime) > 0
}

type cache map[string]*accinfo

// token [ 0 ... 7 ] [ 8 ... 9 ] [ 10 ...  13 ] [ 14 ... 21 ]
//       start time  life time	  random         verify code
func CreateToken(account string, minutes uint16, starttime uint64) string {
	code := make([]byte, 22)
	if starttime == 0 {
		starttime = uint64(time.Now().UTC().Unix())
	}
	binary.LittleEndian.PutUint64(code, uint64(starttime))
	binary.LittleEndian.PutUint16(code[8:], uint16(minutes))
	rand.Read(code[10:14])
	magickey := "abc@123"
	data := fmt.Sprintf("%s%d%d%s", account, starttime, minutes, magickey)
	dec := md5.New()
	dec.Write([]byte(data))
	verify := dec.Sum(nil)
	copy(code[14:], verify[4:12])
	return string(hex.EncodeToString(code))
}

func (c cache) Put(acc string) string {
	if info, dup := c[acc]; dup {
		return info.token
	}

	info := &accinfo{}
	info.expireTime = time.Now().Add(time.Second * 30)
	info.token = CreateToken(acc, 1, 0)
	c[acc] = info
	return info.token
}

func (c cache) Pop(acc string) {
	if _, ok := c[acc]; ok {
		delete(c, acc)
	}
}

func (c cache) Check() {
	for k, v := range c {
		if v.Expired() {
			delete(c, k)
		}
	}
}

func (c cache) Valid(acc string, token string) bool {
	if a, ok := c[acc]; ok {
		if a.token == token {
			delete(c, acc)
			return true
		}
	}

	return false
}
