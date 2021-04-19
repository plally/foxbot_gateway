package role_persister

import (
	"fmt"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	err := CreateRolePersist(
		"794721717209530368",
		"832140863231492116",
		"687411481314459663",
		time.Now().Add(time.Hour*24),
	)
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	roles, err := GetPersistedRoles("420427961867567124", "345010836089339906")
	if err != nil {
		t.Error(err)
	}

	for _, role := range roles {
		fmt.Println(role.InfoKey)
	}
}

func TestDeletePersistedRole(t *testing.T) {
	err := DeletePersistedRole(
		"794721717209530368",
		"687411481314459663",
		"832140863231492116",
		)
	if err != nil {
		t.Error(err)
	}
}