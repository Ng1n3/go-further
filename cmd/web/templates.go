package main

import "github.com/Ng1n3/go-further/pkg/models"

type templateData struct {
  Snippet *models.Snippet
  Snippets []*models.Snippet
}