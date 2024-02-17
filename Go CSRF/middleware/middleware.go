package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bhupesh998/go-csrf/db"
	"github.com/bhupesh998/go-csrf/middleware/myJWT"
	// "github.com/bhupesh998/go-csrf/middleware/myJWT.go"
	"github.com/bhupesh998/go-csrf/server/templates"
	"github.com/justinas/alice"
)

func Newhandler() http.Handler {
	return alice.New(recoverHandler, authHandler).ThenFunc(logicHandler)
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Panic("Recovered! panic:%+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func authHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/restricted", "/logout", "/deleteUser":
			log.Println("IN AUTH RESTRICTED SECTION")
			AuthCookie, authErr := r.Cookie("AuthToken")
			if authErr == http.ErrNoCookie {
				log.Println("Unauthorized attempt! no auth cookie")
				nullifyTokenCookies(&w, r)
				http.Error(w, http.StatusText(401), 401)
				return
			} else if authErr != nil {
				log.Panic("panic: %+v", authErr)
				nullifyTokenCookies(&w, r)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			RefreshCookie, refreshErr := r.Cookie("RefreshToken")
			if refreshErr == http.ErrNoCookie {
				log.Println("Unauthorized attempt! no refresh cookie")
				nullifyTokenCookies(&w, r)
				http.Redirect(w, r, "/login", 302)

				return
			} else if refreshErr != nil {
				log.Panic("panic: %+v", refreshErr)
				nullifyTokenCookies(&w, r)
				http.Error(w, http.StatusText(500), 500)
				return
			}

			requestCsrfToken := grabCsrfFromReq(r)
			log.Println(requestCsrfToken)

			authTokenString, refreshTokenString, csrfSecret, err := myJWT.CheckAndRefreshTokens(AuthCookie.Value, RefreshCookie.Value, requestCsrfToken)
			if err != nil {
				if err.Error() == "Unauthorized" {
					log.Println("Unauthorized Attempt JWT is not valid")
					http.Error(w, http.StatusText(401), 401)
					return
				} else {
					log.Panic("Error: %+v", err)
					http.Error(w, http.StatusText(500), 500)
					return
				}
			}
			log.Println("Successfully reCreated JWT's")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			setAuthAndRefreshCookies(&w, authTokenString, refreshTokenString)
			w.Header().Set("X-CSRF-Token", csrfSecret)

		default:
			// no check necessary
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func logicHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/restricted":
		csrfSecret := grabCsrfFromReq(r)
		templates.RenderTemplate(w, "restricted", &templates.RestrictedPage{ csrfSecret, "Hello Buddy!" })
	case "/login":
		switch r.Method {
		case "GET":
			templates.RenderTemplate(w, "login", &templates.LoginPage{false, ""})
		case "POST":
			r.ParseForm()
			log.Println(r.Form)

			user , uuid , loginErr := db.LogUserIn(strings.Join(r.Form["username"], ""),strings.Join(r.Form["password"], "") )
			log.Println(user, uuid , loginErr)

			if loginErr != nil {
				w.WriteHeader(http.StatusUnauthorized)
			}else{
			authTokenString, refreshTokenString , csrfSecret, err :=myJWT.CreateNewTokens(uuid, user.Role)
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
			}

			setAuthAndRefreshCookies(&w, authTokenString, refreshTokenString)
			w.Header().Set("X-CSRF-Token", csrfSecret)
			w.WriteHeader(http.StatusOK)
		}

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "/register":
		switch r.Method {
		case "GET":
			templates.RenderTemplate(w, "register", &templates.RegisterPage{false, ""})
		case "POST":
			r.ParseForm()
			log.Println("Form data register" , r.Form)
			_, uuid, err := db.FetchUserByUsername(strings.Join(r.Form["username"], ""))

			if err == nil {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				role := "user"
				uuid, err = db.StoreUser(strings.Join(r.Form["username"], ""), strings.Join(r.Form["password"], ""), role)

				if err != nil {
					http.Error(w, http.StatusText(500), 500)
				}
				log.Println("uuid is ", uuid)

				authTokenString, refreshTokenString, csrfSecret, err := myJWT.CreateNewTokens(uuid, role)
				if err != nil {
					http.Error(w, http.StatusText(500), 500)
				}
				setAuthAndRefreshCookies(&w, authTokenString, refreshTokenString)
				log.Println("Before setting CSRF token in register route")
				w.Header().Set("X-CSRF-Token", csrfSecret)
				w.WriteHeader(http.StatusOK)
			}

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	case "/logout":
		nullifyTokenCookies(&w, r)
		http.Redirect(w, r, "/login", 302)
	case "/deleteUser":
		log.Println("Deleting the user")
		AuthCookie , authErr := r.Cookie("AuthToken")
		if authErr == http.ErrNoCookie {
			log.Println("Unauthorized attempt! no auth cookie")
			nullifyTokenCookies(&w, r)
			http.Redirect(w, r , "/login", 302)
			return
		} else if authErr != nil {
			log.Panic("panic: %+v", authErr)
			nullifyTokenCookies(&w, r)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		uuid, uuidErr := myJWT.GrabUUID(AuthCookie.Value)
		if(uuidErr != nil){
			log.Panic("panic: %+v", uuidErr)
			nullifyTokenCookies(&w, r)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		db.DeleteUser(uuid)
		nullifyTokenCookies(&w, r)
		http.Redirect(w, r , "/register", 302)
		
	default:
		w.WriteHeader(http.StatusOK)
	}

}

func nullifyTokenCookies(w *http.ResponseWriter, r *http.Request) {
	authCookie := http.Cookie{
		Name:     "AuthToken",
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(*w, &authCookie)

	refreshCookie := http.Cookie{
		Name:     "RefreshToken",
		Value:    "",
		Expires:  time.Now().Add(-1000 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(*w, &refreshCookie)

	RefreshCookie, RefreshErr := r.Cookie("RefreshToken")
	if RefreshErr == http.ErrNoCookie {
		return
	} else if RefreshErr != nil {
		log.Panic("panic :%+v", RefreshErr)
		http.Error(*w, http.StatusText(500), 500)

	}

	myJWT.RevokeRefreshToken(RefreshCookie.Value)
}

func setAuthAndRefreshCookies(w *http.ResponseWriter, authTokenString string, refreshTokenString string) {
	authCookie := http.Cookie{
		Name:     "AuthToken",
		Value:    authTokenString,
		HttpOnly: true,
	}
	http.SetCookie(*w, &authCookie)

	refreshCookie := http.Cookie{
		Name:     "RefreshToken",
		Value:    refreshTokenString,
		HttpOnly: true,
	}
	http.SetCookie(*w, &refreshCookie)

}

func grabCsrfFromReq(r *http.Request) string {
	csrfFromFrom := r.FormValue("X-CSRF-Token")

	if csrfFromFrom != "" {
		return csrfFromFrom
	} else {
		return r.Header.Get("X-CSRF-Token")
	}

}
