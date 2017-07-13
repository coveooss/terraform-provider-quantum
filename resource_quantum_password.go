package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceQuantumPassword() *schema.Resource {
	return &schema.Resource{
		Read:   resourceQuantumPasswordRead,
		Create: resourceQuantumPasswordCreate,
		Update: resourceQuantumPasswordUpdate,
		Delete: resourceQuantumPasswordDelete,
		Exists: resourceQuantumPasswordExists,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"length": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"expires_in_days": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceQuantumPasswordRead(d *schema.ResourceData, meta interface{}) error {

	log.Printf("resourceQuantumPasswordRead - START")

	reset := false

	// Get parameters
	args := getQuantumPasswordArgs(d)

	log.Printf("Current length inside: %v", args.length)
	log.Printf("Current name inside: %v", args.name)
	log.Printf("Current password inside: %v", args.password)

	// Check if the password is conform
	valid, err := isPasswordConform(args)

	if !valid {
		log.Printf("Password will be reset: %v", err)
		reset = true
	} else {
		// Check last created_date and compare
		t, _ := time.Parse("2006-01-02", args.createdAt)

		diff := time.Now().Sub(t)
		days := int(diff.Hours() / 24)

		log.Printf("Diff Days: %v", days)

		if days >= args.expiresInDays {
			log.Printf("Generate a new password after %v days!", args.expiresInDays)
			reset = true
		}
	}

	if reset {
		password, _ := generatePassword(args)
		d.Set("password", password)
		d.Set("created_at", time.Now().Format("2006-01-02"))
	}

	d.SetId("-")

	log.Printf("resourceQuantumPasswordRead - END")

	return nil
}

func resourceQuantumPasswordCreate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("resourceQuantumPasswordCreate - START")

	// Get parameters
	args := getQuantumPasswordArgs(d)

	log.Printf("Generate a new password!")
	password, _ := generatePassword(args)

	log.Printf("Password: %v", password)

	d.Set("created_at", time.Now().Format("2006-01-02"))
	d.Set("password", password)
	d.SetId("-")

	log.Printf("resourceQuantumPasswordCreate - END")
	return nil
}

func resourceQuantumPasswordUpdate(d *schema.ResourceData, meta interface{}) error {

	log.Printf("resourceQuantumPasswordUpdate - START")

	// Normally, the Read will have set the new password
	// but when attributes are updated, the check is done
	// on previous attributes instead of new ones, so
	// recall it with new one to get password updated
	// with latest values.
	err := resourceQuantumPasswordRead(d, meta)

	log.Printf("resourceQuantumPasswordUpdate - END")
	return err
}

func resourceQuantumPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceQuantumPasswordExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	// Try to create a password with provided info
	// Need to crash early for invalid resource parameters

	log.Printf("resourceQuantumPasswordExists - START")

	pw, err := generatePassword(getQuantumPasswordArgs(d))

	log.Printf("resourceQuantumPasswordExists - END")

	return len(pw) > 0, err

}

// To generate password
const lowercaseBytes = "abcdefghijklmnopqrstuvwxyz"
const uppercaseBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberBytes = "0123456789"
const specialBytes = "!%@#$"

func generatePassword(args *QuantumPasswordArgs) (string, error) {

	log.Printf("generatePassword - START - %-v", args)

	if args.length < 4 {
		return "", errors.New("the password must be at least 4 chars long")
	}

	assign := args.length / 4

	password := ""

	password += randStringBytes(assign, lowercaseBytes)
	password += randStringBytes(assign, uppercaseBytes)
	password += randStringBytes(assign, numberBytes)
	password += randStringBytes(assign+args.length%4, specialBytes)

	password = shuffle(password)

	log.Printf("generatePassword - END - %v", password)

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

func isPasswordConform(args *QuantumPasswordArgs) (bool, error) {

	log.Printf("isPasswordConform pw: %v", len(args.password))
	log.Printf("isPasswordConform lt: %v", args.length)
	if len(args.password) != args.length {
		return false, fmt.Errorf("password does not match length requirement (%v != %v)", len(args.password), args.length)
	}

	return true, nil
}

func getQuantumPasswordArgs(d *schema.ResourceData) *QuantumPasswordArgs {

	args := &QuantumPasswordArgs{
		name:          d.Get("name").(string),
		length:        d.Get("length").(int),
		expiresInDays: d.Get("expires_in_days").(int),
		password:      d.Get("password").(string),
		createdAt:     d.Get("created_at").(string),
	}

	// Setting some default for unspecified values
	if args.length == 0 {
		args.length = 20
	}

	log.Printf("QuantumArgs: %v", args)

	return args

}

// QuantumPasswordArgs contains provided terraform arguments
type QuantumPasswordArgs struct {
	name          string
	length        int
	expiresInDays int
	password      string
	createdAt     string
}
