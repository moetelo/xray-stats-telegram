package models

import "xray-stats-telegram/internal"

type UsersJson struct {
	Admins                *[]int64          `json:"admins,omitempty"`
	TelegramIdToXrayEmail *map[int64]string `json:"usersToXrayEmail,omitempty"`
}

type UserState struct {
	admins          *Set[int64]
	tgIdToXrayEmail *map[int64]string
}

func (m *UserState) IsAdmin(id int64) bool {
	_, ok := (*m.admins)[id]
	return ok
}

func (m *UserState) GetXrayEmail(id int64) (string, bool) {
	email, ok := (*m.tgIdToXrayEmail)[id]
	return email, ok
}

func (m *UserState) GetAllUsers() *[]string {
	return internal.Values(m.tgIdToXrayEmail)
}

func NewState(m *UsersJson) *UserState {
	admins := make(Set[int64])
	for _, admin := range *m.Admins {
		admins[admin] = struct{}{}
	}

	return &UserState{
		admins:          &admins,
		tgIdToXrayEmail: m.TelegramIdToXrayEmail,
	}
}
