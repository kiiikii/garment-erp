package notifications

import (
	"fmt"
	"log"
	"net/smtp"
)

const (
	smtpHost   = "smtp.gmail.com"
	smtpPort   = "587"
	sender     = "shiejacko@gmail.com"
	password   = "evfjrguxyepbnlxg"
	ownerEmail = "shiejacko@gmail.com"
)

type NotificationService struct {
	repo *NotificationRepository
}

func NewNotificationService(repo *NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) sendEmail(to []string, message string) error {
	auth := smtp.PlainAuth("", sender, password, smtpHost)
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, to, []byte(message))
}

func (s *NotificationService) CheckAndSendAlerts() {
	log.Printf("Starting background check for delayed orders...")

	delayedOrders, err := s.repo.GetDelayedProductionOrder()
	if err != nil {
		log.Println("Failed to fetch delayed orders: ", err)
		return
	}

	for _, order := range delayedOrders {
		//! building dynamic receiver list
		receivers := []string{ownerEmail}
		if order.CustomerEmail.Valid && order.CustomerEmail.String != "" {
			receivers = append(receivers, order.CustomerEmail.String)
		}

		reasonText := ""
		if order.WaitingReason.Valid && order.WaitingReason.String != "" {
			reasonText = fmt.Sprintf("Reported Blockers: %s\n", order.WaitingReason.String)
		}

		//! draft professional email for both owner and customer
		subject := fmt.Sprintf("Subject: Production Status Update: Order #%d\n", order.OrderID)
		body := fmt.Sprintf("Dear %s & Production Team,\n\n"+
			"This is an automated system update regarding Order #%d.\n\n"+
			"Current Status: %s\n"+
			"Original Internal Target: %s\n"+
			"%s\n"+ // <-- This will insert the reason (or stay blank if none)
			"We are currently experiencing a slight delay beyond our initial internal timeline. "+
			"Our floor managers have been notified and are actively prioritizing this order to ensure timely final delivery.\n\n"+
			"Thank you,\nGarment ERP Automated System",
			order.CustomerName, order.OrderID, order.Status, order.Deadline.Format("Jan 02, 2006"), reasonText)

		err := s.sendEmail(receivers, subject+"\n"+body)

		if err == nil {
			s.repo.MarkAlertSent(order.OrderID, "PRODUCTION_DELAY")
			log.Printf("Sent delay alert for order #%d to %v", order.OrderID, receivers)
		} else {
			log.Println("Failed to send email:", err)
		}

	}
}
