package fan

import (
	"fmt"
	"unsafe"
)

type hchan struct {
	qcount   uint
	dataqsiz uint
	buf      unsafe.Pointer
	elemsize uint16
}

//go:noinline
func In[T any](ch ...*chan T) chan T {
	var buf []T
	var qcount uint
	for i := range ch {
		var p = (*hchan)(unsafe.Pointer((*(*uintptr)(unsafe.Pointer(ch[i])))))
		qcount += p.qcount

		s := unsafe.Slice((*T)(p.buf), p.qcount)
		buf = append(buf, s...)
	}
	fmt.Println(buf)
	fan := make(chan T, qcount)
	var p = (*hchan)(unsafe.Pointer((*(*uintptr)(unsafe.Pointer(&fan)))))
	p.qcount = qcount
	p.buf = unsafe.Pointer(unsafe.SliceData(buf))

	return fan
}