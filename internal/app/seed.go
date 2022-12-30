package app

import (
	_ "embed"
	"github.com/kamkali/go-timeline/internal/config"
	"github.com/kamkali/go-timeline/internal/domain"
	"golang.org/x/net/context"
	"time"
)

//go:embed apollo11.jpeg
var apollo11Graphic []byte

func (a *app) seedDBWithAdmin(c *config.Config) error {
	u := domain.User{
		Email:    c.Admin.Email,
		Password: c.Admin.Password,
	}
	if err := a.userService.CreateUser(context.Background(), u); err != nil {
		return err
	}
	a.log.Info("Seeded the DB with admin user")
	return nil
}

func (a *app) seedDBWithExampleValues() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	t1 := &domain.Type{
		Name:  "normal",
		Color: "white",
	}
	t2 := &domain.Type{
		Name:  "error",
		Color: "red",
	}
	t3 := &domain.Type{
		Name:  "special",
		Color: "green",
	}
	for _, t := range []*domain.Type{t1, t2, t3} {
		id, err := a.typeService.CreateType(ctx, t)
		if err != nil {
			return err
		}
		t.ID = id
	}

	e1 := &domain.Event{
		Name:                "Launch of Sputnik 1",
		EventTime:           time.Date(1957, 10, 4, 0, 0, 0, 0, time.UTC),
		ShortDescription:    "The Soviet Union launches the first artificial satellite",
		DetailedDescription: "Sputnik 1 was a Soviet artificial satellite that was launched into orbit on October 4, 1957. It was the first satellite to be launched into space, and its successful launch marked the beginning of the Space Age. The satellite was about the size of a beach ball and weighed just under 200 pounds. It was equipped with two radio transmitters, which sent out a series of beeps that could be heard by amateur radio operators around the world. The launch of Sputnik 1 sparked the Space Race between the Soviet Union and the United States, which ultimately led to the first human landing on the Moon in 1969.",
		TypeID:              t1.ID,
	}

	e2 := &domain.Event{
		Name:                "First human spaceflight",
		EventTime:           time.Date(1961, 4, 12, 9, 07, 0, 0, time.UTC),
		ShortDescription:    "Yuri Gagarin becomes the first human to orbit Earth",
		DetailedDescription: "Yuri Gagarin was a Soviet pilot and cosmonaut who became the first human to orbit Earth on April 12, 1961. Gagarin's spacecraft, Vostok 1, circled the Earth once in just under an hour and a half, reaching an altitude of about 186 miles. Gagarin's flight marked a major milestone in the Space Race between the Soviet Union and the United States, and he became an international celebrity after his return to Earth. Gagarin died in a plane crash in 1968 at the age of 34.",
		TypeID:              t2.ID,
	}

	e3 := &domain.Event{
		Name:                "Apollo 11 Moon landing",
		EventTime:           time.Date(1969, 7, 20, 20, 18, 0, 0, time.UTC),
		ShortDescription:    "Neil Armstrong becomes the first human to set foot on the Moon",
		DetailedDescription: `The Apollo 11 mission was the first manned mission to land on the Moon. It was launched on July 16, 1969, and four days later, on July 20, astronauts Neil Armstrong and Edwin "Buzz" Aldrin landed the lunar module Eagle on the surface of the Moon. Armstrong became the first human to set foot on the Moon when he stepped out of the lunar module and onto the surface, saying the famous words, "That's one small step for man, one giant leap for mankind." The Apollo 11 mission marked the end of the Space Race between the United States and the Soviet Union, and it remains one of the most significant achievements in the history of space exploration.`,
		Graphic:             string(apollo11Graphic),
		TypeID:              t3.ID,
	}

	for _, e := range []*domain.Event{e1, e2, e3} {
		id, err := a.eventService.CreateEvent(ctx, e)
		if err != nil {
			return err
		}
		e.ID = id
	}

	return nil
}
