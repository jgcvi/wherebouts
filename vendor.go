package wherebouts

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"strconv"
	//	"strings"
	"time"
)

type Vendor struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  []byte    `json:"-"`
	Salt      string    `json:"-"`
	Joined    time.Time `json:"joined"`
	Open      bool      `json:"open"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

func CreateVendor(name string, email string, pass string) *Vendor {
	f := new(Vendor)

	f.Name = name
	f.Email = email
	f.Joined = time.Now()
	f.Open = false

	f.generatePass(pass)

	return f
}

func (v *Vendor) isOpen() bool {
	return v.Open
}

func generateSalt() (string, error) {
	salt := make([]byte, 16)
	b, err := rand.Reader.Read(salt)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(b), nil
}

func (v *Vendor) generatePass(pass string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	v.Salt = salt

	v.Password = pbkdf2.Key([]byte(pass), []byte(v.Salt), 4096, sha256.Size, sha256.New)

	return "success", nil
}

func (v *Vendor) ValidatePass(pass string) bool {
	fmt.Print("in validatePass")
	hash := pbkdf2.Key([]byte(pass), []byte(v.Salt), 4096, sha256.Size, sha256.New)
	fmt.Printf("%v and %v\n", len(hash), len(v.Password))
	fmt.Printf("%b and \n%b\n", ([]byte(hash)), ([]byte(v.Password)))

	return bytes.Equal(hash, v.Password)
}

type VendorDatabase interface {
	ListVendors() ([]*Vendor, error)
	AddVendor(*Vendor) (int64, error)
	GetVendorByEmail(string) (*Vendor, error)
	GetVendorByID(int64) (*Vendor, error)
	GetNearbyVendors(float64, float64, float64) ([]*Vendor, error)
	ToggleAvailability(string, float64, float64)

	Close()
}
