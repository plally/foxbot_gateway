package main

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/plally/foxbot_gateway/role_persister"
	"log"
	"net/http"
	"time"
)

func onMemberAdd(session *discordgo.Session, newMemberEvent *discordgo.GuildMemberAdd) {
	persistedRoles, err := role_persister.GetPersistedRoles(newMemberEvent.GuildID, newMemberEvent.User.ID)
	if err != nil {
		log.Println("error getting persisted roles", err)
	}


	for _, role := range persistedRoles {
		err = session.GuildMemberRoleAdd(role.GuildID, role.UserID, role.RoleID)

		if time.Now().After(role.Expiration)  {
			err = role_persister.DeletePersistedRole(role.GuildID, role.UserID, role.RoleID)
			if err != nil {
				log.Println("error getting persisted role", err)
			}
		}

		var restError discordgo.RESTError
		if errors.As(err, &restError) {
			if restError.Response.StatusCode == http.StatusForbidden {
				// TODO notify server owner that the bot cant add roles
			}
		}

		if err != nil {
			log.Println("error adding role", err)
		}
	}
}