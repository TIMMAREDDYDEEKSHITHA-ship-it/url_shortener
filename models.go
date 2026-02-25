package main

import (
	"time"

	"github.com/uptrace/bun"
)

type URL struct {
	bun.BaseModel `bun:"table:urls"`

	ID        int64     `bun:",pk,autoincrement"`
	Code      string    `bun:",unique,notnull"`
	LongURL   string    `bun:",notnull"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
