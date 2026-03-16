package cron

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/sultaniman/kpow/server/mailer"
)

type InboxHandler func()

func sendWebhook(message mailer.Message, webhookHandler mailer.Mailer, inboxPath string) error {
	if webhookHandler == nil {
		return nil
	}

	err := webhookHandler.Send(message)
	if err != nil {
		message.Retries += 1
		log.Err(err).Msg("webhook delivery failed")
		message.Save(inboxPath)
		return err
	}

	return nil
}

func InboxCleaner(inboxPath string, sender mailer.Mailer, webhookHandler mailer.Mailer) InboxHandler {
	return func() {
		filepath.Walk(inboxPath, func(path string, item os.FileInfo, err error) error {
			if err != nil {
				log.Err(err).Str("path", path).Msg("unable to read file")
				return nil
			}
			if item.IsDir() && path != inboxPath {
				return filepath.SkipDir
			}

			if item.IsDir() && path == inboxPath {
				return nil
			}

			contents, err := os.ReadFile(path)
			if err != nil {
				log.Err(err).Str("path", path).Msg("unable to read message file")
				return nil
			}

			var message mailer.Message
			err = json.Unmarshal(contents, &message)
			if err != nil {
				log.Err(err).Str("path", path).Msg("unable to load the message")
			}

			// If it is mailer then we try both
			// Else we try only webhook sender.
			if message.Method == "mailer" {
				message.Retries += 1
				err := sender.Send(message)
				if err != nil {
					log.Err(err).Str("path", path).Msg("unable to send message")
					message.Save(inboxPath)
				} else {
					log.Info().Str("path", path).Msg("message successfully sent")
					if err := os.Remove(path); err != nil {
						log.Err(err).Str("path", path).Msg("unable to remove file")
					}
				}

				// Reduce counter because webhook counter should be separate
				message.Retries -= 1
				sendWebhook(message, webhookHandler, inboxPath)
			} else {
				message.Retries += 1
				err := sendWebhook(message, webhookHandler, inboxPath)
				if err != nil {
					log.Err(err).Str("path", path).Msg("unable to send webhook")
				} else {
					log.Info().Str("path", path).Msg("webhook successfully sent")
					if err := os.Remove(path); err != nil {
						log.Err(err).Str("path", path).Msg("unable to remove file")
					}
				}
			}

			return nil
		})
	}
}
