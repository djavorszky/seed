package metadata

import "fmt"

func (m *Metadata) AddRoute(route Route) error {
	for _, r := range m.Routes {
		if r.HandlerName == route.HandlerName {
			return fmt.Errorf("handler with the same name already exists: %v", r.HandlerName)
		}

		if r.Path != route.Path {
			continue
		}

		for _, routeMethod := range route.HttpMethods {
			for _, haveMethod := range r.HttpMethods {
				if routeMethod == haveMethod {
					return fmt.Errorf("route with path %q and method %q already used", r.Path, haveMethod)
				}
			}
		}
	}

	m.Routes = append(m.Routes, route)

	return nil
}

func (m *Metadata) AddMiddleware(mw Middleware) error {
	for _, m := range m.Middlwares {
		if m.HandlerName == mw.HandlerName {
			return fmt.Errorf("handler with the same name already exists: %v", m.HandlerName)
		}
	}

	m.Middlwares = append(m.Middlwares, mw)

	return nil
}
