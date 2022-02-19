package cron

import (
	"log"

	"github.com/eyecuelab/go-api/internal/models"
	"github.com/eyecuelab/kit/db/psql"
)

func processNotificationEmail(user *models.User) {
	defer wg.Done()
	// if err := notifications.SomethingSomething(user); err != nil {
	// 	log.Fatalf("Error: %+v\n\n", err)
	// }
	if err := psql.DB.Exec("update something set something_at = now() where id = ?", user.ID).Error; err != nil {
		log.Fatalf("Error: %+v\n\n", err)
	}
}
