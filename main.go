package main

import (
	"reflect"
	"syscall/js"
	"unsafe"

	"github.com/Query-Interface/gofpdf"
)

func main() {

	// define the function in the Javascript scope
	js.Global().Set("CreateBasicPdf", js.FuncOf(CreateBasicPdf))

	// prevent the function from returning, which is required in a WASM module
	select {}
}

func CreateBasicPdf(this js.Value, args []js.Value) interface{} {
	// import the JS console object
	console := js.Global().Get("console")

	// create a basic PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "hello world ! (generated from go-wasm binary) ...")

	// gets the output into a PDF
	buffer := pdf.OutputToBuffer()

	// TODO: try to replace this with linear memory pointer
	dst := js.Global().Get("Uint8Array").New(len(buffer.Bytes()))
	js.CopyBytesToJS(dst, buffer.Bytes())
	console.Call("log", "bytes copied:", Bytes2StrImp(buffer.Bytes()))

	// callback JS to create the download link
	js.Global().Call("downloadPdf", dst)
	return ""
}

// from https://play.golang.org/p/G6rRGp8jwtW
func Bytes2StrImp(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	var s string
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh.Data = sliceHeader.Data
	sh.Len = sliceHeader.Len
	return s
}
