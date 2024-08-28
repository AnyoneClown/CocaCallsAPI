package handlers

import "github.com/AnyoneClown/CocaCallsAPI/storage"

type DefaultHandler struct {
	Storage *storage.CockroachDB
}
