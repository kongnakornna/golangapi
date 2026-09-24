package usecase

import (
	"icmongolang/pkg/helpers"
)

type AlarmUseCase interface {
	ValidateAlarm(dto helpers.AlarmDetailDto) helpers.AlarmDetailResult
	ValidateAlarmEn(dto helpers.AlarmDetailDto) helpers.AlarmDetailResult
	ValidateAlarmTh(dto helpers.AlarmDetailDto) helpers.AlarmDetailResult
}

type alarmUseCase struct{}

func NewAlarmUseCase() AlarmUseCase {
	return &alarmUseCase{}
}

func (u *alarmUseCase) ValidateAlarm(dto helpers.AlarmDetailDto) helpers.AlarmDetailResult {
	return helpers.AlarmDetailValidate(dto)
}
func (u *alarmUseCase) ValidateAlarmEn(dto helpers.AlarmDetailDto) helpers.AlarmDetailResult {
	return helpers.AlarmDetailValidateEn(dto)
}
func (u *alarmUseCase) ValidateAlarmTh(dto helpers.AlarmDetailDto) helpers.AlarmDetailResult {
	return helpers.AlarmDetailValidateTh(dto)
}
