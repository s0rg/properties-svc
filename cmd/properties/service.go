package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type service struct {
	addr      string
	dbUser    *sql.DB
	dbSetting *sql.DB
	h         handler
}

type apiReq struct {
	UserID int        `json:"user_id"`
	Items  []string   `json:"items"`
	Expire *time.Time `json:"expire,omitempty"`
}

func mAPI(method string, fn func(w io.Writer, r *http.Request) int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			buf  bytes.Buffer
			code int
		)

		if r.Method != method {
			code = http.StatusMethodNotAllowed
			http.Error(w, http.StatusText(code), code)

			return
		}

		if code = fn(&buf, r); code != 0 {
			http.Error(w, http.StatusText(code), code)

			return
		}

		w.Header().Set("Content-Type", "application/json")

		if _, err := buf.WriteTo(w); err != nil {
			log.Println("api response error:", err)
		}
	}
}

func mREQ(next func(w io.Writer, r *apiReq) int) func(w io.Writer, r *http.Request) int {
	return func(w io.Writer, r *http.Request) int {
		var rq apiReq

		if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
			return http.StatusBadRequest
		}

		if rq.UserID == 0 || len(rq.Items) == 0 {
			return http.StatusBadRequest
		}

		return next(w, &rq)
	}
}

func newService(addr string, dbu, dbs *sql.DB) *service {
	return &service{
		addr:      addr,
		dbUser:    dbu,
		dbSetting: dbs,
	}
}

func (svc *service) handleGetSettings(w io.Writer, r *http.Request) int {
	var uIDStr string

	if uIDStr = strings.TrimSpace(r.URL.Path[len("/settings/"):]); uIDStr == "" {
		return http.StatusBadRequest
	}

	uid, err := strconv.Atoi(uIDStr)
	if err != nil {
		return http.StatusBadRequest
	}

	when := time.Now()

	if whs := r.URL.Query().Get("when"); whs != "" {
		when, err = time.Parse(time.RFC3339, whs)
		if err != nil {
			log.Println("get date parse error:", err)

			return http.StatusBadRequest
		}
	}

	ctx := context.Background()

	res, err := svc.h.GetSettings(ctx, uid, when)
	if err != nil {
		log.Println("get handler error:", err)

		return http.StatusInternalServerError
	}

	_ = json.NewEncoder(w).Encode(res)

	return 0
}

func (svc *service) handleListSettings(w io.Writer, _ *http.Request) int {
	ctx := context.Background()

	res, err := svc.h.ListSettings(ctx)
	if err != nil {
		log.Println("list-settings handler error:", err)

		return http.StatusInternalServerError
	}

	_ = json.NewEncoder(w).Encode(res)

	return 0
}

func (svc *service) handleListBundles(w io.Writer, _ *http.Request) int {
	ctx := context.Background()

	res, err := svc.h.ListBundles(ctx)
	if err != nil {
		log.Println("list-bundles handler error:", err)

		return http.StatusInternalServerError
	}

	_ = json.NewEncoder(w).Encode(res)

	return 0
}

func (svc *service) handleListTags(w io.Writer, _ *http.Request) int {
	ctx := context.Background()

	res, err := svc.h.ListTags(ctx)
	if err != nil {
		log.Println("list-tags handler error:", err)

		return http.StatusInternalServerError
	}

	_ = json.NewEncoder(w).Encode(res)

	return 0
}

func (svc *service) handleSetTag(w io.Writer, req *apiReq) int {
	ctx := context.Background()

	if err := svc.h.SetTag(ctx, req.UserID, req.Items[0], req.Expire); err != nil {
		log.Println("set-tag handler error:", err)

		return http.StatusInternalServerError
	}

	return http.StatusCreated
}

func (svc *service) handleSetBundle(w io.Writer, req *apiReq) int {
	ctx := context.Background()

	if err := svc.h.SetBundles(ctx, req.UserID, req.Items, req.Expire); err != nil {
		log.Println("set-bundle handler error:", err)

		return http.StatusInternalServerError
	}

	return http.StatusCreated
}

func (svc *service) handleUnSetTag(w io.Writer, req *apiReq) int {
	ctx := context.Background()

	if err := svc.h.UnSetTag(ctx, req.UserID, req.Items[0]); err != nil {
		log.Println("unset-tag handler error:", err)

		return http.StatusInternalServerError
	}

	return http.StatusCreated
}

func (svc *service) handleUnSetBundle(w io.Writer, req *apiReq) int {
	ctx := context.Background()

	if err := svc.h.UnSetBundles(ctx, req.UserID, req.Items); err != nil {
		log.Println("unset-bundle handler error:", err)

		return http.StatusInternalServerError
	}

	return http.StatusCreated
}

func (svc *service) Serve() error {
	svc.h.user = NewUserStore(svc.dbUser)
	svc.h.setting = NewSettingStore(svc.dbSetting)

	http.HandleFunc("/settings", mAPI(http.MethodGet, svc.handleListSettings))
	http.HandleFunc("/bundles", mAPI(http.MethodGet, svc.handleListBundles))
	http.HandleFunc("/tags", mAPI(http.MethodGet, svc.handleListTags))

	http.HandleFunc("/settings/", mAPI(http.MethodGet, svc.handleGetSettings))

	http.HandleFunc("/set-tag", mAPI(http.MethodPost, mREQ(svc.handleSetTag)))
	http.HandleFunc("/set-bundle", mAPI(http.MethodPost, mREQ(svc.handleSetBundle)))

	http.HandleFunc("/unset-tag", mAPI(http.MethodPost, mREQ(svc.handleUnSetTag)))
	http.HandleFunc("/unset-bundle", mAPI(http.MethodPost, mREQ(svc.handleUnSetBundle)))

	srv := http.Server{
		Addr:         svc.addr,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return srv.ListenAndServe()
}
