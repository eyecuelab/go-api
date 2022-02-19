package email

import (
	"github.com/eyecuelab/kit/mailman"
	"github.com/eyecuelab/kit/mailman/sendgrid"
	"github.com/spf13/viper"
)

// InitEmail mailer configuration init
func InitEmail() error {
	sendgrid.Configure(sendgrid.MinConfig(
		viper.GetString("from_email_name"),
		viper.GetString("from_email_address"),
		viper.GetString("email_domain"),
	))
	return mailman.RegisterTemplate("forgot", "confirm_email")
}
