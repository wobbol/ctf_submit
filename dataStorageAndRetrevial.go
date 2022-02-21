package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

//dataStorageConsistencyAndRetrevial.go

type Flag struct {
	Data string
	Desc string
	Pts  int
}
type User struct {
	Name  string
	Flags []Flag
	Pts   int
	//Cookie string
}
type Flagboard struct {
	Flags []Flag
}
type Scoreboard struct {
	Users []User
}
type CTF struct {
	Scoreboard Scoreboard
	Flagboard  Flagboard
}

// User
func (u *User) AddFlag(f Flag) {
	for _, flag := range u.Flags {
		if flag.Data == f.Data {
			return // flag list should be unique
		}
	}
	u.Pts += f.Pts
	u.Flags = append(u.Flags, f)
}

// Flagboard
func (b *Flagboard) AddFlag(f Flag) {
	for _, flag := range b.Flags {
		if flag.Data == f.Data {
			return // flag list should be unique
		}
	}
	b.Flags = append(b.Flags, f)
}
func (b *Flagboard) FindFlag(f string) (*Flag, error) {
	for i, flag := range b.Flags {
		if flag.Data == f {
			return &b.Flags[i], nil
		}
	}
	return nil, errors.New("Flag does not exist")
}

// Scoreboard
func (b *Scoreboard) AddUser(u User) {
	for _, user := range b.Users {
		if user.Name == u.Name {
			return // User list should be unique
		}
	}
	b.Users = append(b.Users, u)
}
func (b *Scoreboard) FindUser(n string) (*User, error) {
	for i, user := range b.Users {
		if user.Name == n {
			return &b.Users[i], nil
		}
	}
	return nil, errors.New("Username not registered")
}

// CTF
func NewCTF(dirPath string) *CTF {
	scorePath := dirPath + "/scoreboard.json"
	flagPath := dirPath + "/flags.json"
	var ctf CTF
	jsonToStruct(&ctf.Scoreboard, scorePath)
	jsonToStruct(&ctf.Flagboard, flagPath)

	err := ctf.validate()
	if err != nil {
		log.Println(err)
	}
	return &ctf
}
func (ctf CTF) SaveCTF(dirPath string) {
	structToJson(ctf.Scoreboard, dirPath+"/scoreboard.json")
	structToJson(ctf.Flagboard, dirPath+"/flags.json")
}
func (ctf *CTF) SubmitFlag(name string, data string) error {
	// Check for valid input
	log.Println(name)
	log.Println(data)
	if name == "" {
		return errors.New("Empty name not allowed.")
	}
	if data == "" {
		return errors.New("Empty flag not allowed.")
	}
	user, err := ctf.Scoreboard.FindUser(name)
	if err != nil {
		ctf.Scoreboard.AddUser(User{Name: name})
		//return err
		user, err = ctf.Scoreboard.FindUser(name)
		if err != nil {
			//eh give up
			return err
		}
	}
	flag, err := ctf.Flagboard.FindFlag(data)
	if err != nil {
		return err
	}
	// User has found flag. give it to them
	user.AddFlag(*flag)
	return nil
}
func (ctf CTF) validate() error {
	var err string
	for i, user := range ctf.Scoreboard.Users {
		if user.Name == "" {
			err = fmt.Sprintf("%sscoreboard user #%d: invalid name\n", err, i)
		}
		var acc int
		for _, flag := range user.Flags {
			acc += flag.Pts
		}
		if acc != user.Pts {
			log.Printf("user %s: Pts inconsistent with users captured flags\n", user.Name)
		}
		ctf.Scoreboard.Users[i].Pts = acc
	}
	for i, flag := range ctf.Flagboard.Flags {
		if flag.Data == "" {
			err = fmt.Sprintf("%sflaglist item #%d: invalid flag\n", err, i)
		}
	}
	if err != "" {
		return errors.New(err)
	}
	return nil
}

// Utilities
func structToJson(v interface{}, path string) {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(path, data, 0700)
	if err != nil {
		log.Fatal(err)
	}
}
func jsonToStruct(v interface{}, path string) {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(fileData, v)
	if err != nil {
		_, ok := err.(*json.InvalidUnmarshalError)
		if ok {
			log.Println("jsonToStruct(): need a pointer for v")
		}
		log.Fatal(err)
	}
}
