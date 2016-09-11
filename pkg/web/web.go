package web

import (
	"encoding/json"
	"net/http"

	"github.com/husobee/vestigo"
	"github.com/juju/errors"
	r "gopkg.in/dancannon/gorethink.v2"

	"github.com/ernestoalejo/tfg-fn/pkg/context"
	"github.com/ernestoalejo/tfg-fn/pkg/models"
)

func Register(r *vestigo.Router, db *r.Session) {
	r.Get("/", middleware(db, homeHandler))
	r.Get("/functions", middleware(db, functionsHandler))
	r.Get("/monitoring", middleware(db, monitoringHandler))
	r.Get("/monitoring/api", middleware(db, monitoringApiHandler))
	r.Get("/static/:file", middleware(db, staticHandler))
}

func homeHandler(db *r.Session, w http.ResponseWriter, req *http.Request) error {
	return renderTemplate(w, "home", nil)
}

func functionsHandler(db *r.Session, w http.ResponseWriter, req *http.Request) error {
	cursor, err := r.Table(models.TableFunctions).Run(db)
	if err != nil {
		return errors.Trace(err)
	}
	defer cursor.Close()
	functions := []*models.Function{}
	if err := cursor.All(&functions); err != nil {
		return errors.Trace(err)
	}

	return renderTemplate(w, "functions", map[string]interface{}{
		"Functions": functions,
	})
}

func monitoringHandler(db *r.Session, w http.ResponseWriter, req *http.Request) error {
	return renderTemplate(w, "monitoring", nil)
}

func monitoringApiHandler(db *r.Session, w http.ResponseWriter, req *http.Request) error {
	if err := json.NewEncoder(w).Encode(context.Stats()); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func staticHandler(db *r.Session, w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, r.URL.Path[1:])
	return nil
}
