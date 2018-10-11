// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package devkit

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"openpitrix.io/openpitrix/pkg/devkit/opapp"
)

const DefaultServeAddr = "127.0.0.1:9191"

const indexHTMLTemplate = `
<html>
<head>
	<title>OpenPitrix Repository</title>
</head>
<h1>OpenPitrix Apps Repository</h1>
<h6>Thanks to <a href="https://helm.sh/">Helm Project</a></h6>
<ul>
{{range $name, $ver := .Index.Entries}}
  <li>{{$name}}<ul>{{range $ver}}
    <li><a href="{{index .URLs 0}}">{{.Name}}-{{.OpVersion}}</a></li>
  {{end}}</ul>
  </li>
{{end}}
</ul>
<body>
<p>Last Generated: {{.Index.Generated}}</p>
</body>
</html>
`

// RepositoryServer is an HTTP handler for serving a app repository.
type RepositoryServer struct {
	RepoPath string
}

// ServeHTTP implements the http.Handler interface.
func (s *RepositoryServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	switch uri {
	case "/", "/apps/", "/apps/index.html", "/apps/index":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		s.htmlIndex(w, r)
	default:
		file := strings.TrimPrefix(uri, "/apps/")
		http.ServeFile(w, r, filepath.Join(s.RepoPath, file))
	}
}

// StartLocalRepo starts a web server and serves files from the given path
func StartLocalRepo(path, address string) error {
	if address == "" {
		address = DefaultServeAddr
	}
	s := &RepositoryServer{RepoPath: path}
	return http.ListenAndServe(address, s)
}

func (s *RepositoryServer) htmlIndex(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("index.html").Parse(indexHTMLTemplate))
	// load index
	lrp := filepath.Join(s.RepoPath, "index.yaml")
	i, err := opapp.LoadIndexFile(lrp)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data := map[string]interface{}{
		"Index": i,
	}
	if err := t.Execute(w, data); err != nil {
		fmt.Fprintf(w, "Template error: %s", err)
	}
}
