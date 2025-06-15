package cron

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/goforj/godump"
	"github.com/sultaniman/kpow/server/mailer"
)

type InboxHander func()

func InboxCleaner(inboxPath string, sender mailer.Mailer, webhookHandler mailer.Mailer) InboxHander {
	return func() {
		filepath.Walk(inboxPath, func(path string, item os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
				return nil
			}
			if !item.IsDir() {
				contents, err := os.ReadFile(path)
				if err != nil {
					log.Println("err", err)
				} else {
					var message mailer.Message
					err = json.Unmarshal(contents, &message)
					if err != nil {
						log.Println("err", err)
					}
					if message.Method == "mailer" {
						sender.Send(message)
					}
					godump.Dump(message)
				}
				log.Println(path, contents)
			}
			return nil
		})
	}
}
