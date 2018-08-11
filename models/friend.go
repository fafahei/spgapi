package models

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/utils"
)

var db *sql.DB

const (
	connection = 1 + iota
	subscription
	block
)

type Friends struct {
	Friends []string `json:"friends"`
}

type Email struct {
	Email string `json:"email"`
}

type Request struct {
	Requestor string `json:"requestor"`
	Target string `json:"target"`
}

type RequestUpdate struct {
	Sender string `json:"sender"`
	Text string `json:"text"`
}

type Status struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

type Response struct {
	Status Status `json:"status"`
	Friends []string `json:"friends"`
	Count int `json:"count"`
}

type ResponseUpdate struct {
	Status Status `json:"status"`
	Recipients []string `json:"recipients"`
}

type Info struct {
	Id int
	Email string
}

type Relationship struct {
	Id int
	Requestor *Info
	Target *Info
	Category int
}

func init() {
	host := "localhost"
	port := strconv.Itoa(3306)
	user := "root"
	password := "xxxxxxxxx"
	db_name := "spg"
	// connect to MySQL database
	db, _ = sql.Open("mysql", user+":"+password+"@tcp("+host+":"+port+")/"+db_name+"?charset=utf8")
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func Connect(f [2]*Info) (bool, error) {
	var status bool
	var r *Relationship
	for _, friend := range f {
		if err := checkFriend(friend); err != nil {
			return status, err
		}
	}
	r = &Relationship{0, f[0], f[1], connection}
	if err := checkRelationship(r); err != nil {
		return status, err
	}
	r.Category = block
	if err := checkRelationship(r); err != nil {
		//check block
		err = errors.New(r.Requestor.Email + " blocks " + r.Target.Email + ", so can't create friends connection!")
		return status, err
	}
	r.Requestor = f[1]
	r.Target = f[0]
	if err := checkRelationship(r); err != nil {
		//check block
		err = errors.New(r.Requestor.Email + " blocks " + r.Target.Email + ", so can't create friends connection!")
		return status, err
	}
	r.Category = connection
	//create connection
	sqlStr := "INSERT INTO friend_relationship(requestor, target, relationship) VALUES (?, ?, ?), (?, ?, ?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(strconv.Itoa(r.Requestor.Id),strconv.Itoa(r.Target.Id),strconv.Itoa(r.Category),strconv.Itoa(r.Target.Id),strconv.Itoa(r.Requestor.Id),strconv.Itoa(r.Category))
	if err != nil {
		log.Fatal(err)
	}
	status = true
	return status, nil
}

func GetFriends(friend *Info) (bool, []string, error) {
	var status bool
	var f []string
	if err := checkFriend(friend); err != nil {
		return status, f, err
	}
	stmt, err := db.Prepare("SELECT target FROM friend_relationship WHERE requestor = ? AND relationship = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(friend.Id, connection)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		_ = rows.Scan(&id)
		f = append(f, getEmail(id))
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	status = true
	if len(f) == 0 {
		err = errors.New("Can't find any friend for " + friend.Email + "!")
		return status, f, err
	}
	return status, f, nil
}

func GetCommonFriends(f [2]*Info) (bool, []string, error) {
	var status bool
	m := make(map[int][]string)
	var c []string
	var err error
	for i, friend := range f {
		status, m[i], err = GetFriends(friend)
		if status == false {
			return status, c, err
		} else if status == true && err != nil {
			return status, c, err
		}
	}
	c = intersect(m[0], m[1])
	status = true
	if len(c) == 0 {
		err = errors.New("No common friends!")
		return status, c, err
	}
	return status, c, nil
}

func Subscribe(r *Relationship) (bool, error){
	var status bool
	if err := checkFriend(r.Requestor); err != nil {
		return status, err
	}
	if err := checkFriend(r.Target); err != nil {
		return status, err
	}
	r.Category = block
	if err := checkRelationship(r); err != nil {
		//update subscription
		stmt, _ := db.Prepare("update friend_relationship set relationship = ? where id = ?")
		_, _ = stmt.Exec(subscription, r.Id)
		status = true
		return status, nil
	} else {
		r.Category = subscription
	}
	if err := checkRelationship(r); err != nil {
		return status, err
	}
	//create subscription
	sqlStr := "INSERT INTO friend_relationship(requestor, target, relationship) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(strconv.Itoa(r.Requestor.Id),strconv.Itoa(r.Target.Id),strconv.Itoa(r.Category))
	if err != nil {
		log.Fatal(err)
	}
	status = true
	return status, nil
}

func Block(r *Relationship) (bool, error) {
	var status bool
	if err := checkFriend(r.Requestor); err != nil {
		return status, err
	}
	if err := checkFriend(r.Target); err != nil {
		return status, err
	}
	r.Category = subscription
	if err := checkRelationship(r); err != nil {
		//update block
		stmt, _ := db.Prepare("update friend_relationship set relationship = ? where id = ?")
		_, _ = stmt.Exec(block, r.Id)
		status = true
		return status, nil
	} else {
		r.Category = block
	}
	if err := checkRelationship(r); err != nil {
		return status, err
	}
	//create block
	sqlStr := "INSERT INTO friend_relationship(requestor, target, relationship) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(strconv.Itoa(r.Requestor.Id),strconv.Itoa(r.Target.Id),strconv.Itoa(r.Category))
	if err != nil {
		log.Fatal(err)
	}
	status = true
	return status, nil
}

func GetUpdateFriends(target *Info, r []*Info) (bool, []string, error) {
	var status bool
	var f []string
	var b []string
	if err := checkFriend(target); err != nil {
		return status, f, err
	}
	for _, friend := range r {
		if err := checkFriend(friend); err == nil {
			if !utils.InSlice(friend.Email, f) {
				f = append(f, friend.Email)
			}
		}
	}
	//get connected or subscribed requestors
	stmt, err := db.Prepare("SELECT requestor FROM friend_relationship WHERE target = ? AND relationship IN (?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(target.Id, connection, subscription)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		_ = rows.Scan(&id)
		email := getEmail(id)
		if !utils.InSlice(email, f) {
			f = append(f, email)
		}
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	//get blocked requestors
	stmt, err = db.Prepare("SELECT requestor FROM friend_relationship WHERE target = ? AND relationship = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err = stmt.Query(target.Id, block)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		_ = rows.Scan(&id)
		b = append(b, getEmail(id))
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	f = diff(f, b)
	status = true
	return status, f, nil
}

//check friend
func checkFriend(f *Info) error {
	stmt, err := db.Prepare("SELECT id FROM friend_info WHERE email = ?")
	if err != nil {
		log.Fatal(err)
	}
	err = stmt.QueryRow(f.Email).Scan(&f.Id)
	if err != nil {
		err = errors.New("This friend " + f.Email + " is non-exist!")
		return err
	}
	return nil
}

//get email of friend
func getEmail(id int) string {
	var email string
	stmt, err := db.Prepare("SELECT email FROM friend_info WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	err = stmt.QueryRow(id).Scan(&email)
	if err != nil {
		return email
	}
	return email
}

//check friend relationship
func checkRelationship(r *Relationship) error {
	m := map[int]string{1: "connection", 2: "subscription", 3: "block"}
	stmt, err := db.Prepare("SELECT id FROM friend_relationship WHERE requestor = ? AND target = ? AND relationship = ?")
	if err != nil {
		log.Fatal(err)
	}
	err = stmt.QueryRow(r.Requestor.Id, r.Target.Id, r.Category).Scan(&r.Id)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	} else if err == nil {
		err = errors.New("This " + m[r.Category] + " relationship has been created!")
		return err
	}
	return nil
}

func diff(slice1, slice2 []string) (diff []string) {
	for _, v := range slice1 {
		if !utils.InSlice(v, slice2) {
			diff = append(diff, v)
		}
	}
	return
}

func intersect(slice1, slice2 []string) (intersect []string) {
	for _, v := range slice1 {
		if utils.InSlice(v, slice2) {
			intersect = append(intersect, v)
		}
	}
	return
}

/***

//get intersection of two arrays
func intersect(a, b []string) []string {
	m := make(map[string]bool)
	var c []string
	for _, item := range a {
		m[item] = true
	}
	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return c
}

//check element
func contain(slice []string, search string) bool {
	for _, value := range slice {
		if value == search {
			return true
		}
	}
	return false
}

***/