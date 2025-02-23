package yahoo

import "net/http"

// emulate a browser by appending headers like a browser to a request
func addBrowserHeaders(r *http.Request) {
	r.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36")
	r.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	r.Header.Add("accept-language", "en-US,en;q=0.9")
	r.Header.Add("cache-control", "max-age=0")
	r.Header.Add("priority", "u=0, i")
	r.Header.Add("sec-ch-ua", "\"Chromium\";v=\"130\", \"Google Chrome\";v=\"130\", \"Not?A_Brand\";v=\"99\"")
	r.Header.Add("sec-ch-ua-mobile", "?0")
	r.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	r.Header.Add("sec-fetch-dest", "document")
	r.Header.Add("sec-fetch-mode", "navigate")
	r.Header.Add("sec-fetch-site", "none")
	r.Header.Add("sec-fetch-user", "?1")
	r.Header.Add("upgrade-insecure-requests", "1")
	r.Header.Add("referrerPolicy", "strict-origin-when-cross-origin")
}
