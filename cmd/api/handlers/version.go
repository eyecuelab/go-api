package handlers

import (
	"net/http"

	"github.com/eyecuelab/go-api/version"
	"github.com/eyecuelab/kit/web"
)

// Version api version handler
func Version(c web.ApiContext) error {
	data := struct {
		Version   string `json:"version"`
		GitRev    string `json:"gitRev"`
		Timestamp string `json:"timeStamp"`
	}{
		version.Version,
		version.GitRev,
		version.Date,
	}

	return c.JSON(http.StatusOK, &data)
}
