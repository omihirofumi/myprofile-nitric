package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nitrictech/go-sdk/faas"
	"github.com/nitrictech/go-sdk/nitric"
)

func main() {
	profilesApi, err := nitric.NewApi("public")
	if err != nil {
		return
	}

	profiles, err := nitric.NewCollection("profiles").With(nitric.CollectionReading, nitric.CollectionWriting)
	if err != nil {
		return
	}

	profilesApi.Post("/profiles", func(ctx *faas.HttpContext, next faas.HttpHandler) (*faas.HttpContext, error) {
		id := uuid.New().String()

		var profileRequest map[string]interface{}
		err := json.Unmarshal(ctx.Request.Data(), &profileRequest)
		if err != nil {
			return ctx, err
		}

		err = profiles.Doc(id).Set(ctx.Request.Context(), profileRequest)
		if err != nil {
			return ctx, err
		}

		ctx.Response.Body = []byte(id)

		return ctx, nil
	})

	profilesApi.Get("/profiles/:id", func(ctx *faas.HttpContext, next faas.HttpHandler) (*faas.HttpContext, error) {
		profiles, err := profiles.Query().Fetch(ctx.Request.Context())
		if err != nil {
			return ctx, err
		}

		var profileContent []map[string]interface{}
		for _, doc := range profiles.Documents {
			profileContent = append(profileContent, doc.Content())
		}

		ctx.Response.Body, err = json.Marshal(profileContent)

		return ctx, err
	})

	profilesApi.Delete("/profiles/:id", func(ctx *faas.HttpContext, next faas.HttpHandler) (*faas.HttpContext, error) {
		id := ctx.Request.PathParams()["id"]

		err := profiles.Doc(id).Delete(ctx.Request.Context())
		if err != nil {
			ctx.Response.Status = 404
			ctx.Response.Body = []byte(fmt.Sprintf("profile with id '%s' not found", id))

			return ctx, err
		}

		return ctx, err
	})

	if err := nitric.Run(); err != nil {
		fmt.Println(err)
	}

}
