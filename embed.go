package mows

import (
	"io/fs"
	"net/http"
	"path"
)

// StaticFS serves static files from an embedded filesystem (embed.FS).
//
// prefix: URL path to serve under (e.g., "/static")  
// filesystem: embedded FS containing static files  
// root: subdirectory in the FS to serve
//
// Example:
//
//    app.StaticFS("/static", publicFS, "public")
//
// Then accessing /static/css/app.css will serve embedded public/css/app.css
func (e *Engine) StaticFS(prefix string, filesystem fs.FS, root string) {
	subFS, err := fs.Sub(filesystem, root)
	if err != nil {
		panic(err)
	}

	fsHandler := http.FileServer(http.FS(subFS))
	routePath := path.Join(prefix, "/*filepath")

	e.GET(routePath, func(c *Context) error {
		http.StripPrefix(prefix, fsHandler).ServeHTTP(c.Writer, c.Request)
		return nil
	})
}
