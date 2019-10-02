package model

type (
	Guild struct {
		ID        string `json:"id" db:"id"`
		Name      string `json:"name" db:"name"`
		Icon      string `json:"icon" db:"icon"`
		BotExists bool   `json:"-" db:"is_bot_exists"`
	}

	UserGuild struct {
		UserID      string `db:"user_id"`
		GuildID     string `json:"id" db:"guild_id"`
		IsOwner     bool   `json:"owner" db:"is_owner"`
		Permissions uint   `json:"permissions" db:"permissions"`
		CanInvite   bool   `db:"can_invite"`
	}
)

func (m *model) AddGuilds(guilds []*Guild) (err error) {
	tx, err := m.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for i := range guilds {
		_, err = tx.Exec(`INSERT INTO 
						discord_guilds (id, name, icon, is_bot_exists) 
						VALUES (?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						name = VALUES(name),
						icon = VALUES(icon),
						is_bot_exists = VALUES(is_bot_exists)`,
			guilds[i].ID,
			guilds[i].Name,
			guilds[i].Icon,
			guilds[i].BotExists)
		if err != nil {
			return
		}
	}

	return
}

func (m *model) AddGuild(guild *Guild) (err error) {
	_, err = m.db.Exec(`INSERT INTO 
						discord_guilds (id, name, icon, is_bot_exists) 
						VALUES (?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						name = VALUES(name),
						icon = VALUES(icon),
						is_bot_exists = VALUES(is_bot_exists)`,
		guild.ID,
		guild.Name,
		guild.Icon,
		guild.BotExists)
	return
}

func (m *model) AddUserGuilds(userGuilds []*UserGuild) (err error) {
	tx, err := m.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for i := range userGuilds {
		_, err = tx.Exec(`INSERT INTO users__discord_guilds
						VALUES (?, ?, ?, ?, ?)
						ON DUPLICATE KEY UPDATE
						is_owner = VALUES(is_owner),
						permissions = VALUES(permissions),
						can_invite = VALUES(can_invite)`,
			userGuilds[i].UserID,
			userGuilds[i].GuildID,
			userGuilds[i].IsOwner,
			userGuilds[i].Permissions,
			userGuilds[i].CanInvite)
		if err != nil {
			return
		}
	}

	return
}

func (m *model) GetBotExistsGuilds() ([]*Guild, error) {
	guilds := []*Guild{}
	if err := m.db.Select(&guilds, `SELECT * FROM discord_guilds WHERE is_bot_exists = true`); err != nil {
		return nil, err
	}
	return guilds, nil
}

func (m *model) UpdateGuild(guild *Guild) (err error) {
	_, err = m.db.NamedQuery(`UPDATE discord_guilds SET 
							name=:name, 
							icon=:icon, 
							is_bot_exists=:is_bot_exists
							WHERE id=:id`, &guild)
	return err
}

func (m *model) UpdateGuildBotExists(id string, exists bool) (err error) {
	_, err = m.db.Exec(`UPDATE discord_guilds SET 
						is_bot_exists=? 
						WHERE id=?`, exists, id)
	return err
}

func (m *model) GetBotAndUserExistsGuilds(userid string) ([]*Guild, error) {
	guilds := []*Guild{}
	if err := m.db.Select(&guilds, `SELECT id, name, icon FROM discord_guilds 
									WHERE is_bot_exists = true 
									AND id IN (SELECT guild_id FROM users__discord_guilds WHERE user_id = ?)`, userid); err != nil {
		return nil, err
	}
	return guilds, nil
}
