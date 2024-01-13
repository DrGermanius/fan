package fan

import (
	"unsafe"
)

type hchan struct {
	qcount   uint
	dataqsiz uint
	buf      unsafe.Pointer
	elemsize uint16
	elemtype uint
	sendx    uint   
}

//go:noinline
func In[T any](ch ...*chan T) chan T {
	var buf []T
	var qcount uint
	fan := make(chan T, 1)
	for i := range ch {
		p := (*hchan)(unsafe.Pointer((*(*uintptr)(unsafe.Pointer(ch[i])))))
		qcount += p.qcount

		*ch[i] = fan

		s := unsafe.Slice((*T)(p.buf), p.qcount)
		buf = append(buf, s...)
	}
	p := (*hchan)(unsafe.Pointer((*(*uintptr)(unsafe.Pointer(&fan)))))
	p.qcount = qcount 
	p.dataqsiz = qcount + 1 
	p.buf = unsafe.Pointer(unsafe.SliceData(buf))
	p.sendx = qcount
	return fan
}
