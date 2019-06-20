package db

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/satori/go.uuid"
	"os"
	"strings"
)

type User struct {
	UserId   int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Username string
	Password string
	Token    string
	Role     string
	Status   string
}

type UserWithAuth struct {
	User
	SubDomain string
	Port   int64
}

type SubDomainAuth struct {
	SubDomain string `gorm:"primary_key"`
	UserId int64
}

type PortAuth struct {
	Port   int64 `gorm:"primary_key"`
	UserId int64
}

type ApplyAuth struct {
	ApplyId         int64  `gorm:"primary_key;AUTO_INCREMENT"`
	ApplySubDomains string
	ApplyPorts      string
	Status          string
	PassSubDomains  string
	PassPorts       string
}

func InitTables() {
	conn, err := GetConnection()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.AutoMigrate(&User{}).AddUniqueIndex("uniq_username", "username")
	conn.AutoMigrate(&SubDomainAuth{})
	conn.AutoMigrate(&PortAuth{})

	var userCount int64
	err = conn.Model(&User{}).Count(&userCount).Error
	if err != nil {
		panic(err)
	}

	if userCount == 0 {
		// 生成管理员用户
		fmt.Println("creating an admin user...")
		var username, password, token string

		fmt.Print("admin username: ")
		_, _ = fmt.Scanln(&username)
		if username == "" {
			username = "admin"
			fmt.Println("use default username admin")
		}

		fmt.Print("admin password: ")
		_, _ = fmt.Scanln(&password)
		if password == "" {
			password = strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)[:16]
			fmt.Println("use random password " + password)
		}
		hash := sha1.Sum([]byte(password))

		token = strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
		fmt.Println("use random token " + token)

		err = conn.Create(&User{
			Username: username,
			Password: hex.EncodeToString(hash[:]),
			Token: token,
			Role: "admin",
		}).Error
		if err != nil {
			panic(err)
		}

		f, err := os.OpenFile("./password.txt", os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0644)
		if err == nil {
			defer f.Close()
			_, _ = f.WriteString(fmt.Sprintf("username: %s\npassword: %s\ntoken: %s\n", username, password, token))
		}
	}
}

func CreateUser(user *User) error {
	conn, err := GetConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	if user.Token == "" {
		user.Token = strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
	}
	return conn.Create(user).Error
}

func GetUserById(id int64) (*User, error) {
	conn, err := GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var user User
	return &user, conn.First(&user, "user_id = ?", id).Error
}

func GetUserByToken(token string) (*User, error) {
	conn, err := GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var user User
	return &user, conn.First(&user, "token = ?", token).Error
}

func LoginGetUser(username string, password string) (*User, error) {
	conn, err := GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var user User
	return &user, conn.First(&user, "username = ? AND password = ?", username, password).Error
}

func ListUsers() (users []User, err error) {
	conn, err := GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return users, conn.Select("user_id, username, role, status").Find(&users).Error
}

func AuthSubDomainToUser(subDomain string, userId int64) error {
	conn, err := GetConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Create(&SubDomainAuth{
		SubDomain: subDomain,
		UserId: userId,
	}).Error
}

func AuthPortToUser(port int64, userId int64) error {
	conn, err := GetConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Create(&PortAuth{
		Port: port,
		UserId: userId,
	}).Error
}

func GetAllSubDomainAuth(userId int64) (subDomains []UserWithAuth, err error) {
	conn, err := GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var user User
	if err := conn.First(&user, "user_id = ?", userId).Error; err != nil {
		return nil, err
	}
	if user.Role == "admin" {
		return subDomains, conn.Raw("SELECT sub_domain_auths.user_id, username, sub_domain" +
			" FROM sub_domain_auths LEFT JOIN users ON users.user_id = sub_domain_auths.user_id").
			Scan(&subDomains).Error
	}
	return subDomains, conn.Raw("SELECT sub_domain_auths.user_id, username, sub_domain" +
		" FROM sub_domain_auths LEFT JOIN users ON users.user_id = sub_domain_auths.user_id WHERE sub_domain_auths.user_id = ?", userId).
		Scan(&subDomains).Error
}

func GetAllPortAuth(userId int64) (ports []UserWithAuth, err error) {
	conn, err := GetConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var user User
	if err := conn.First(&user, "user_id = ?", userId).Error; err != nil {
		return nil, err
	}
	if user.Role == "admin" {
		return ports, conn.Raw("SELECT port_auths.user_id, username, port" +
			" FROM port_auths LEFT JOIN users ON users.user_id = port_auths.user_id").
			Scan(&ports).Error
	}
	return ports, conn.Raw("SELECT port_auths.user_id, username, port" +
		" FROM port_auths LEFT JOIN users ON users.user_id = port_auths.user_id WHERE port_auths.user_id = ?", userId).
		Scan(&ports).Error
}

func CanUserUsePort(userId int64, port int64) bool {
	conn, err := GetConnection()
	if err != nil {
		return false
	}
	defer conn.Close()
	var portAuth PortAuth
	return conn.First(&portAuth, "user_id = ? AND port = ?", userId, port).Error == nil
}

func CanUserUseSubDomain(userId int64, subDomain string) bool {
	conn, err := GetConnection()
	if err != nil {
		return false
	}
	defer conn.Close()
	var subDomainAuth SubDomainAuth
	return conn.First(&subDomainAuth, "user_id = ? AND sub_domain = ?", userId, subDomain).Error == nil
}
