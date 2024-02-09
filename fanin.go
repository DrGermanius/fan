package fan

import (
	"unsafe"
)

type hchan struct {
	qcount   uint
	dataqsiz uint
	buf      unsafe.Pointer
	elemsize uint16
	closed   uint32
	elemtype uint
	sendx    uint
	recvx    uint
	recvq    waitq
	sendq    waitq

	lock mutex
}

type lockRankStruct struct {}

type mutex struct {
	lockRankStruct
	key uintptr
}

//go:linkname lock runtime.lock
func lock(l *mutex)

//go:linkname unlock runtime.unlock
func unlock(l *mutex)

type waitq struct {
	first *sudog
	last  *sudog
}

type sudog struct {
	g uint

	next *sudog
	prev *sudog
}

//go:noinline
func In[T any](ch ...*chan T) chan T {
	if len(ch) < 1 {
		panic("no chances")
	}

	var buf []T
	var qcount uint
	fan := make(chan T, 1)

	var send waitq
	var recv waitq

	for i := range ch {
		p := (*hchan)(unsafe.Pointer((*(*uintptr)(unsafe.Pointer(ch[i])))))

		*ch[i] = fan


		lock(&p.lock)
		qcount += p.qcount

		if i == 0 {
			send = p.sendq
		} else if p.sendq.first != nil  {
			if send.first == nil {
				send = p.sendq
			} else {
				send.last.next = p.sendq.first
				send.last = p.sendq.last
			}
		}
		
		if i == 0 {
			recv = p.recvq
		} else if p.recvq.first != nil  {
			if recv.first == nil {
				recv = p.recvq
			} else {
				recv.last.next = p.recvq.first
				recv.last = p.recvq.last
			}
		}

		s := unsafe.Slice((*T)(p.buf), p.qcount)
		buf = append(buf, s...)
		unlock(&p.lock)
	}
	if qcount == 0  { 
		return fan
	}
	
	p := (*hchan)(unsafe.Pointer((*(*uintptr)(unsafe.Pointer(&fan)))))
	p.qcount = qcount
	p.dataqsiz = qcount
	p.buf = unsafe.Pointer(unsafe.SliceData(buf))
	// p.sendx = qcount
	p.sendq = send
	return fan
}
