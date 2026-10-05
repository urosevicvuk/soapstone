// Paket web sadrži HTML šablon stranice. go:embed ga ugrađuje u binarni
// fajl pri build-u, pa aplikacija ne čita ništa sa diska.
package web

import "embed"

//go:embed index.html
var FS embed.FS
