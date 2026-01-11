package model

import (
	"database/sql"
	"time"

	"github.com/konsultin/project-goes-here/dto"
)

type AuthSession struct {
	BaseField
	Id                    int64                        `db:"id" json:"id"`
	Xid                   string                       `db:"xid" json:"xid"`
	SubjectId             string                       `db:"subject_id" json:"subjectId"`
	SubjectTypeId         dto.Role_Enum                `db:"subject_type_id" json:"subjectTypeId"`
	AuthProviderId        dto.AuthProvider_Enum        `db:"auth_provider_id" json:"authProviderId"`
	DevicePlatformId      dto.DevicePlatform_Enum      `db:"device_platform_id" json:"devicePlatformId"`
	DeviceId              string                       `db:"device_id" json:"deviceId"`
	Device                *AuthSessionDevice           `db:"device" json:"device"`
	NotificationChannelId dto.NotificationChannel_Enum `db:"notification_channel_id" json:"notificationChannelId"`
	NotificationToken     sql.NullString               `db:"notification_token" json:"notificationToken"`
	ExpiredAt             time.Time                    `db:"expired_at" json:"expiredAt"`
	StatusId              dto.ControlStatus_Enum       `db:"status_id" json:"statusId"`
}

type AuthSessionDevice struct {
	DeviceId         string                  `json:"deviceId"`
	DevicePlatformId dto.DevicePlatform_Enum `json:"devicePlatformId"`
	ClientIp         string                  `json:"clientIp"`
}
