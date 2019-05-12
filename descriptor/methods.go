package descriptor

import "fmt"

func (d *ServiceDescriptor) AddRoute(route Route) error {
	for _, r := range d.Routes {
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

	d.Routes = append(d.Routes, route)

	return nil
}

func (d *ServiceDescriptor) AddMiddleware(mw Middleware) error {
	for _, m := range d.Middlwares {
		if m.HandlerName == mw.HandlerName {
			return fmt.Errorf("handler with the same name already exists: %v", m.HandlerName)
		}
	}

	d.Middlwares = append(d.Middlwares, mw)

	return nil
}
