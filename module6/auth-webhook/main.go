package main

import (
	"encoding/json"
	"github.com/xanzy/go-gitlab"
	authv1 "k8s.io/api/authentication/v1"
	authentication "k8s.io/api/authentication/v1beta1"
	"log"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// Check User
		log.Print("receving request")
		var opts []gitlab.ClientOptionFunc
		//if g.opts.BaseUrl != "" {
		//	opts = append(opts, gitlab.WithBaseURL(g.opts.BaseUrl))
		//}
		token := <PAT>
		client, err := gitlab.NewClient(token, opts...)
		if err != nil {
			log.Fatalln(err.Error())
		}

		user, _, err := client.Users.CurrentUser()
		if err != nil {
			log.Fatalln(err.Error())
		}

		resp := &authv1.UserInfo{
			Username: user.Username,
			UID:      strconv.Itoa(user.ID),
		}
		//log.Println(resp)
		w.WriteHeader(http.StatusOK)
		trs := authentication.TokenReviewStatus{
			Authenticated: true,
			User: authentication.UserInfo{
				Username: resp.Username,
				UID:      resp.UID,
			},
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"apiVersion": "authentication.k8s.io/v1beta1",
			"kind":       "TokenReview",
			"status":     trs,
		})
	})
	http.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var tr authentication.TokenReview
		err := decoder.Decode(&tr)
		if err != nil {
			log.Println("[Error]", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]interface{}{
				"apiVersion": "authentication.k8s.io/v1beta1",
				"kind":       "TokenReview",
				"status": authentication.TokenReviewStatus{
					Authenticated: false,
				},
			})
			if err != nil {
				return
			}
			return
		}
		// Check User
		log.Print("receving request")
		var opts []gitlab.ClientOptionFunc
		//if g.opts.BaseUrl != "" {
		//	opts = append(opts, gitlab.WithBaseURL(g.opts.BaseUrl))
		//}
		//token := ":q"
		client, err := gitlab.NewClient(tr.Spec.Token, opts...)
		if err != nil {
			log.Println(err.Error())
		}

		user, _, err := client.Users.CurrentUser()
		if err != nil {
			log.Println(err.Error())
		}

		resp := &authv1.UserInfo{
			Username: user.Username,
			UID:      strconv.Itoa(user.ID),
		}
		log.Println(resp)
		//ts := oauth2.StaticTokenSource(
		//	&oauth2.Token{AccessToken: tr.Spec.Token},
		//)
		//tc := oauth2.NewClient(context.Background(), ts)
		//client := github.NewClient(tc)
		//user, _, err := client.Users.Get(context.Background(), "")
		//if err != nil {
		//	log.Println("[Error]", err.Error())
		//	w.WriteHeader(http.StatusUnauthorized)
		//	err := json.NewEncoder(w).Encode(map[string]interface{}{
		//		"apiVersion": "authentication.k8s.io/v1beta1",
		//		"kind":       "TokenReview",
		//		"status": authentication.TokenReviewStatus{
		//			Authenticated: false,
		//		},
		//	})
		//	if err != nil {
		//		return
		//	}
		//	return
		//}

		log.Printf("[Success] login as %s", resp.Username)
		w.WriteHeader(http.StatusOK)
		trs := authentication.TokenReviewStatus{
			Authenticated: true,
			User: authentication.UserInfo{
				Username: resp.Username,
				UID:      resp.UID,
			},
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"apiVersion": "authentication.k8s.io/v1beta1",
			"kind":       "TokenReview",
			"status":     trs,
		})
	})
	log.Println(http.ListenAndServe(":3000", nil))
}
