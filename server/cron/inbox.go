package cron

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/sultaniman/kpow/server/mailer"
)

type InboxHander func()

func sendWebhook(message mailer.Message, webhookHandler mailer.Mailer, inboxPath string) error {
	if webhookHandler == nil {
		return nil
	}

	err := webhookHandler.Send(message)
	if err != nil {
		message.Retries += 1
		log.Println("err", err)
		message.Save(inboxPath)
		return err
	}

	return nil
}

func InboxCleaner(inboxPath string, sender mailer.Mailer, webhookHandler mailer.Mailer) InboxHander {
	return func() {
		filepath.Walk(inboxPath, func(path string, item os.FileInfo, err error) error {
			if err != nil {
				log.Println(err)
				return nil
			}

			if item.IsDir() {
				return filepath.SkipDir
			}

			contents, err := os.ReadFile(path)
			if err != nil {
				log.Println("err", err)
				return nil
			}

			var message mailer.Message
			err = json.Unmarshal(contents, &message)
			if err != nil {
				log.Println("err", err)
			}

			// If it is mailer then we try both
			// Else we try only webhook sender.
			if message.Method == "mailer" {
				message.Retries += 1
				err := sender.Send(message)
				if err != nil {
					log.Println("err", err)
					message.Save(inboxPath)
				}

				// Reduce counter because webhook counter should be separate
				message.Retries -= 1
				sendWebhook(message, webhookHandler, inboxPath)
			} else {
				message.Retries += 1
				err := sendWebhook(message, webhookHandler, inboxPath)
				if err != nil {
					log.Println("err", err)
				}
			}

			return nil
		})
	}
}
