package main

import "database/sql"

type User struct {
	Id			string
	Username	string
	Password	string
}

type Post struct {
	Id				string		
	Author			string
	Time_stamp		string
	Title			string
	Context			string
}

type DB struct {
	db			 *sql.DB
}