package idgenerate

import (
	"github.com/marmotedu/component-base/pkg/util/iputil"
	"github.com/marmotedu/component-base/pkg/util/stringutil"
	"github.com/sony/sonyflake"
	"github.com/speps/go-hashids"
)

// Defiens alphabet.
const (
	Alphabet62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	Alphabet36 = "abcdefghijklmnopqrstuvwxyz1234567890"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	st.MachineID = func() (uint16, error) {
		ip := iputil.GetLocalIP()

		return uint16([]byte(ip)[2])<<8 + uint16([]byte(ip)[3]), nil
	}

	sf = sonyflake.NewSonyflake(st)
}

// GetIntID returns uint64 uniq id.
func GetIntID() (uint64, error) {
	id, err := sf.NextID()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetUUID36 returns id format like: 300m50zn91nwz5.
func GetUUID36(prefix string) string {
	id, _ := GetIntID()
	hd := hashids.NewData()
	hd.Alphabet = Alphabet36

	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}

	i, err := h.Encode([]int{int(id)})
	if err != nil {
		panic(err)
	}

	return prefix + stringutil.Reverse(i)
}
