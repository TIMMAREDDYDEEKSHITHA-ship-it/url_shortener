//user structure

//this file defines what a user is

package main

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users"`
	ID            string `json:"id" bun:"id,pk"`
	Name          string `json:"name"`
	Email         string `json:"email"`
}
