package handlers

import (
	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

func AddGroupWateringInterval(bot *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)

			return err
		}

		////temp, err := useCases.GetUserTemporary(int(context.Sender().ID))
		////if err != nil {
		////	return err
		////}
		////
		////if temp.MessageID != nil {
		////	err = context.Bot().Delete(&telebot.Message{ID: *temp.MessageID, Chat: context.Chat()})
		////	if err != nil {
		////		logger.Error("Failed to delete message", "Error", err)
		////
		////		return err
		////	}
		////}
		//
		// wateringInterval, err := strconv.Atoi(context.Data())
		// if err != nil {
		//	logger.Error("Failed to parse watering interval", "Error", err)
		//
		//	return err
		//}
		//
		//group, err := useCases.AddGroupWateringInterval(int(context.Sender().ID), wateringInterval)
		//if err != nil {
		//	return err
		//}
		//
		//// TODO переделать меню на актуальное
		//menu := &telebot.ReplyMarkup{
		//	ResizeKeyboard: true,
		//	InlineKeyboard: [][]telebot.InlineButton{
		//		{
		//			buttons.SkipGroupDescriptionButton,
		//		},
		//		{
		//			buttons.BackToAddGroupTitleButton,
		//			buttons.MenuButton,
		//		},
		//	},
		//}
		//
		//// Получаем бота, чтобы при отправке получить messageID для дальнейшего удаления:
		//msg, err := context.Bot().Send(
		//	context.Chat(),
		//	&telebot.Photo{
		//		File:    telebot.FromDisk(paths.AddGroupLastWateringDateImagePath),
		//		Caption: fmt.Sprintf(texts.AddGroupLastWateringDateText, group.Title, group.Description, group.Title),
		//	},
		//	menu,
		//)
		//if err != nil {
		//	logger.Error("Failed to send message", "Error", err)
		//
		//	return err
		//}

		return nil
	}
}
