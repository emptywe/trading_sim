package router

import (
	"net/http"
)

func (h *Handler) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		return
	})
}

//
//import (
//	"log"
//	"net/http"
//	"strings"
//)
//
//
//func (h *Handler) Middleware(handler http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		headInfo := r.Header.Get("Authorization")
//		if r.URL.Path == "/login" {
//
//			if headInfo == "" {
//				handler.ServeHTTP(w, r)
//				return
//			}
//			tokenStruct := strings.Split(headInfo, " ")
//			if len(tokenStruct) < 2 {
//				handler.ServeHTTP(w, r)
//				return
//			}
//			token := tokenStruct[1]
//
//			err := h.services.ValidateToken(token)
//			if err == nil {
//				log.Println("Already authorized")
//				w.WriteHeader(400)
//				_ = json.NewEncoder(w).Encode(map[string]string{"code": "400", "message": "Already authorized", "status": "Error"})
//				return
//			} else {
//				r.Header.Set("Authorization", "")
//				handler.ServeHTTP(w, r)
//				return
//			}
//		}
//
//		if r.URL.Path == "/refresh" {
//			handler.ServeHTTP(w, r)
//			return
//		}
//
//		if headInfo == "" {
//			log.Println("Empty Authorization Header")
//			w.WriteHeader(400)
//			_ = json.NewEncoder(w).Encode(map[string]string{"code": "400", "message": "Empty Authorization Header", "status": "Error"})
//			return
//		}
//
//		tokenStruct := strings.Split(headInfo, " ")
//		if len(tokenStruct) < 2 {
//			log.Println("Invalid token structure")
//			w.WriteHeader(400)
//			_ = json.NewEncoder(w).Encode(map[string]string{"code": "400", "message": "Invalid token structure", "status": "Error"})
//			return
//		}
//		token := tokenStruct[1]
//		err := h.services.ValidateToken(token)
//
//		if err != nil {
//			log.Println(err)
//			if err.Error() == "Token is expired" {
//				log.Println("Empty Authorization Header")
//				w.WriteHeader(300)
//				_ = json.NewEncoder(w).Encode(map[string]string{"code": "300", "message": "Refresh token", "status": "Error"})
//				return
//			} else {
//				log.Println("Invalid token")
//				w.WriteHeader(400)
//				_ = json.NewEncoder(w).Encode(map[string]string{"code": "400", "message": "Invalid token", "status": "Error"})
//				return
//			}
//		}
//
//		handler.ServeHTTP(w, r)
//		return
//	})
//}
//
//func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
//	headInfo := r.Header.Get("Authorization")
//	if headInfo == "" {
//		log.Println("Empty Authorization Header")
//		w.WriteHeader(400)
//		_ = json.NewEncoder(w).Encode(map[string]string{"code": "400", "message": "Empty Authorization Header", "status": "Error"})
//		return
//	}
//	tokenStruct := strings.Split(headInfo, " ")
//	expiredToken := tokenStruct[1]
//
//	var tokenReqBody struct {
//		RefreshToken string `json:"refresh_token"`
//	}
//
//	err := json.NewDecoder(r.Body).Decode(&tokenReqBody)
//	if err != nil {
//		log.Println(err.Error())
//		w.WriteHeader(400)
//		_ = json.NewEncoder(w).Encode(map[string]string{"code": "400", "message": "Can't parse request", "status": "Error"})
//		return
//	}
//
//	err = h.services.ValidateToken(tokenReqBody.RefreshToken)
//	if err != nil {
//		log.Println(err.Error())
//		w.WriteHeader(400)
//		_ = json.NewEncoder(w).Encode(map[string]string{"code": "400", "message": "Invalid token", "status": "Error"})
//		return
//	}
//
//	sUUID, UserName, Uid, err := h.services.CheckSession(context.Background(), tokenReqBody.RefreshToken, expiredToken)
//	if err != nil {
//		log.Println(err.Error())
//		w.WriteHeader(400)
//		_ = json.NewEncoder(w).Encode(map[string]string{"code": "400", "message": "Invalid token", "status": "Error"})
//		return
//	}
//
//	token, err := h.services.Auth.ReGenerateToken(sUUID, UserName, Uid)
//	if err != nil {
//		log.Println(err.Error())
//		w.WriteHeader(401)
//		_ = json.NewEncoder(w).Encode(map[string]string{"code": "401", "message": "User Unauthorized", "status": "Error"})
//		return
//	}
//
//	_ = json.NewEncoder(w).Encode(token)
//
//}
