package notifications

import (
	"strings"

	"github.com/eyecuelab/go-api/internal/models"
	"github.com/eyecuelab/kit/mailman"
)

// Welcome user signup welcome email
func Welcome(u *models.User) error {
	vars := new(mailman.MergeVars)
	vars.SubjectVars = u.MergeVars()
	vars.BodyVars = u.MergeVars()

	return send("welcome", u.FullName(), u.Email, vars)
}

// ForgotPassword forgot password email
func ForgotPassword(u *models.User) error {
	vars := new(mailman.MergeVars)
	data := u.ForgotMergeVars()
	data["reset_url"] = strings.Replace(u.URL, ":token", data["token"], 1)
	vars.SubjectVars = data
	vars.BodyVars = data

	return send("forgot", u.FullName(), u.Email, vars)
}
