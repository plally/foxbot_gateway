package role_persister

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"time"
)


type Info struct {
	InfoKey string
	InfoType  string
}

type PersistedRole struct {
	Info
	GuildID string
	RoleID string
	UserID string
	Expiration time.Time
}

func CreateRolePersist(guildID string, userID string, roleID string,  expiration time.Time) error {
	table := GetTable()
	return table.Put(PersistedRole{
		Info:       Info{
			InfoKey: "member#" + guildID + "#" + userID,
			InfoType:  "persisted_roles#"+roleID,
		},

		GuildID: guildID,
		RoleID: roleID,
		UserID: userID,
		Expiration: expiration,
	}).Run()
}

func GetPersistedRoles(guildID string, userID string) ([]PersistedRole, error) {
	table := GetTable()
	var out []PersistedRole
	err := table.Get("InfoKey", "member#" + guildID + "#" + userID).
		Range("InfoType", dynamo.BeginsWith, "persisted_roles#").
		All(&out)

	return out, err
}

func DeletePersistedRole(guildID string, userID string, roleID string) error {
	table := GetTable()

	return table.Delete("InfoKey", "member#" + guildID + "#" + userID).Range("InfoType", "persisted_roles#"+roleID).Run()
}

func GetNamedRoleID(guildID string, name string) (string, error) {
	table := GetTable()
	var data struct {
		Info
		RoleID string
	}

	err := table.Get("ItemType", "guild#"+guildID).Range("InfoType", dynamo.Equal, "named_role#"+name).One(&data)

	return data.RoleID, err
}

func GetTable() dynamo.Table {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // TODO this  region shouldnt be hardcoded
	})
	if err != nil {
		panic(err)
	}

	db := dynamo.New(s)
	return db.Table("bot_info_database")
}