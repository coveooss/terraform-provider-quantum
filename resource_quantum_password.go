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

	log.Printf("resourceQuantumPasswordUpdate - END")
	return nil
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

	length := args.length
	lowercase := args.lowercase
	uppercase := args.uppercase
	numbers := args.numbers
	specials := args.specials

	// To avoid some chars, -1 can be passed for this group, skip them from count

	l := 0
	if lowercase > 0 {
		l += lowercase
	}
	if uppercase > 0 {
		l += uppercase
	}
	if numbers > 0 {
		l += numbers
	}
	if specials > 0 {
		l += specials
	}

	r := length - l

	password := ""

	if r < 0 || (lowercase+uppercase+numbers+specials) == -4 {
		return "", errors.New("the password length cannot meet minimal requirement")
	}

	rand.Seed(int64(time.Now().Nanosecond()))
	for r > 0 {
		if lowercase > -1 {
			add := rand.Intn(r)
			if add == 0 {
				add = 1
			}
			lowercase += add
			r -= add
		}
		if r > 0 && uppercase > -1 {
			add := rand.Intn(r)
			if add == 0 {
				add = 1
			}
			uppercase += add
			r -= add
		}
		if r > 0 && numbers > -1 {
			add := rand.Intn(r)
			if add == 0 {
				add = 1
			}
			numbers += add
			r -= add
		}
		if r > 0 && specials > -1 {
			add := rand.Intn(r)
			if add == 0 {
				add = 1
			}
			specials += add
			r -= add
		}
	}

	if lowercase > 0 {
		password += randStringBytes(lowercase, lowercaseBytes)
	}
	if uppercase > 0 {
		password += randStringBytes(uppercase, uppercaseBytes)
	}
	if numbers > 0 {
		password += randStringBytes(numbers, numberBytes)
	}
	if specials > 0 {
		password += randStringBytes(specials, specialBytes)
	}

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

	log.Printf("bleh pw: %v", len(args.password))
	log.Printf("bleh lt: %v", args.length)
	if len(args.password) < args.length {
		return false, fmt.Errorf("password does not match minimum length requirement (%v < %v)", len(args.password), args.length)
	}

	re := regexp.MustCompile("[a-z]")
	l := len(re.FindAllString(args.password, -1))

	if l < args.lowercase {
		return false, fmt.Errorf("password does not match minimum requirement for lowercase length (%v < %v)", l, args.lowercase)
	}

	re = regexp.MustCompile("[A-Z]")
	u := len(re.FindAllString(args.password, -1))

	if u < args.uppercase {
		return false, fmt.Errorf("password does not match minimum requirement for uppercase length (%v < %v)", u, args.uppercase)
	}

	re = regexp.MustCompile("[0-9]")
	n := len(re.FindAllString(args.password, -1))

	if n < args.numbers {
		return false, fmt.Errorf("password does not match minimum requirement for numbers length (%v < %v)", n, args.numbers)
	}

	re = regexp.MustCompile(fmt.Sprintf("[%s]", specialBytes))
	s := len(re.FindAllString(args.password, -1))

	if s < args.specials {
		return false, fmt.Errorf("password does not match minimum requirement for special characters length (%v < %v)", s, args.specials)
	}

	return true, nil
}

func getQuantumPasswordArgs(d *schema.ResourceData) *QuantumPasswordArgs {

	args := &QuantumPasswordArgs{
		name:          d.Get("name").(string),
		length:        d.Get("length").(int),
		lowercase:     d.Get("lowercase").(int),
		uppercase:     d.Get("uppercase").(int),
		numbers:       d.Get("numbers").(int),
		specials:      d.Get("specials").(int),
		expiresInDays: d.Get("expires_in_days").(int),
		password:      d.Get("password").(string),
		createdAt:     d.Get("created_at").(string),
	}

	// Setting some default for unspecified values
	if args.length == 0 {
		args.length = 20
	}

	// If specified 0, exclude some criterias, otherwise set default
	// if len(length) == 0 {
	// 	args.length = 20
	// }
	// if len(lowercase) == 0 {
	// 	args.lowercase = -1
	// } else {
	// 	args.lowercase = d.Get("lowercase").(int)
	// }
	// if len(uppercase) == 0 {
	// 	args.uppercase = -1
	// } else {
	// 	args.uppercase = d.Get("uppercase").(int)
	// }
	// if len(numbers) == 0 {
	// 	args.numbers = -1
	// } else {
	// 	args.numbers = d.Get("numbers").(int)
	// }
	// if len(specials) == 0 {
	// 	args.specials = -1
	// } else {
	// 	args.specials = d.Get("specials").(int)
	// }
	// if len(expiresInDays) == 0 {
	// 	args.expiresInDays = 30
	// } else {
	// 	args.expiresInDays = d.Get("expires_in_days").(int)
	// }

	log.Printf("QuantumArgs: %-v", args)

	return args

}

// QuantumPasswordArgs contains provided terraform arguments
type QuantumPasswordArgs struct {
	name          string
	length        int
	lowercase     int
	uppercase     int
	numbers       int
	specials      int
	expiresInDays int
	password      string
	createdAt     string
}
