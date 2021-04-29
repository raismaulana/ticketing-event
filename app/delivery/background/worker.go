package background

import (
	"log"
	"math/rand"

	"github.com/raismaulana/ticketing-event/app/helper"
	"github.com/raismaulana/ticketing-event/app/usecase"
)

type BackgroundTask interface {
	SendReminderPayment()
	SendPromotionEvent()
}

type backgroundTask struct {
	backgroundCase usecase.BackgroundCase
}

func NewBackgroundTask(backgroundCase usecase.BackgroundCase) BackgroundTask {
	return &backgroundTask{
		backgroundCase: backgroundCase,
	}
}

func (b *backgroundTask) SendReminderPayment() {
	// a := []string{"a", "b", "c"}
	// for _, v := range a {
	// 	log.Println(v)
	// 	time.Sleep(time.Second * 2)
	// }
	log.Println("Reminder start for today")
	res, err := b.backgroundCase.GetPendingTransaction()
	if err != nil {
		log.Println("worker.go : ", err)
	}
	for _, v := range res {
		helper.SendMail(v.Email, "Reminder Payment", "v.Eid will be held tommorow, complete your payment now to get the webinar's link.")
		log.Println("send reminder payment to ", v.Email)
	}
	log.Println("Reminder end for today")
}

func (b *backgroundTask) SendPromotionEvent() {
	// a := []string{"1", "2", "3"}
	// for _, v := range a {
	// 	log.Println(v)
	// 	time.Sleep(time.Second * 1)
	// }
	events, users, err := b.backgroundCase.GetPromotionEvent()
	log.Println(events)
	log.Println(users)
	log.Println(err)
	log.Println("Promotion start for today")
	if err == nil && len(events) > 0 && len(users) > 0 {
		for _, v := range users {
			i := rand.Intn(len(events))
			event := events[i]
			helper.SendMail(v.Email, "Cool Event Will Be Held", "Join this event "+event.TitleEvent+".")
			log.Println("send ", event.ID, " to ", v.Email)
		}
	}
	log.Println("Promotion end for today")

}
