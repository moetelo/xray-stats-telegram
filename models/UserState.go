package models

import (
	"os"
	"strconv"
	"strings"
	"xray-stats-telegram/internal"
)

type UserState struct {
	admins          *Set[int64]
	tgIdToXrayEmail *map[int64]string
}

func (m UserState) IsAdmin(id int64) bool {
	_, ok := (*m.admins)[id]
	return ok
}

func (m UserState) GetXrayEmail(id int64) (string, bool) {
	email, ok := (*m.tgIdToXrayEmail)[id]
	return email, ok
}

func (m UserState) GetAllUsers() *[]string {
	return internal.Values(m.tgIdToXrayEmail)
}

func NewState(adminsSlice *[]int64, telegramIdToXrayEmail *map[int64]string) *UserState {
	admins := make(Set[int64])
	for _, admin := range *adminsSlice {
		admins[admin] = struct{}{}
	}

	return &UserState{
		admins:          &admins,
		tgIdToXrayEmail: telegramIdToXrayEmail,
	}
}

func NewStateFromConfigs(adminsPath, usersPath string) *UserState {
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
		admins:          &admins,
		tgIdToXrayEmail: &telegramIdToXrayEmail,
	}
}
