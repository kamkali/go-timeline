package app

import (
	_ "embed"
	"github.com/kamkali/go-timeline/internal/config"
	"github.com/kamkali/go-timeline/internal/domain"
	"golang.org/x/net/context"
	"time"
)

//go:embed logoPW.jpeg
var pwLogo []byte

func (a *app) seedDBWithAdmin(c *config.Config) error {
	u := domain.User{
		Email:    c.AdminEmail,
		Password: c.AdminPassword,
	}
	if err := a.userService.CreateUser(context.Background(), u); err != nil {
		return err
	}
	a.log.Info("Seeded the DB with admin user")
	return nil
}

func (a *app) seedDBWithExampleValues(c *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Server.TimeoutSeconds)*time.Second)
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
		Name:                "Zamknięcie Szkoły Przygotowawczej",
		EventTime:           time.Date(1831, 1, 1, 0, 0, 0, 0, time.UTC),
		ShortDescription:    "W roku 1831 została zamknięta Szkoła Przygotowawcza",
		DetailedDescription: "Za datę powstania szkolnictwa technicznego w Warszawie przyjmuje się rok 1826, w którym została otwarta Szkoła Przygotowawcza do studiów technicznych. Inicjatorem powstania szkoły i autorem programu nauczania był działający w Komisji Wyznań Religijnych i Oświecenia Publicznego Stanisław Staszic - wszechstronny uczony i działacz oświaty.",
		TypeID:              t1.ID,
	}

	e2 := &domain.Event{
		Name:                "Powstanie szkolnictwa technicznego",
		EventTime:           time.Date(1836, 1, 1, 0, 0, 0, 0, time.UTC),
		ShortDescription:    "W roku 1826 została otwarta Szkoła Przygotowawcza",
		DetailedDescription: "Po kilku zaledwie latach działania, szkoła ta została zamknięta w roku 1831, w ramach represji po wybuchu Powstania Listopadowego.",
		TypeID:              t2.ID,
	}

	e3 := &domain.Event{
		Name:                "Powstanie Politechniki Warszawskiej",
		EventTime:           time.Date(1915, 11, 1, 19, 20, 18, 0, time.UTC),
		ShortDescription:    "Data powstania PW",
		DetailedDescription: "Wybuch I wojny światowej przyniósł ze sobą możliwość realizacji postulatu nauczania we własnym języku. Po zajęciu Warszawy przez Niemców uzyskano zgodę na inaugurację działalności polskiej Politechniki, co nastąpiło w listopadzie 1915 roku. Studia prowadzone były na czterech wydziałach: Architektury, Budowy Maszyn i Elektrotechniki, Chemicznym oraz Inżynierii Budowlanej i Rolnej. Pierwszym rektorem został profesor Zygmunt Straszewicz.",
		Graphic:             pwLogo,
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
