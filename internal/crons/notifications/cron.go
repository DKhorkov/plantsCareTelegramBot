package notifications

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/handlers"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/utils"
)

const (
	dateFormat = "02.01.2006"
	sendHour   = 12
)

type Cron struct {
	bot      *telebot.Bot
	useCases interfaces.UseCases
	logger   logging.Logger
	stopChan chan struct{}
}

func New(bot *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) *Cron {
	return &Cron{
		bot:      bot,
		useCases: useCases,
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

func (c *Cron) Run(limit, offset int, interval time.Duration) error {
	go func() {
		if r := recover(); r != nil {
			c.logger.Error("Recovered from panic", "Recovered", r)
		}
	}()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return nil
		case <-ticker.C:
			if err := c.notify(limit, offset); err != nil {
				return err
			}
		}
	}
}

func (c *Cron) Stop() error {
	c.stopChan <- struct{}{}

	return nil
}

// TODO при необходимости добавить еще кроны - сделать этот универсальным и передавать функцию для исполнения.
// TODO добавить sync.Map для проверки отправки уведомления по сценарию в течение дня.
func (c *Cron) notify(limit, offset int) error {
	// Не отправляем уведомления до нужного часа:
	var (
		now      = time.Now()
		sendTime = time.Date(now.Year(), now.Month(), now.Day(), sendHour, 0, 0, 0, now.Location())
	)

	if now.Before(sendTime) {
		return nil
	}

	groups, err := c.useCases.GetGroupsForNotify(limit, offset)
	if err != nil {
		return err
	}

	for _, group := range groups {
		// TODO при проблеме с производительностью - сделать кэширование
		user, err := c.useCases.GetUserByID(group.UserID)
		if err != nil {
			return err
		}

		groupPlants, err := c.useCases.GetGroupPlants(group.ID)
		if err != nil {
			return err
		}

		//// Обновляем сценарий на месте, потому что растений нет и нет смысла в отправке уведомления:
		// if len(groupPlants) == 0 {
		//	group.LastWateringDate = group.NextWateringDate
		//	= group.NextWateringDate.AddDate(0, 0, group.WateringInterval)
		//	if err = c.useCases.UpdateGroup(group); err != nil {
		//		return err
		//	}
		//
		//	continue
		//}

		plantsText, err := c.preparePlantsText(groupPlants)
		if err != nil {
			c.logger.Error("Failed to prepare plants text", "Error", err)

			return err
		}

		btn := telebot.InlineButton{
			Unique: "groupWatered",
			Text:   "Растения в данном сценарии политы ✅",
			Data:   strconv.Itoa(group.ID),
		}

		c.bot.Handle(&btn, handlers.GroupWateredCallback(c.bot, c.useCases, c.logger))

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					btn,
				},
			},
		}

		msg, err := c.bot.Send(
			&telebot.Chat{ID: int64(user.TelegramID)},
			fmt.Sprintf(
				texts.Notify,
				group.Title,
				group.Description,
				group.LastWateringDate.Format(dateFormat),
				utils.GetWateringInterval(group.WateringInterval),
				plantsText,
			),
			menu,
		)
		if err != nil {
			c.logger.Error("Failed to send message", "Error", err)

			return err
		}

		notification := &entities.Notification{
			GroupID:   group.ID,
			MessageID: msg.ID,
			Text:      msg.Text,
			SentAt:    msg.Time(),
		}

		if _, err = c.useCases.SaveNotification(*notification); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cron) preparePlantsText(plants []entities.Plant) (string, error) {
	if len(plants) == 0 {
		return "В данный сценарий полива пока что не было добавлено ни одно растение!\n", nil
	}

	builder := strings.Builder{}
	for i, plant := range plants {
		_, err := builder.WriteString(fmt.Sprintf("%d) %s\n", i+1, plant.Title))
		if err != nil {
			return "", err
		}
	}

	return builder.String(), nil
}
