package preparers

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

var notifiedGroups = make(map[int]time.Time)

type NotificationsPreparer struct {
	bot      *telebot.Bot
	useCases interfaces.UseCases
	logger   logging.Logger
	limit    int
	offset   int
}

func NewNotificationsPreparer(
	bot *telebot.Bot,
	useCases interfaces.UseCases,
	logger logging.Logger,
	limit int,
	offset int,
) *NotificationsPreparer {
	return &NotificationsPreparer{
		bot:      bot,
		useCases: useCases,
		logger:   logger,
		limit:    limit,
		offset:   offset,
	}
}

func (p *NotificationsPreparer) GetCallback() interfaces.Callback {
	return func() error {
		if !p.canNotifyByTime() {
			return nil
		}

		groups, err := p.useCases.GetGroupsForNotify(p.limit, p.offset)
		if err != nil {
			return err
		}

		for _, group := range groups {
			if p.alreadyNotified(group) {
				continue
			}

			if err = p.notify(group); err != nil {
				return err
			}
		}

		return nil
	}
}

func (p *NotificationsPreparer) canNotifyByTime() bool {
	// Не отправляем уведомления до нужного часа:
	now := time.Now()
	sendTimeThreshold := time.Date(now.Year(), now.Month(), now.Day(), sendHour, 0, 0, 0, now.Location())

	return now.After(sendTimeThreshold)
}

func (p *NotificationsPreparer) alreadyNotified(group entities.Group) bool {
	date, exists := notifiedGroups[group.ID]
	if !exists {
		return false
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	return date.After(today)
}

func (p *NotificationsPreparer) notify(group entities.Group) error {
	// TODO при проблеме с производительностью - сделать кэширование
	user, err := p.useCases.GetUserByID(group.UserID)
	if err != nil {
		return err
	}

	groupPlants, err := p.useCases.GetGroupPlants(group.ID)
	if err != nil {
		return err
	}

	//// Обновляем сценарий на месте, потому что растений нет и нет смысла в отправке уведомления:
	// if len(groupPlants) == 0 {
	//	group.LastWateringDate = group.NextWateringDate
	//	= group.NextWateringDate.AddDate(0, 0, group.WateringInterval)
	//	if err = p.useCases.UpdateGroup(group); err != nil {
	//		return err
	//	}
	//
	//	continue
	//}

	plantsText, err := p.preparePlantsText(groupPlants)
	if err != nil {
		p.logger.Error("Failed to prepare plants text", "Error", err)

		return err
	}

	btn := telebot.InlineButton{
		Unique: "groupWatered",
		Text:   "Растения в данном сценарии политы ✅",
		Data:   strconv.Itoa(group.ID),
	}

	p.bot.Handle(&btn, handlers.GroupWateredCallback(p.bot, p.useCases, p.logger))

	menu := &telebot.ReplyMarkup{
		ResizeKeyboard: true,
		InlineKeyboard: [][]telebot.InlineButton{
			{
				btn,
			},
		},
	}

	msg, err := p.bot.Send(
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
		p.logger.Error("Failed to send message", "Error", err)

		return err
	}

	// Сохраняем информаци об отправке уведомления пользователю по сценарию:
	notifiedGroups[group.ID] = time.Now()

	notification := &entities.Notification{
		GroupID:   group.ID,
		MessageID: msg.ID,
		Text:      msg.Text,
		SentAt:    msg.Time(),
	}

	if _, err = p.useCases.SaveNotification(*notification); err != nil {
		return err
	}

	return nil
}

func (p *NotificationsPreparer) preparePlantsText(plants []entities.Plant) (string, error) {
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
