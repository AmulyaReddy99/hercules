package department

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/kshitij10496/hercules/common"
)

var Routes = common.Routes{
	common.Route{
		Name:        "Department Info",
		Method:      "GET",
		Pattern:     "/info",
		HandlerFunc: ServiceDepartment.departmentsHandler,
		PathPrefix:  common.VERSION + "/department",
	},
}

// serviceDepartment implements the server interface
//
type serviceDepartment common.Service

func (s serviceDepartment) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: What should go in here?
	log.Printf("[request initiate] %s - %v\n", s.Name, r.URL)
	s.Router.ServeHTTP(w, r)
	log.Printf("[request end] %s - %v\n", s.Name, r.URL)
}

func (s serviceDepartment) GetDBConnection(ctx context.Context) (*sql.Conn, error) {
	return s.DB.Conn(ctx)
}

func (s serviceDepartment) GetURL() string {
	return s.URL
}

func (s serviceDepartment) SetDB(db *sql.DB) common.Server {
	s.DB = db
	return s
}

// ServiceDepartment represents the course service.
var ServiceDepartment serviceDepartment

// Initialise the service with no DB.
func init() {
	ServiceDepartment = serviceDepartment{
		Name:   "service-course",
		URL:    "/course",
		DB:     nil,
		Router: common.NewSubRouter(Routes),
	}
}
