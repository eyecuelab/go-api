package cron

import (
	"log"

	"github.com/eyecuelab/kit/db/psql"
)

func processSomethingExample() {
	wg.Add(1)
	defer wg.Done()
	// if err := notifications.SomethingSomething(user); err != nil {
	// 	log.Fatalf("Error: %+v\n\n", err)
	// }
	if err := psql.DB.Exec("update users set updated_at = now()").Error; err != nil {
		log.Fatalf("Error: %+v\n\n", err)
	}
}
