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
		if saveErr := message.Save(inboxPath); saveErr != nil {
			log.Err(saveErr).Msg("unable to save message to inbox")
		}
		return err
	}

	return nil
}

func InboxCleaner(inboxPath string, sender mailer.Mailer, webhookHandler mailer.Mailer) InboxHandler {
	return func() {
		err := filepath.Walk(inboxPath, func(path string, item os.FileInfo, err error) error {
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

			cleanPath := filepath.Clean(path)
			contents, err := os.ReadFile(cleanPath) // #nosec G304 -- path from controlled inbox directory
			if err != nil {
				log.Err(err).Str("path", path).Msg("unable to read message file")
				return nil
			}

			var message mailer.Message
			err = json.Unmarshal(contents, &message)
			if err != nil {
				log.Err(err).Str("path", path).Msg("unable to load the message")
				return nil
			}

			// If it is mailer then we try both
			// Else we try only webhook sender.
			if message.Method == "mailer" {
				message.Retries += 1
				err := sender.Send(message)
				if err != nil {
					log.Err(err).Str("path", path).Msg("unable to send message")
					if saveErr := message.Save(inboxPath); saveErr != nil {
						log.Err(saveErr).Str("path", path).Msg("unable to save message")
					}
				} else {
					log.Info().Str("path", path).Msg("message successfully sent")
					if err := os.Remove(path); err != nil {
						log.Err(err).Str("path", path).Msg("unable to remove file")
					}
				}

				// Reduce counter because webhook counter should be separate
				message.Retries -= 1
				if err := sendWebhook(message, webhookHandler, inboxPath); err != nil {
					log.Err(err).Str("path", path).Msg("unable to send webhook")
				}
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
		if err != nil {
			log.Err(err).Str("path", inboxPath).Msg("unable to walk inbox")
		}
	}
}
