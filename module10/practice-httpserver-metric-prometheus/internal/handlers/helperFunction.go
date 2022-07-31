package handlers

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func TestHandler(w http.ResponseWriter, res *http.Request) {
	/* Header part*/
	header := res.Header          // get Header from a request
	w.Header().Set("cncamp", "J") // self-defined

	showHeader := make(map[string]string, len(header))
	for k, v := range header {
		sb := strings.Builder{}
		for i := range v {
			sb.WriteString(v[i])
			//w.Header().Set(k, v[i])
			showHeader[k] = v[i]
		}
		w.Header().Set(k, sb.String())
	}

	/* Environment part*/
	key := "Version"
	//os.Setenv(key, "1.0")
	val, ok := os.LookupEnv(key)
	if !ok {
		out := key + " not set"
		fmt.Println(out)
		fmt.Fprintln(w, out)
		return
	}
	w.Header().Set(key, val)
	statusCode := http.StatusOK
	w.WriteHeader(statusCode)
	out := key + " = " + val
	fmt.Printf("%s = %s\n", key, val)

	io.WriteString(w, "<h1>Hello, cncamp</h1>\n")
	io.WriteString(w, "<h2>Header of Request: </h2>\n")
	for k, v := range showHeader {
		//io.WriteString(w, k+"= "+v)
		fmt.Fprintln(w, k+"= "+v)
	}
	fmt.Fprintln(w, "Status Code:", statusCode)
	fmt.Fprintf(w, "<h2>Env Variable: </h2>\n")
	fmt.Fprintln(w, out)

}

func LogHandler(code int, r *http.Request) {
	// add more level of log for check clearly
	startTimer := time.Now()

	// duration from request to response
	ip := "127.0.0.1"
	if i := r.Header.Get("x-forwarded-for"); isIP(i) && ip != i {
		ip = i
	}
	log.Printf("[%d] %s - %s - %v\n", code, ip, r.RequestURI, time.Since(startTimer))
	// Care more information about ip header
	//log.Printf("[%d] %s - %s - %v\n", code, GetClientIP(r), r.RequestURI, time.Since(startTimer))

}

/* Parse all headers' information */
// IsIP : Check if the given ip address is valid
func isIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// GetClientIPFromXForwardedFor : Parse x-forwarded-for headers.
func getClientIPFromXForwardedFor(header string) string {
	if header == "" {
		return ""
	}

	// x-forwarded-for may return multiple IP addresses in the format
	// @see https://en.wikipedia.org/wiki/X-Forwarded-For#Format
	proxies := strings.Split(header, ", ")

	var ips []string

	if len(proxies) > 0 {
		for _, proxy := range proxies {
			ip := proxy
			// make sure we only use this if it's ipv4 (ip:port)
			if strings.Contains(ip, ":") {
				splitted := strings.Split(ip, ":")
				ips = append(ips, splitted[0])
				continue
			}
			ips = append(ips, ip)
		}
	}

	// Sometimes IP addresses in this header can be 'unknown' (http://stackoverflow.com/a/11285650).
	// Therefore taking the left-most IP address that is not unknown
	// A Squid configuration directive can also set the value to "unknown" (http://www.squid-cache.org/Doc/config/forwarded_for/)
	for _, ip := range ips {
		if isIP(ip) {
			return ip
		}
	}

	return ""
}

// GetClientIP : Parse all headers.
func GetClientIP(r *http.Request) string {
	headers := r.Header

	if len(headers) > 0 {
		checklist := []string{
			"x-client-ip",         // Standard headers used by Amazon EC2, Heroku, and others.
			"x-forwarded-for",     // Load-balancers (AWS ELB) or proxies.
			"cf-connecting-ip",    // @see https://support.cloudflare.com/hc/en-us/articles/200170986-How-does-Cloudflare-handle-HTTP-Request-headers-
			"fastly-client-ip",    // Fastly and Firebase hosting header (When forwared to cloud function)
			"true-client-ip",      // Akamai and Cloudflare: True-Client-IP.
			"x-real-ip",           // Default nginx proxy/fcgi; alternative to x-forwarded-for, used by some proxies.
			"x-cluster-client-ip", // (Rackspace LB and Riverbed's Stingray) http://www.rackspace.com/knowledge_center/article/controlling-access-to-linux-cloud-sites-based-on-the-client-ip-address
			"x-forwarded",
			"forwarded-for",
			"forwarded",
		}

		for _, h := range checklist {
			if h == "x-forwarded-for" {
				if ip := getClientIPFromXForwardedFor(r.Header.Get(h)); isIP(ip) {
					return ip
				}
				continue
			}

			if ip := r.Header.Get(h); isIP(ip) {
				return ip
			}
		}
	}

	if ip := r.RemoteAddr; isIP(ip) {
		return ip
	}

	return ""
}