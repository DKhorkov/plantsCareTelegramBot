package handlers

import (
	"fmt"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"gopkg.in/telebot.v4"
)

func AddGroupTitle(useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)
			return err
		}

		temp, err := useCases.GetUserTemporary(int(context.Sender().ID))
		if err != nil {
			return err
		}

		if err = context.Bot().Delete(&telebot.Message{ID: temp.MessageID, Chat: context.Chat()}); err != nil {
			logger.Error("Failed to delete message", "Error", err)
			return err
		}

		group, err := useCases.AddGroupTitle(int(context.Sender().ID), context.Message().Text)
		if err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					backToAddGroupTitleButton,
					menuButton,
				},
			},
		}

		// Получаем бота, чтобы при отправке получить messageID для дальнейшего удаления:
		msg, err := context.Bot().Send(
			context.Chat(),
			&telebot.Photo{
				File:    telebot.FromDisk(addGroupDescriptionImagePath),
				Caption: fmt.Sprintf(addGroupDescriptionText, group.Title, group.Title),
			},
			menu,
		)

		if err != nil {
			logger.Error("Failed to send message", "Error", err)
			return err
		}

		if err = useCases.SetTemporaryMessage(int(context.Sender().ID), msg.ID); err != nil {
			return err
		}

		return nil
	}
}
