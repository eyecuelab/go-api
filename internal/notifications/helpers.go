package notifications

import (
	"strings"

	"github.com/eyecuelab/kit/mailman"
	"github.com/spf13/viper"
)

func send(tmpl string, name string, email string, vars *mailman.MergeVars) error {
	substr := viper.GetString("email_whitelist_pattern")
	if substr != "" && !strings.Contains(email, substr) {
		return nil
	}
	addr := &mailman.Address{Name: name, Email: email}
	if vars.BodyVars == nil {
		vars.BodyVars = map[string]string{}
	}
	vars.BodyVars["portal_url"] = viper.GetString("portal_url")

	return mailman.Send(addr, tmpl, vars)
}
