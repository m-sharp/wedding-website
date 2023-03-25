package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const (
	InsertRSVP  = `INSERT INTO rsvp (content, ctime) VALUES (?, ?);`
	GetAllRSVPs = `SELECT id, content, ctime FROM rsvp;`
)

type DinnerType int

// ToDo - add dinner DB table and expose via API for form building? Or just hardcode everything?
const (
	UnknownDinner = DinnerType(iota)
	DinnerA
	DinnerB
	DinnerC
)

type RSVP struct {
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Email        string     `json:"email"`
	IsAttending  bool       `json:"is_attending"`
	DinnerChoice DinnerType `json:"dinner_choice"`
	Comments     string     `json:"comments"`
	Guests       []*PlusOne `json:"guests"`
}

type PlusOne struct {
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	DinnerChoice DinnerType `json:"dinner_choice"`
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
