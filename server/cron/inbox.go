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
				log.Println("Unable to read file", path, err)
				return nil
			}

			if item.IsDir() {
				return filepath.SkipDir
			}

			contents, err := os.ReadFile(path)
			if err != nil {
				log.Println("Unable to read file with message", path, err)
				return nil
			}

			var message mailer.Message
			err = json.Unmarshal(contents, &message)
			if err != nil {
				log.Println("Unable to load the message", path, err)
			}

			// If it is mailer then we try both
			// Else we try only webhook sender.
			if message.Method == "mailer" {
				message.Retries += 1
				err := sender.Send(message)
				if err != nil {
					log.Println("Unable to send a message", path, err)
					message.Save(inboxPath)
				}

				// Reduce counter because webhook counter should be separate
				message.Retries -= 1
				sendWebhook(message, webhookHandler, inboxPath)
			} else {
				message.Retries += 1
				err := sendWebhook(message, webhookHandler, inboxPath)
				if err != nil {
					log.Println("unable to send webhook", path, err)
				} else {
					err := os.Remove(path)
					if err != nil {
						log.Println("unable to remove file from inbox", path, err)
					}
				}
			}

			return nil
		})
	}
}
