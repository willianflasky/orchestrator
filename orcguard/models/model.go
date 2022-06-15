package models

type Info struct {
	Id     int    `db:"id"`
	Domain string `db:"domain"`
	Ip     string `db:"ip"`
	Rw     string `db:"rw"`
	Ts     string `db:"ts"`
}
