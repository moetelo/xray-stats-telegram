package models

import (
	"os"
	"strconv"
	"strings"
	"xray-stats-telegram/internal"
)

type UserState struct {
	admins          Set[int64]
	tgIdToXrayEmail map[int64]string

	usersPath  string
	adminsPath string
}

func (m UserState) IsAdmin(id int64) bool {
	_, ok := m.admins[id]
	return ok
}

func (m UserState) SetUser(id int64, email string) {
	m.tgIdToXrayEmail[id] = email
}

func (m UserState) Save() {
	os.WriteFile(m.adminsPath, []byte(m.serializeAdmins()), 0644)
	os.WriteFile(m.usersPath, []byte(m.serializeUsers()), 0644)
}

func (m UserState) serializeAdmins() string {
	admins := make([]string, 0, len(m.admins))
	for id := range m.admins {
		admins = append(admins, strconv.FormatInt(id, 10))
	}

	return strings.Join(admins, "\n")
}

func (m UserState) serializeUsers() string {
	users := make([]string, 0, len(m.tgIdToXrayEmail))
	for id, email := range m.tgIdToXrayEmail {
		users = append(users, strconv.FormatInt(id, 10)+":"+email)
	}

	return strings.Join(users, "\n")
}

func (m UserState) GetXrayEmail(id int64) (string, bool) {
	email, ok := m.tgIdToXrayEmail[id]
	return email, ok
}

func (m UserState) GetAllUsers() []string {
	return internal.Values(m.tgIdToXrayEmail)
}

func NewState(adminsPath, usersPath string) *UserState {
	adminsFile, err := os.ReadFile(adminsPath)
	admins := make(Set[int64])
	if err == nil {
		for _, line := range strings.Split(string(adminsFile), "\n") {
			adminId, _ := strconv.ParseInt(line, 10, 64)
			admins[adminId] = struct{}{}
		}
	}

	usersFile, err := os.ReadFile(usersPath)
	if err != nil {
		panic(err)
	}

	telegramIdToXrayEmail := make(map[int64]string)
	for _, line := range strings.Split(string(usersFile), "\n") {
		if len(line) == 0 {
			continue
		}

		idEmail := strings.Split(line, ":")

		id, _ := strconv.ParseInt(idEmail[0], 10, 64)
		telegramIdToXrayEmail[id] = idEmail[1]
	}

	return &UserState{
		admins:          admins,
		tgIdToXrayEmail: telegramIdToXrayEmail,

		usersPath:  usersPath,
		adminsPath: adminsPath,
	}
}
