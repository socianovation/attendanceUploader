package models

type AttendanceClock struct {
	UserId    string `json:"user_id"`
	CheckTime string `json:"checktime"`
	CheckType string `json:"checktype"`
	SensorId  string `json:"sensor_id"`
}
