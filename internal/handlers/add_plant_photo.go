package handlers

import (
	"fmt"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/paths"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
)

func AcceptAddPlantPhotoCallback(
	_ *telebot.Bot,
	useCases interfaces.UseCases,
	logger logging.Logger,
) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)

			return err
		}

		temp, err := useCases.GetUserTemporary(int(context.Sender().ID))
		if err != nil {
			return err
		}

		// Получаем растение для корректного отображения данных прошлых этапов:
		plant, err := temp.GetPlant()
		if err != nil {
			logger.Error("Failed to get Plant from Temporary", "Error", err)

			return err
		}

		group, err := useCases.GetGroup(plant.GroupID)
		if err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.BackToAddPlantPhotoQuestionButton,
					buttons.MenuButton,
				},
			},
		}

		// Получаем бота, чтобы при отправке получить messageID для дальнейшего удаления:
		msg, err := context.Bot().Send(
			context.Chat(),
			&telebot.Photo{
				File: telebot.FromDisk(paths.AddPlantPhotoImagePath),
				Caption: fmt.Sprintf(
					texts.AddPlantPhotoText,
					plant.Title,
					plant.Description,
					group.Title,
					plant.Title,
				),
			},
			menu,
		)
		if err != nil {
			logger.Error("Failed to send message", "Error", err)

			return err
		}

		// TODO при проблемах логики следует сделать в рамках транзакции
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.AddPlantPhotoStep); err != nil {
			return err
		}

		if err = useCases.SetTemporaryMessage(int(context.Sender().ID), &msg.ID); err != nil {
			return err
		}

		return nil
	}
}

func BackToAddPlantPhotoQuestionCallback(
	_ *telebot.Bot,
	useCases interfaces.UseCases,
	logger logging.Logger,
) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)

			return err
		}

		temp, err := useCases.GetUserTemporary(int(context.Sender().ID))
		if err != nil {
			return err
		}

		// Получаем растение для корректного отображения данных прошлых этапов:
		plant, err := temp.GetPlant()
		if err != nil {
			logger.Error("Failed to get Plant from Temporary", "Error", err)

			return err
		}

		group, err := useCases.GetGroup(plant.GroupID)
		if err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.AcceptAddPlantPhotoButton,
				},
				{
					buttons.RejectAddPlantPhotoButton,
				},
				{
					buttons.BackToAddPlantGroupButton,
					buttons.MenuButton,
				},
			},
		}

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromDisk(paths.AddPlantPhotoQuestionImagePath),
				Caption: fmt.Sprintf(
					texts.AddPlantPhotoQuestionText,
					plant.Title,
					plant.Description,
					group.Title,
					group.Title,
					plant.Title,
				),
			},
			menu,
		)
		if err != nil {
			logger.Error("Failed to send message", "Error", err)

			return err
		}

		// TODO при проблемах логики следует сделать в рамках транзакции
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.AddPlantPhotoQuestionStep); err != nil {
			return err
		}

		// Обнуляем сообщение для удаления:
		if err = useCases.SetTemporaryMessage(int(context.Sender().ID), nil); err != nil {
			return err
		}

		return nil
	}
}
