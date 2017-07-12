package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceQuantumPassword() *schema.Resource {
	return &schema.Resource{
		Read:   resourceQuantumPasswordRead,
		Create: resourceQuantumPasswordCreate,
		Update: resourceQuantumPasswordCreate,
		Delete: resourceQuantumPasswordDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"length": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"lowercase": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"uppercase": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"numbers": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"specials": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			// "force_rotate": &schema.Schema{
			// 	Type:     schema.TypeBool,
			// 	Optional: true,
			// },
			"expires_in_days": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceQuantumPasswordRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("resourceQuantumPasswordRead - START")

	// Get parameters
	// name := d.Get("name").(string)
	// length := d.Get("length").(int)
	// lowercase := d.Get("lowercase").(int)
	// uppercase := d.Get("uppercase").(int)
	// numbers := d.Get("numbers").(int)
	// specials := d.Get("specials").(int)

	// expires_in_days := d.Get("expires_in_days").(string)
	// created_at := d.Get("created_at").(string)
	password := d.Get("password").(string)

	// Check last created_date and compare
	// if created_at - now in days > expired_in_days {
	// 	log.Printf("Generate a new password!")
	// 	password, _ := generatePassword(length, lowercase, uppercase, numbers, specials)
	// }

	// log.Printf("Password: %v\n", password)

	// password := "ABCE"
	// meta.key_to_save = "Yé"

	d.Set("password", password)
	d.SetId("-")

	log.Printf("resourceQuantumPasswordRead - END")

	return nil
}

func resourceQuantumPasswordCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("resourceQuantumPasswordCreate - START")

	// Get parameters
	// name := d.Get("name").(string)
	length := d.Get("length").(int)
	lowercase := d.Get("lowercase").(int)
	uppercase := d.Get("uppercase").(int)
	numbers := d.Get("numbers").(int)
	specials := d.Get("specials").(int)

	// expires_in_days := d.Get("expires_in_days").(string)
	// created_at := d.Get("created_at").(string)
	// password := d.Get("password").(string)

	// Check last created_date and compare
	// if created_at - now in days > expired_in_days {
	// 	log.Printf("Generate a new password!")
	password, _ := generatePassword(length, lowercase, uppercase, numbers, specials)
	// }

	log.Printf("Password: %v\n", password)

	// password := "ABCE"

	d.Set("password", password)
	d.SetId("-")

	log.Printf("resourceQuantumPasswordCreate - END")
	return nil
}

func resourceQuantumPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// To generate password
const lowercaseBytes = "abcdefghijklmnopqrstuvwxyz"
const uppercaseBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberBytes = "0123456789"
const specialBytes = "!%@#$"

func generatePassword(length int, lowercase int, uppercase int, numbers int, specials int) (string, error) {
	l := lowercase + uppercase + numbers + specials
	r := length - l

	password := ""

	if r < 0 {
		return "", errors.New("the password length cannot meet minimal requirement")
	}

	if r > 0 {
		rand.Seed(int64(time.Now().Nanosecond()))
		if lowercase > 0 {
			add := rand.Intn(r)
			lowercase += add
			r -= add
		}
		if r > 0 && uppercase > 0 {
			add := rand.Intn(r)
			uppercase += add
			r -= add
		}
		if r > 0 && numbers > 0 {
			add := rand.Intn(r)
			numbers += add
			r -= add
		}
		if r > 0 && specials > 0 {
			specials += r
		}
	}

	password += randStringBytes(lowercase, lowercaseBytes)
	password += randStringBytes(uppercase, uppercaseBytes)
	password += randStringBytes(numbers, numberBytes)
	password += randStringBytes(specials, specialBytes)

	password = shuffle(password)

	return password, nil
}

func randStringBytes(n int, chars string) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func shuffle(password string) string {
	arr := strings.Split(password, "")
	t := time.Now()
	rand.Seed(int64(t.Nanosecond())) // no shuffling without this line

	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i)
		arr[i], arr[j] = arr[j], arr[i]
	}

	return strings.Join(arr, "")
}

func isPasswordConform(password string, length int, lowercase int, uppercase int, numbers int, specials int) (bool, error) {

	if len(password) < length {
		return false, fmt.Errorf("password does not match minimum length requirement (%v < %v)", len(password), length)
	}

	re := regexp.MustCompile("[a-z]")
	l := len(re.FindAllString(password, -1))

	if l < lowercase {
		return false, fmt.Errorf("password does not match minimum requirement for lowercase length (%v < %v)", l, lowercase)
	}

	re = regexp.MustCompile("[A-Z]")
	u := len(re.FindAllString(password, -1))

	if u < uppercase {
		return false, fmt.Errorf("password does not match minimum requirement for uppercase length (%v < %v)", u, uppercase)
	}

	re = regexp.MustCompile("[0-9]")
	n := len(re.FindAllString(password, -1))

	if n < numbers {
		return false, fmt.Errorf("password does not match minimum requirement for numbers length (%v < %v)", n, numbers)
	}

	re = regexp.MustCompile(fmt.Sprintf("[%s]", specialBytes))
	s := len(re.FindAllString(password, -1))

	if s < specials {
		return false, fmt.Errorf("password does not match minimum requirement for special characters length (%v < %v)", s, specials)
	}

	return true, nil
}
