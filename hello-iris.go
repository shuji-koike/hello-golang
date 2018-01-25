package main

import (
	"archive/zip"
	"github.com/kataras/iris"
	"io"
	"io/ioutil"
	"regexp"
)

func main() {
	app := iris.New()

	// Load all templates from the "./views" folder
	// where extension is ".html" and parse them
	// using the standard `html/template` package.

	app.RegisterView(iris.HTML("./views", ".html"))

	// Method:    GET
	// Resource:  http://localhost:8080
	app.Get("/", func(ctx iris.Context) {
		// Bind: {{.message}} with "Hello world!"
		ctx.ViewData("message", "Hello world!")
		// Render template file: ./views/hello.html
		ctx.View("hello.html")
	})

	app.Get("/finder/{path:path regexp(\\.zip)}", func(ctx iris.Context) {
		ctx.Writef(ctx.Params().Get("path"))
	})

	app.Get("/finder/{path:path}", func(ctx iris.Context) {
		path := ctx.Params().Get("path")
		rep := regexp.MustCompile(`(.*)(\.zip)/(.*)`).FindStringSubmatch(path)
		println(path)
		var async, _ = ctx.Params().GetBool("async")

		if rep != nil {
			zipPath := rep[1] + rep[2]
			name := rep[3]

			r, err := zip.OpenReader(zipPath)
			if err != nil {
				println("error zip.OpenReader()")
			}
			defer r.Close()

			var file *zip.File
			for _, f := range r.File {
				if f.Name == name {
					file = f
					break
				}
			}

			ctx.Header("Content-Type", "image/jpeg")

			rr, err := file.Open()
			if err != nil {
				println("error file.Open()")
			}
			defer rr.Close()

			if async {
				ctx.StreamWriter(func(w io.Writer) bool {
					var buf = make([]byte, 1024*8)
					n, _ := rr.Read(buf)
					if n == 0 {
						return false
					} else {
						w.Write(buf)
					}
					return true
				})
			} else {
				bytes, _ := ioutil.ReadAll(rr)
				ctx.Write(bytes)
			}
		}
	})

	// Start the server using a network address.
	app.Run(iris.Addr(":8081"))
}
