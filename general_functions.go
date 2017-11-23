/****** MIT License **********
Copyright (c) 2017 Datzer Rainer
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
***************************/

package main

import (
	"fmt"
	"runtime"
	//"io"
	//"io/ioutil"
	//"reflect"
	"github.com/davecgh/go-spew/spew"
)

/// @see https://stackoverflow.com/questions/17640360/file-or-line-similar-in-golang
func file_line() string {
	_, fileName, fileLine, ok := runtime.Caller(1)
	var s string
	if ok {
		s = fmt.Sprintf("%s:%d", fileName, fileLine)
	} else {
		s = ""
	}
	return s
}

func byte_array_to_string(byteArray []byte) string {
	s := string(byteArray[:])
	return s
}

func string_to_byte_array(str string) []byte {
	byteArray := []byte(str)
	return byteArray
}

func dump_r(data interface{}) string {
	var s string
	s = fmt.Sprintf("%v", data)
	return s
}

func dump_r_v(data interface{}) string {
	var s string
	//s = fmt.Sprintf("%v", reflect.ValueOf(data))
	s = spew.Sdump(data)
	return s
}
