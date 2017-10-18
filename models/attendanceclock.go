package models

type AttendanceClock struct {
	UserId    string `json:"badgeno"`
	Checktime string `json:"checktime"`
	CompanyId string `json:"company_id"`
}
