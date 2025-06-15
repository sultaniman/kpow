package mailer

import (
	"errors"

	"github.com/rs/zerolog/log"
)

func SendMessage(
	message Message,
	mailer Mailer,
	webhookHandler Mailer,
	inboxPath string,
) error {
	failed := false
	if mailer != nil {
		if err := mailer.Send(message); err != nil {
			failed = true
			log.
				Err(err).
				Str("method", "mailer").
				Msg("Unable to send the message")

			message.Method = "mailer"
			message.Retries = 0
			err = message.Save(inboxPath)

			if err != nil {
				log.
					Err(err).
					Msg("Unable to save message")
			}
		}
	}

	if webhookHandler != nil {
		if err := webhookHandler.Send(message); err != nil {
			failed = true
			log.
				Err(err).
				Str("method", "webhook").
				Msg("Unable to send the message")

			message.Method = "webhook"
			message.Retries = 0
			err = message.Save(inboxPath)
			if err != nil {
				log.
					Err(err).
					Msg("Unable to save message")
			}
		}
	}

	if failed {
		return errors.New("unable to send the message")
	}

	return nil
}
