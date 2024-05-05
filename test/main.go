package test

//
//import (
//	"fmt"
//	"unsafe"
//)
//
//func main() {
//	a, b := 1, 1
//	NilIsNil(a)
//
//	c := 1
//	fmt.Println("c:", c)
//	fmt.Println(a, b, c)
//
//}
//
//func NilIsNil(a any) {
//	// 打开行代码输出 10 10 10 [0 0 0 0 0 0 0 0 0 0] 10
//	// 注释输出 1 1 1 [0 0 0 0 0 0 0 0 0 0] 1
//	//fmt.Println(&a)
//	fmt.Println(a)
//	//a = 10
//	type rtype struct {
//		t    unsafe.Pointer
//		data unsafe.Pointer
//	}
//	t := (*rtype)(unsafe.Pointer(&a))
//	*(*int)(t.data) = 10
//	fmt.Println("data pointer:", t.data)
//}
