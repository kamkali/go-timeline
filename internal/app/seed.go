package app

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"github.com/kamkali/go-timeline/internal/config"
	timeline2 "github.com/kamkali/go-timeline/internal/timeline"
	"golang.org/x/net/context"
	"time"
)

var (
	//go:embed apollo11.jpeg
	apollo11Graphic []byte
	//go:embed apollo11-crew.jpeg
	apollo11CrewGraphic []byte
)

func (a *app) seedDBWithAdmin(c *config.Config) error {
	u := timeline2.User{
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
	t1 := &timeline2.Type{
		Name:  "normal",
		Color: "white",
	}
	t2 := &timeline2.Type{
		Name:  "error",
		Color: "red",
	}
	t3 := &timeline2.Type{
		Name:  "special",
		Color: "green",
	}
	for _, t := range []*timeline2.Type{t1, t2, t3} {
		id, err := a.typeService.CreateType(ctx, t)
		if err != nil {
			return err
		}
		t.ID = id
	}

	e1 := &timeline2.Event{
		Name:                "Launch of Sputnik 1",
		EventTime:           time.Date(1957, 10, 4, 0, 0, 0, 0, time.UTC),
		ShortDescription:    "The Soviet Union launches the first artificial satellite",
		DetailedDescription: "Sputnik 1 was a Soviet artificial satellite that was launched into orbit on October 4, 1957. It was the first satellite to be launched into space, and its successful launch marked the beginning of the Space Age. The satellite was about the size of a beach ball and weighed just under 200 pounds. It was equipped with two radio transmitters, which sent out a series of beeps that could be heard by amateur radio operators around the world. The launch of Sputnik 1 sparked the Space Race between the Soviet Union and the United States, which ultimately led to the first human landing on the Moon in 1969.",
		TypeID:              t1.ID,
	}

	e2 := &timeline2.Event{
		Name:                "First human spaceflight",
		EventTime:           time.Date(1961, 4, 12, 9, 07, 0, 0, time.UTC),
		ShortDescription:    "Yuri Gagarin becomes the first human to orbit Earth",
		DetailedDescription: "Yuri Gagarin was a Soviet pilot and cosmonaut who became the first human to orbit Earth on April 12, 1961. Gagarin's spacecraft, Vostok 1, circled the Earth once in just under an hour and a half, reaching an altitude of about 186 miles. Gagarin's flight marked a major milestone in the Space Race between the Soviet Union and the United States, and he became an international celebrity after his return to Earth. Gagarin died in a plane crash in 1968 at the age of 34.",
		TypeID:              t2.ID,
	}

	e3 := &timeline2.Event{
		Name:                "Apollo 11 Moon landing",
		EventTime:           time.Date(1969, 7, 20, 20, 18, 0, 0, time.UTC),
		ShortDescription:    "Neil Armstrong becomes the first human to set foot on the Moon",
		DetailedDescription: `The Apollo 11 mission was the first manned mission to land on the Moon. It was launched on July 16, 1969, and four days later, on July 20, astronauts Neil Armstrong and Edwin "Buzz" Aldrin landed the lunar module Eagle on the surface of the Moon. Armstrong became the first human to set foot on the Moon when he stepped out of the lunar module and onto the surface, saying the famous words, "That's one small step for man, one giant leap for mankind." The Apollo 11 mission marked the end of the Space Race between the United States and the Soviet Union, and it remains one of the most significant achievements in the history of space exploration.`,
		TypeID:              t3.ID,
	}

	e4 := &timeline2.Event{
		Name:                "Launch of Apollo 8",
		EventTime:           time.Date(1968, 12, 21, 12, 51, 0, 0, time.UTC),
		ShortDescription:    "First manned spacecraft to leave Earth's orbit and reach the Moon",
		DetailedDescription: "Apollo 8 was the first manned spacecraft to leave Earth's orbit, and it successfully orbited the Moon on December 24, 1968. The crew of Apollo 8, consisting of Frank Borman, James Lovell, and William Anders, became the first humans to see the far side of the Moon. This mission was a major stepping stone in the Apollo program, paving the way for the eventual landing on the Moon.",
		TypeID:              t3.ID,
	}

	e5 := &timeline2.Event{
		Name:                "Launch of Apollo 11",
		EventTime:           time.Date(1969, 7, 16, 9, 32, 0, 0, time.UTC),
		ShortDescription:    "Launch of the Apollo 11 mission to the Moon",
		DetailedDescription: "Apollo 11 was the fifth manned mission in the Apollo program and the first mission to land humans on the Moon. The mission was launched on July 16, 1969, and was crewed by Neil Armstrong, Buzz Aldrin, and Michael Collins. The mission is most famous for the first human landing on the Moon on July 20, 1969, when Armstrong and Aldrin became the first humans to walk on the Moon's surface.",
		Graphic:             fmt.Sprintf("data:text/plain;base64,%s", base64.StdEncoding.EncodeToString(apollo11Graphic)),
		TypeID:              t3.ID,
	}

	e6 := &timeline2.Event{
		Name:                "Lunar Landing Training Vehicle Crash",
		EventTime:           time.Date(1968, 5, 6, 15, 35, 13, 0, time.UTC),
		ShortDescription:    "Crash of the Lunar Landing Training Vehicle during a training mission",
		DetailedDescription: "During a training mission for the Apollo 11 landing, the Lunar Landing Training Vehicle (LLTV) being piloted by Neil Armstrong experienced a malfunction and crashed to the ground. Armstrong was able to eject from the vehicle just seconds before impact, saving his life. The incident raised concerns about the safety of the Apollo program and highlighted the risks involved in space exploration.",
		TypeID:              t3.ID,
	}

	e7 := &timeline2.Event{
		Name:                "Return of Apollo 11 Crew to Earth",
		EventTime:           time.Date(1969, 7, 24, 16, 50, 35, 0, time.UTC),
		ShortDescription:    "Apollo 11 crew successfully returns to Earth",
		DetailedDescription: "After spending eight days in space, the Apollo 11 crew successfully returned to Earth on July 24, 1969. The crew, consisting of Neil Armstrong, Buzz Aldrin, and Michael Collins, were hailed as heroes upon their return and were celebrated around the world for their historic achievement.",
		Graphic:             fmt.Sprintf("data:text/plain;base64,%s", base64.StdEncoding.EncodeToString(apollo11CrewGraphic)),
		TypeID:              t3.ID,
	}

	for _, e := range []*timeline2.Event{e1, e2, e3, e4, e5, e6, e7} {
		id, err := a.eventService.CreateEvent(ctx, e)
		if err != nil {
			return err
		}
		e.ID = id
	}

	return nil
}
