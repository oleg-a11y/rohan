package service

import (
	"github.com/robfig/cron/v3"
	"rohan/internal/logger"
	"time"
)

type CronService struct {
	notionService   *NotionService
	telegramService *TelegramService
	log             *logger.Logger
}

func NewCronService(notionService *NotionService, telegramService *TelegramService, log *logger.Logger) *CronService {
	return &CronService{
		notionService:   notionService,
		telegramService: telegramService,
		log:             log,
	}
}

func (c *CronService) Start() {
	c.scheduleDailyInterviews()
	c.scheduleUpcomingInterviews()
}

func (c *CronService) scheduleDailyInterviews() {
	cronJob := cron.New()
	cronJob.AddFunc("30 9 * * *", func() {
		todayFilter := map[string]interface{}{
			"property": "Дата",
			"date": map[string]string{
				"equals": time.Now().Format("2006-01-02"),
			},
		}

		interviews, err := c.notionService.FetchInterviews(todayFilter)
		if err != nil {
			c.log.Error("Error fetching today's interviews: " + err.Error())
			return
		}

		if err := c.telegramService.SendInterviews(interviews, false); err != nil {
			c.log.Error("Error sending today's interviews to Telegram: " + err.Error())
		}
	})
	cronJob.Start()
}

func (c *CronService) scheduleUpcomingInterviews() {
	cronJob := cron.New()
	cronJob.AddFunc("* * * * *", func() {
		currentTime := time.Now()
		upcomingTime := currentTime.Add(10 * time.Minute).Format("15:04")
		todayFilter := map[string]interface{}{
			"property": "Дата",
			"date": map[string]string{
				"equals": currentTime.Format("2006-01-02"),
			},
		}

		interviews, err := c.notionService.FetchInterviews(todayFilter)
		if err != nil {
			c.log.Error("Error fetching today's interviews: " + err.Error())
			return
		}

		var upcomingInterviews []Interview
		for _, interview := range interviews {
			if interview.Date == upcomingTime {
				upcomingInterviews = append(upcomingInterviews, interview)
			}
		}

		if len(upcomingInterviews) > 0 {
			if err := c.telegramService.SendInterviews(upcomingInterviews, true); err != nil {
				c.log.Error("Error sending upcoming interviews to Telegram: " + err.Error())
			}
		}
	})
	cronJob.Start()
}
