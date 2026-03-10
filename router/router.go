package router

import (
	"github.com/Shubhouy1/asset-management/handlers"
	"github.com/Shubhouy1/asset-management/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/register", handlers.RegisterUser)
	r.Post("/login", handlers.LoginUser)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/get-assets", handlers.TotalAssets)
		r.Post("/logout", handlers.LogoutUser)
		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.RequiredRoles("admin"))
			r.Delete("/delete/{id}", handlers.DeleteUser)
		})
		r.Route("/assets", func(r chi.Router) {
			r.Use(middleware.RequiredRoles("admin", "asset-manager"))
			r.Post("/", handlers.CreateAsset)

			r.Put("/assign/{id}", handlers.AssignAsset)
			r.Put("/sent-to-service/{id}", handlers.SentToService)
			r.Get("/", handlers.ShowAssets)
			r.Put("/{id}", handlers.UpdateAsset)
		})
		r.Route("/employee", func(r chi.Router) {
			r.Use(middleware.RequiredRoles("admin", "asset-manager"))
			r.Get("/", handlers.GetAllUsers)
		})
	})
	return r
}
