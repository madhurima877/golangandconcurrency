package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("method=%s path=%s duration=%s",
			r.Method,
			r.URL.Path,
			time.Since(start))
	})

}

// func Logger() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		start := time.Now()

// 		c.Next()
// 		c.Set()
// 		c.Get()
// 		log.Printf(
// 			"method=%s path=%s duration=%s",
// 			c.Request.Method,
// 			c.Request.URL.Path,
// 			time.Since(start),
// 		)
// 	}
// }
