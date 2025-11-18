package models

type JwtServiceRole uint

const (
	JwtServiceRoleAdmin JwtServiceRole = iota
	JwtServiceRoleProductionLine
)
