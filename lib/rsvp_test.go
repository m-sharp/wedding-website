package lib

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

var (
	testCfg = map[string]string{
		DBHost:     "host.docker.internal",
		DBUsername: "root",
		DBPass:     "REDACTED",
		DBPort:     "3306",
	}
)

func TestRSVPProvider_Add(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	cfg := &Config{cfg: testCfg}

	dbClient, err := NewDBClient(cfg, log)
	if err != nil {
		t.Fatalf("Failed to open DB client: %s", err.Error())
	}

	inst := NewRSVPProvider(dbClient)

	if err := inst.Add(ctx, &RSVP{
		FirstName:    "Mike",
		LastName:     "Harp",
		Email:        "msharp185@gmail.com",
		IsAttending:  true,
		DinnerChoice: 1,
		Comments:     "Hello there",
		Guests: []*PlusOne{
			{
				FirstName:    "L",
				LastName:     "G",
				DinnerChoice: 2,
			},
		},
	}); err != nil {
		t.Fatalf("Failed to insert RSVP record: %s", err.Error())
	}

	if err := inst.Add(ctx, &RSVP{
		FirstName:    "Neal",
		LastName:     "Tyson",
		Email:        "durr@hotmail.com",
		IsAttending:  false,
		DinnerChoice: 3,
		Comments:     "Sorry bout ya",
		Guests: []*PlusOne{
			{
				FirstName:    "John",
				LastName:     "McClain",
				DinnerChoice: 2,
			},
		},
	}); err != nil {
		t.Fatalf("Failed to insert RSVP record: %s", err.Error())
	}
}

func TestRSVPProvider_GetAll(t *testing.T) {
	ctx := context.Background()
	log := zap.NewNop()
	cfg := &Config{cfg: testCfg}

	dbClient, err := NewDBClient(cfg, log)
	if err != nil {
		t.Fatalf("Failed to open DB client: %s", err.Error())
	}

	inst := NewRSVPProvider(dbClient)

	rsvps, err := inst.GetAll(ctx)
	if err != nil {
		t.Fatalf("Failed to get RSVP records: %s", err.Error())
	}

	for i, rsvp := range rsvps {
		t.Logf(
			"Record #%d\n\tFirst Name: %s\n\tLast Name: %s\n\tEmail: %s\n\tAttending?: %t\n\tDinner: %d\n\tComments: %q",
			i, rsvp.FirstName, rsvp.LastName, rsvp.Email, rsvp.IsAttending, rsvp.DinnerChoice, rsvp.Comments,
		)

		for j, guest := range rsvp.Guests {
			t.Logf(
				"Guest #%d\n\tFirst Name: %s\n\tLast Name: %s\n\tDinner: %d",
				j, guest.FirstName, guest.LastName, guest.DinnerChoice,
			)
		}
	}
}
