package middleware

import (
	"net"
	"net/http"
)

// CIDRMiddleware - middleware для валидации IP адреса
func CIDRMiddleware(trustedSubnet string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if trustedSubnet == "" {
				next.ServeHTTP(w, r)
				return
			}
			// смотрим заголовок запроса X-Real-IP
			ipStr := r.Header.Get("X-Real-IP")
			ok, err := CIDRValidate(trustedSubnet, ipStr)

			if err != nil || !ok {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CIDRValidate - проверяем входит ли указанный адрес в нужную подсеть
func CIDRValidate(cidr string, ipStr string) (bool, error) {
	ip := net.ParseIP(ipStr)

	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}

	return ipnet.Contains(ip), nil
}
