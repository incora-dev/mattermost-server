// Copyright (c) 2018-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const (
	GroupTypeLdap             = "ldap"
	GroupNameMaxLength        = 64
	GroupTypeMaxLength        = 64
	GroupDisplayNameMaxLength = 128
	GroupDescriptionMaxLength = 1024
	GroupRemoteIdMaxLength    = 2048
)

var groupTypesRequiringRemoteId = []string{
	GroupTypeLdap,
}

var groupTypes = []string{
	GroupTypeLdap,
}

type Group struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	RemoteId    string `json:"remote_id"`
	CreateAt    int64  `json:"create_at"`
	UpdateAt    int64  `json:"update_at"`
	DeleteAt    int64  `json:"delete_at"`
}

func (group *Group) IsValidForCreate() *AppError {
	if l := len(group.Name); l == 0 || l > GroupNameMaxLength {
		return NewAppError("Group.IsValidForCreate", "model.group.name.app_error", map[string]interface{}{"GroupNameMaxLength": GroupNameMaxLength}, "", http.StatusBadRequest)
	}

	if l := len(group.DisplayName); l == 0 || l > GroupDisplayNameMaxLength {
		return NewAppError("Group.IsValidForCreate", "model.group.display_name.app_error", map[string]interface{}{"GroupDisplayNameMaxLength": GroupDisplayNameMaxLength}, "", http.StatusBadRequest)
	}

	if len(group.Description) > GroupDescriptionMaxLength {
		return NewAppError("Group.IsValidForCreate", "model.group.description.app_error", map[string]interface{}{"GroupDescriptionMaxLength": GroupDescriptionMaxLength}, "", http.StatusBadRequest)
	}

	isValidType := false
	for _, groupType := range groupTypes {
		if group.Type == groupType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return NewAppError("Group.IsValidForCreate", "model.group.type.app_error", map[string]interface{}{"ValidGroupTypes": strings.Join(groupTypes, ", ")}, "", http.StatusBadRequest)
	}

	if len(group.RemoteId) > GroupRemoteIdMaxLength || len(group.RemoteId) == 0 && group.DoesRequireRemoteId() {
		return NewAppError("Group.IsValidForCreate", "model.group.remote_id.app_error", nil, "", http.StatusBadRequest)
	}

	return nil
}

func (group *Group) DoesRequireRemoteId() bool {
	for _, groupType := range groupTypesRequiringRemoteId {
		if groupType == group.Type {
			return true
		}
	}
	return false
}

func (group *Group) IsValidForUpdate() *AppError {
	if len(group.Id) != 26 {
		return NewAppError("Group.IsValidForUpdate", "model.group.id.app_error", nil, "", http.StatusBadRequest)
	}
	if group.CreateAt == 0 {
		return NewAppError("Group.IsValidForCreate", "model.group.create_at.app_error", nil, "", http.StatusBadRequest)
	}
	if group.UpdateAt == 0 {
		return NewAppError("Group.IsValidForCreate", "model.group.update_at.app_error", nil, "", http.StatusBadRequest)
	}
	if err := group.IsValidForCreate(); err != nil {
		return err
	}
	return nil
}

func (group *Group) ToJson() string {
	b, _ := json.Marshal(group)
	return string(b)
}

func GroupFromJson(data io.Reader) *Group {
	var group *Group
	json.NewDecoder(data).Decode(&group)
	return group
}

func GroupsToJson(groups []*Group) string {
	b, _ := json.Marshal(groups)
	return string(b)
}
