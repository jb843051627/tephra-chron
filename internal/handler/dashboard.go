package handler

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(page))
}

const page = `<!doctype html><html><head><title>Tephra Chron</title></head><body><main><h1>Tephra Chron</h1><p>年代序列执行状态</p><pre id="metrics"></pre></main><script>fetch('/api/metrics').then(r=>r.json()).then(v=>document.querySelector('#metrics').textContent=JSON.stringify(v,null,2))</script></body></html>`
