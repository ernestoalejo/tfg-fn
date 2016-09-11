package web

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/juju/errors"
	r "gopkg.in/dancannon/gorethink.v2"
)

type Handler func(db *r.Session, w http.ResponseWriter, r *http.Request) error

func middleware(db *r.Session, handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(db, w, r); err != nil {
			logrus.WithFields(logrus.Fields{"error": errors.ErrorStack(err)}).Error("handler error")
			http.Error(w, errors.ErrorStack(err), http.StatusInternalServerError)
		}
	}
}

func renderTemplate(w io.Writer, name string, data interface{}) error {
	tmpl, err := template.ParseFiles("templates/base.gotmpl", fmt.Sprintf("templates/%s.gotmpl", name))
	if err != nil {
		return errors.Trace(err)
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		return errors.Trace(err)
	}

	return nil
}
