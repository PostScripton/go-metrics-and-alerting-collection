package middlewares

import (
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

// TrustedSubnet проверяет IP адрес отправителя, что он входит в группу доверенных
func TrustedSubnet(cidr string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if cidr == "" {
				next.ServeHTTP(w, r)
				return
			}

			_, ipNet, err := net.ParseCIDR(cidr)
			if err != nil {
				log.Error().Err(err).Msg("Parsing CIDR")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			realIP := r.Header.Get("X-Real-IP")
			log.Debug().Str("real-ip", realIP).Msg("IP from agent")

			ip := net.ParseIP(realIP)
			if !ipNet.Contains(ip) {
				log.Warn().Str("mask", cidr).Str("ip", ip.String()).Msg("Untrusted subnet")
				http.Error(w, "Untrusted subnet", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
