package lib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	InsertRSVP  = `INSERT INTO rsvp (content, ctime) VALUES (?, ?);`
	GetAllRSVPs = `SELECT id, content, ctime FROM rsvp;`

	rsvpValidationErr = "invalid RSVP: %s"
)

type DinnerType int

// ToDo - the dinners should probably be in their own DB table...but this will do
const (
	UnknownDinner = DinnerType(iota)
	BeefShortRib
	HoneySalmon
	Vegetarian
	Vegan
)

func (d DinnerType) ToString() string {
	switch d {
	case BeefShortRib:
		return "Short Rib"
	case HoneySalmon:
		return "Salmon"
	case Vegetarian:
		return "Vegetarian"
	case Vegan:
		return "Vegan"
	default:
		return "Unknown"
	}
}

type RSVP struct {
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	IsAttending  bool       `json:"is_attending"`
	DinnerChoice DinnerType `json:"dinner_choice"`
	Comments     string     `json:"comments,omitempty"`
	Guests       []*PlusOne `json:"guests,omitempty"`
	WantsAccomm  bool       `json:"accommodations"`
}

func (r *RSVP) Validate() error {
	if r.Name == "" {
		return fmt.Errorf(rsvpValidationErr, "missing Name")
	}
	if r.Email == "" {
		return fmt.Errorf(rsvpValidationErr, "missing Email Address")
	}
	if r.DinnerChoice == UnknownDinner {
		return fmt.Errorf(rsvpValidationErr, "invalid dinner selection")
	}
	if len(r.Guests) > 1 {
		return fmt.Errorf(rsvpValidationErr, "cannot have more than a single +1 guest")
	}
	for _, guest := range r.Guests {
		if err := guest.Validate(); err != nil {
			return fmt.Errorf("invalid RSVP guest: %w", err)
		}
	}

	return nil
}

type PlusOne struct {
	Name         string     `json:"name"`
	DinnerChoice DinnerType `json:"dinner_choice"`
	IsAttending  bool       `json:"is_attending"`
}

func (p *PlusOne) Validate() error {
	if p.Name == "" {
		return errors.New("missing Name")
	}
	if p.DinnerChoice == UnknownDinner {
		return errors.New("invalid dinner selection")
	}

	return nil
}

type RSVPRow struct {
	Id      int       `db:"id"`
	Content string    `db:"content"`
	CTime   time.Time `db:"ctime"`
}

func (r *RSVPRow) toRSVP() (*RSVP, error) {
	var rsvp RSVP
	if err := json.Unmarshal([]byte(r.Content), &rsvp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal RSVP Row content: %w", err)
	}

	return &rsvp, nil
}

type RSVPProvider struct {
	client *DBClient
}

func NewRSVPProvider(client *DBClient) *RSVPProvider {
	return &RSVPProvider{
		client: client,
	}
}

func (r *RSVPProvider) GetAll(ctx context.Context) ([]*RSVP, error) {
	var rows []RSVPRow
	if err := r.client.Db.SelectContext(ctx, &rows, GetAllRSVPs); err != nil {
		return nil, fmt.Errorf("failed to get RSVP records: %w", err)
	}

	var rsvps []*RSVP
	for _, row := range rows {
		rsvp, err := row.toRSVP()
		if err != nil {
			return nil, err
		}

		rsvps = append(rsvps, rsvp)
	}

	// Return an empty list instead of nil
	if rsvps == nil {
		return []*RSVP{}, nil
	}
	return rsvps, nil
}

func (r *RSVPProvider) Add(ctx context.Context, toAdd *RSVP) error {
	marshaled, err := json.Marshal(toAdd)
	if err != nil {
		return fmt.Errorf("failed to marshal RSVP entry for add: %w", err)
	}

	_, err = r.client.Db.ExecContext(ctx, InsertRSVP, marshaled, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert RSVP record: %w", err)
	}

	return nil
}
