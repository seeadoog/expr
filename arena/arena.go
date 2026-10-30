package arena

import "unsafe"

type Arena struct {
	buf    []byte
	offset int
}

func NewArena(size int) *Arena {
	return &Arena{
		//buf: make([]byte, size),
	}
}
func (a *Arena) alloc(size int) unsafe.Pointer {
	data := a.allocBytes(size)
	return unsafe.Pointer(&data[0])
}

func (a *Arena) allocBytes(size int) []byte {
	if a.offset+size > len(a.buf) {
		a.buf = make([]byte, size+len(a.buf)*2)
		a.offset = 0
	}
	data := a.buf[a.offset : a.offset+size]
	a.offset += size
	return data
}
func (a *Arena) free() {
	a.offset = 0
}

func (a *Arena) Reset() {
	a.offset = 0
}

var (
	stringType = typt("")
	floatType1 = typt(1.0)
)

func (a *Arena) StringToAny(s string) any {
	if a == nil {
		return s
	}
	ptr := New[string](a)
	*ptr = s
	si := iface{
		t: stringType.t,
		d: uintptr(unsafe.Pointer(ptr)),
	}
	return *(*any)(unsafe.Pointer(&si))
}
func (a *Arena) Float64ToAny(s float64) any {
	if a == nil {
		return s
	}
	ptr := New[float64](a)
	*ptr = s
	si := iface{
		t: floatType1.t,
		d: uintptr(unsafe.Pointer(ptr)),
	}
	return *(*any)(unsafe.Pointer(&si))
}

func (a *Arena) MakeBytes(size int) []byte {
	data := a.allocBytes(size)

	sh := &sliceHeader{
		data: unsafe.Pointer(&data[0]),
		len:  uintptr(len(data)),
		cap:  uintptr(len(data)),
	}
	return *(*[]byte)(unsafe.Pointer(sh))
}

func (a *Arena) Append(b []byte, b2 ...byte) []byte {
	if a.offset >= len(a.buf) {
		return append(b, b2...)
	}
	if len(b2) == 0 {
		return b
	}
	curptr := uintptr(unsafe.Pointer(&a.buf[a.offset]))
	byteLastPtr := uintptr(unsafe.Pointer(&b[0])) + uintptr(len(b))
	if byteLastPtr == curptr { //
		if a.offset+len(b2) <= len(a.buf) {
			sh := &sliceHeader{
				data: unsafe.Pointer(&b[0]),
				len:  uintptr(len(b) + len(b2)),
				cap:  uintptr(len(b) + len(b2)),
			}
			ret := *(*[]byte)(unsafe.Pointer(sh))
			copy(ret[len(b):], b2)
			return ret
		}
	}
	data := a.MakeBytes(len(b) + len(b2))
	copy(data, b)
	copy(data[len(b):], b2)
	return data
}

func New[T any](a *Arena) *T {
	var t T
	ptr := a.alloc(int(unsafe.Sizeof(t)))
	return (*T)(ptr)
}

type iface struct {
	t, d uintptr
}

func typf(i any) uintptr {
	p := (*iface)(unsafe.Pointer(&i))
	return p.d
}

func typt(i interface{}) iface {
	p := (*iface)(unsafe.Pointer(&i))
	return *p
}

type sliceHeader struct {
	data unsafe.Pointer
	len  uintptr
	cap  uintptr
}
