package router

import (
	"log/slog"
	"main/internal/config"
	"main/internal/http-server/handler/node/add"
	deleteNode "main/internal/http-server/handler/node/delete"
	updatecontent "main/internal/http-server/handler/node/update-content"
	uploadimage "main/internal/http-server/handler/node/upload-image"
	"main/internal/http-server/middleware/authenticator"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
)

type NoteNoder interface {
	add.NodeAdder
	deleteNode.NodeDeleter
	updatecontent.NodeUpdater
	uploadimage.ImageUploader
}

func (r *Router) InitNoteNodesRoutes(storage Storage, logger *slog.Logger, cfg *config.Config) {
	// node routes
	r.Route("/node", func(nodeRouter chi.Router) {
		nodeRouter.Use(jwtauth.Verifier(r.jwtauth))
		nodeRouter.Use(authenticator.Authenticator(r.jwtauth, logger))

		// create
		nodeRouter.Post("/", add.New(logger, storage))

		// update
		nodeRouter.Patch("/{id}", updatecontent.New(logger, storage))
		nodeRouter.Patch("/{id}/image", uploadimage.New(cfg, logger, storage))

		// delete
		nodeRouter.Delete("/{id}", deleteNode.New(logger, storage))
	})
}
