package copyutil

import (
	"bytes"
	"encoding/gob"
)

// 只能拷贝基础变量
// 如果结构体里面还有个属性是结构体，这个属性要单独拎出来拷贝
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)

}
