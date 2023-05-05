package index

import (
	"net/http"
)

type IndexHandler struct {
	Name    string
	Version string
	Path    string
}

func (i IndexHandler) Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<html>
		<head><title>` + i.Name + `</title></head>
		<body>
		<h1>` + i.Name + ` @ ` + i.Version + `</h1>
		<p><a href="` + i.Path + `">Metrics</a></p>
		<p><a href="/prober?target=">Prober (multi-target export pattern)</a></p>
		<p><a href="https://github.com/Luzilla/dnsbl_exporter">Code on Github</a></p>
		</body>
		</html>`))
}
