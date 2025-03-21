package main

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/bcrypt"
)

const minimumCharsPerCategory = 2
const specialChars = '!'

func resourceQuantumPassword() *schema.Resource {
	return &schema.Resource{
		Create: resourceQuantumPasswordCreate,
		Read:   func(d *schema.ResourceData, m interface{}) error { return update(d, false) },
		Update: func(d *schema.ResourceData, m interface{}) error { return update(d, true) },
		Delete: func(*schema.ResourceData, interface{}) error { return nil },

		Schema: map[string]*schema.Schema{
			"length": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"rotation": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"special_chars": &schema.Schema{
				Type:     schema.TypeString,
				Default:  categories[specialChars],
				Optional: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"previous_password": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"last_update": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"bcrypt": &schema.Schema{
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceQuantumPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	// Get parameters
	args := getQuantumPasswordArgs(d)

	password, genDate, bcrypt, err := generatePassword(args)

	if err == nil {
		d.Set("password", password)
		d.Set("last_update", genDate.Format(time.RFC3339))
		d.SetId(getMD5Hash(fmt.Sprintf("%s-%v", password, d.Get("last_update"))))
		d.Set("bcrypt", bcrypt)
	}

	return err
}

func update(d *schema.ResourceData, update bool) error {
	reset := d.HasChange("length") || d.HasChange("special_chars")

	// Get parameters
	args := getQuantumPasswordArgs(d)

	t, err := time.Parse(time.RFC3339, args.lastUpdate)
	if err != nil {
		log.Printf("Unable to parse the last generation date (%s), resetting password", args.lastUpdate)
		reset = true
	}

	if args.rotation != 0 && int(time.Now().Sub(t).Hours()/24) >= args.rotation {
		log.Printf("Generate a new password after %v days", args.rotation)
		reset = true
	}

	if reset {
		if !update {
			// If the reset is caused by the read operation, we keep the previous password
			// in order to be able to restore it if the reset should not have been done.
			// This could happen if according to the previous rotation period, the password
			// was expired, but was not really expired if we consider the new rotation (that
			// is only available during the update phase).
			d.Set("previous_password", d.Get("password"))
			d.Set("previous_bcrypt", d.Get("bcrypt"))
		}
		err = resourceQuantumPasswordCreate(d, nil)
	} else if update {
		previous := d.Get("previous_password")
		if previous != "" {
			// This was a false update, so we bring back the previous password
			d.Set("password", previous)
			d.Set("previous_password", "")
			d.Set("bcrypt", d.Get("previous_bcrypt"))
		}
	}

	return err
}

func generatePassword(args *QuantumPasswordArgs) (string, *time.Time, string, error) {
	charSets := make([]string, 0, len(categories))
	for category, charSet := range categories {
		if category == specialChars {
			if len(args.specialChars) > 0 {
				charSets = append(charSets, args.specialChars)
			}
		} else {
			charSets = append(charSets, charSet)
		}
	}

	if args.length < len(charSets) {
		return "", nil, "", fmt.Errorf("The password must be at least %d chars long", len(charSets))
	}

	var password string
	for i := 0; i < args.length; i++ {
		var group int
		if i < len(charSets)*minimumCharsPerCategory {
			// We take at least a minimum number of characters of each categories
			group = i % len(charSets)
		} else {
			// Afterwhile, we pick them randomly
			group = randInt(len(charSets))
		}
		chars := charSets[group]
		password += string(chars[randInt(len(chars))])
	}
	shuffled := shuffle(password)[:args.length]

	bcrypt, err := bcrypt.GenerateFromPassword([]byte(shuffled), 12)
	if err != nil {
		return "", nil, "", fmt.Errorf("Could not create hash %s", err)
	}

	generated := time.Now()
	return shuffled, &generated, string(bcrypt), nil
}

func shuffle(password string) string {
	arr := []byte(password)

	for i := 0; i < len(arr); i++ {
		j := randInt(len(arr))
		arr[i], arr[j] = arr[j], arr[i]
	}

	return string(arr)
}

func randInt(length int) int {
	i, _ := rand.Int(rand.Reader, big.NewInt(int64(length)))
	return int(i.Int64())
}

func getQuantumPasswordArgs(d *schema.ResourceData) *QuantumPasswordArgs {
	args := &QuantumPasswordArgs{
		length:       d.Get("length").(int),
		rotation:     d.Get("rotation").(int),
		lastUpdate:   d.Get("last_update").(string),
		specialChars: d.Get("special_chars").(string),
	}

	// Setting some default for unspecified values
	if args.length == 0 {
		args.length = 20
	}

	return args
}

// QuantumPasswordArgs contains provided terraform arguments
type QuantumPasswordArgs struct {
	length       int
	rotation     int
	lastUpdate   string
	specialChars string
}

var (
	baseSet    = map[rune]int{'a': 26, 'A': 26, '0': 10, specialChars: 15}
	categories = initializeCharSet()
)

func initializeCharSet() map[rune]string {
	categories := make(map[rune]string)

	for char, count := range baseSet {
		for i := 0; i < count; i++ {
			categories[char] += string(char + rune(i))
		}
	}
	return categories
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
