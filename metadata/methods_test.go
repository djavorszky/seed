package metadata

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceDescriptor_AddRoute(t *testing.T) {
	defInfo := Info{
		Name:        "default",
		Description: "defaultDesc",
		Summary:     "defaultSummary",
	}

	tests := []struct {
		name      string
		addRoutes []Route
		route     Route
		wantErr   bool
	}{
		{
			name:      "All right",
			addRoutes: []Route{},
			route: Route{
				HandlerName: "testName",
				Path:        "/",
				HttpMethods: []string{http.MethodGet},
				StrictSlash: false,
				Info:        defInfo,
			},
			wantErr: false,
		},
		{
			name: "same path, different method",
			addRoutes: []Route{
				{
					HandlerName: "someRandomNAme",
					Path:        "/",
					HttpMethods: []string{http.MethodGet},
					StrictSlash: false,
					Info:        defInfo,
				},
			},
			route: Route{
				HandlerName: "testName",
				Path:        "/",
				HttpMethods: []string{http.MethodPost},
				StrictSlash: false,
				Info:        defInfo,
			},
			wantErr: false,
		},
		{
			name: "same path, same method",
			addRoutes: []Route{
				{
					HandlerName: "someRandomNAme",
					Path:        "/",
					HttpMethods: []string{http.MethodGet},
					StrictSlash: false,
					Info:        defInfo,
				},
			},
			route: Route{
				HandlerName: "testName",
				Path:        "/",
				HttpMethods: []string{http.MethodGet},
				StrictSlash: false,
				Info:        defInfo,
			},
			wantErr: true,
		},
		{
			name: "same path, overlapping methods",
			addRoutes: []Route{
				{
					HandlerName: "someRandomNAme",
					Path:        "/",
					HttpMethods: []string{http.MethodGet, http.MethodPost},
					StrictSlash: false,
					Info:        defInfo,
				},
			},
			route: Route{
				HandlerName: "testName",
				Path:        "/",
				HttpMethods: []string{http.MethodPost, http.MethodPut},
				StrictSlash: false,
				Info:        defInfo,
			},
			wantErr: true,
		},
		{
			name: "same path, different methods",
			addRoutes: []Route{
				{
					HandlerName: "someRandomNAme",
					Path:        "/",
					HttpMethods: []string{http.MethodGet, http.MethodPut},
					StrictSlash: false,
					Info:        defInfo,
				},
			},
			route: Route{
				HandlerName: "testName",
				Path:        "/",
				HttpMethods: []string{http.MethodPost, http.MethodDelete},
				StrictSlash: false,
				Info:        defInfo,
			},
			wantErr: false,
		},
		{
			name: "same handlerName",
			addRoutes: []Route{
				{
					HandlerName: "testName",
					Path:        "/someother",
					HttpMethods: []string{http.MethodGet},
					StrictSlash: false,
					Info:        defInfo,
				},
			},
			route: Route{
				HandlerName: "testName",
				Path:        "/",
				HttpMethods: []string{http.MethodGet},
				StrictSlash: false,
				Info:        defInfo,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Base(Info{})

			d.Routes = append(d.Routes, tt.addRoutes...)

			if err := d.AddRoute(tt.route); (err != nil) != tt.wantErr {
				t.Errorf("Metadata.AddRoute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			assert.Equal(t, len(d.Routes), len(tt.addRoutes)+1, "length mismatch")
		})
	}
}

func TestServiceDescriptor_AddMiddleware(t *testing.T) {
	defInfo := Info{
		Name:        "default",
		Description: "defaultDesc",
		Summary:     "defaultSummary",
	}

	tests := []struct {
		name           string
		addMiddlewares []Middleware
		mw             Middleware
		wantErr        bool
	}{
		{
			name:           "All right",
			addMiddlewares: []Middleware{},
			mw: Middleware{
				HandlerName: "someRandomNAme",
				Paths:       []string{"/"},
				Priority:    1,
				Info:        defInfo,
			},
			wantErr: false,
		},
		{
			name: "same handler name",
			addMiddlewares: []Middleware{
				{
					HandlerName: "someRandomNAme",
					Paths:       []string{"/"},
					Priority:    1,
					Info:        defInfo,
				},
			},
			mw: Middleware{
				HandlerName: "someRandomNAme",
				Paths:       []string{"/"},
				Priority:    1,
				Info:        defInfo,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Base(Info{})

			d.Middlwares = append(d.Middlwares, tt.addMiddlewares...)

			if err := d.AddMiddleware(tt.mw); (err != nil) != tt.wantErr {
				t.Errorf("Metadata.AddRoute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			assert.Equal(t, len(d.Middlwares), len(tt.addMiddlewares)+1, "length mismatch")
		})
	}
}
