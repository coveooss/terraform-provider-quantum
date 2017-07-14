package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

const minimumCharsPerCategory = 2

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
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"previous_password": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_update": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceQuantumPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	// Get parameters
	args := getQuantumPasswordArgs(d)

	password, genDate, err := generatePassword(args)

	if err == nil {
		d.Set("password", password)
		d.Set("last_update", genDate.Format(time.RFC3339))
		d.SetId(getMD5Hash(fmt.Sprintf("%s-%v", password, d.Get("last_update"))))
	}

	return err
}

func update(d *schema.ResourceData, update bool) error {
	reset := d.HasChange("length")

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
		}
		err = resourceQuantumPasswordCreate(d, nil)
	} else if update {
		previous := d.Get("previous_password")
		if previous != "" {
			// This was a false update, so we bring back the previous password
			d.Set("password", previous)
			d.Set("previous_password", "")
		}
	}

	return err
}

func generatePassword(args *QuantumPasswordArgs) (string, *time.Time, error) {
	rand.Seed(int64(time.Now().Nanosecond()))

	if args.length < len(categories) {
		return "", nil, fmt.Errorf("The password must be at least %d chars long", len(categories))
	}

	var password string
	for i := 0; i < args.length; i++ {
		var group int
		if i < len(categories)*minimumCharsPerCategory {
			// We take at least a minimum number of characters of each categories
			group = i % len(categories)
		} else {
			// Afterwhile, we pick them randomly
			group = rand.Intn(len(categories))
		}
		chars := categories[group]
		password += string(chars[rand.Intn(len(chars))])
	}

	generated := time.Now()
	return shuffle(password)[:args.length], &generated, nil
}

func shuffle(password string) string {
	rand.Seed(int64(time.Now().Nanosecond()))

	arr := []byte(password)

	for i := 0; i < len(arr); i++ {
		j := rand.Intn(len(arr))
		arr[i], arr[j] = arr[j], arr[i]
	}

	return string(arr)
}

func getQuantumPasswordArgs(d *schema.ResourceData) *QuantumPasswordArgs {
	args := &QuantumPasswordArgs{
		length:     d.Get("length").(int),
		rotation:   d.Get("rotation").(int),
		lastUpdate: d.Get("last_update").(string),
	}

	// Setting some default for unspecified values
	if args.length == 0 {
		args.length = 20
	}

	return args
}

// QuantumPasswordArgs contains provided terraform arguments
type QuantumPasswordArgs struct {
	length     int
	rotation   int
	lastUpdate string
}

var (
	baseSet    = map[rune]int{'a': 26, 'A': 26, '0': 10, '!': 15}
	categories = initializeCharSet()
)

func initializeCharSet() []string {
	categories := make([]string, len(baseSet))
	categoryCount := 0
	for char, count := range baseSet {
		for i := 0; i < count; i++ {
			categories[categoryCount] += string(char + rune(i))
		}
		categoryCount++
	}
	return categories
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
