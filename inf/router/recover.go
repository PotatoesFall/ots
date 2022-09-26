package router

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			v := recover()
			if v == nil {
				return
			}

			fmt.Fprintf(os.Stderr, "panic: %v\n", v)
			printPanicStack()

			w.WriteHeader(http.StatusInternalServerError)
		}()

		buf := &responseBuffer{ResponseWriter: w}
		next.ServeHTTP(buf, r)

		if buf.code != 0 {
			w.WriteHeader(buf.code)
		}
		if buf.body.Len() != 0 {
			if _, err := io.Copy(w, &buf.body); err != nil {
				panic(err)
			}
		}
	})
}

func printPanicStack() {
	_, filename, _, _ := runtime.Caller(1)

	pc := make([]uintptr, 100)
	entries := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:entries])

	var buf bytes.Buffer
	more, firstLine := true, false
	for more {
		var frame runtime.Frame
		frame, more = frames.Next()

		if !firstLine && frame.File == filename { // Skip frames from this file
			continue
		}
		firstLine = true

		if frame.Function == `runtime.gopanic` { // If a panic occurred, start at the frame that called panic
			buf.Reset()
			continue
		}

		fmt.Fprintf(&buf, "\t%s %s:%d\n", frame.Function, frame.File, frame.Line)
	}

	_, err := io.Copy(os.Stderr, &buf)
	if err != nil {
		panic(err)
	}
}

type responseBuffer struct {
	http.ResponseWriter

	code int
	body bytes.Buffer
}

func (buf *responseBuffer) WriteHeader(code int) {
	buf.code = code
}

func (buf *responseBuffer) Write(data []byte) (int, error) {
	return buf.body.Write(data)
}
