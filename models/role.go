package models

type Role string

const (
	RoleUser  Role = "USER"
	RoleAgent Role = "AGENT"
	RoleAdmin Role = "ADMIN"
)
