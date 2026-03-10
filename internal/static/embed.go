package static

import _ "embed"

// IndexHTML은 빌드 시 embed된 index.html 내용이다.
//go:embed index.html
var IndexHTML []byte
