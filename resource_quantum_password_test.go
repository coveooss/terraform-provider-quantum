package main

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"golang.org/x/crypto/bcrypt"
)

func testQuantumPasswordBasic(t *testing.T) {
	passwordRegex12 := regexp.MustCompile(`^.{12}$`)
	passwordRegex10 := regexp.MustCompile(`^.{10}$`)
	bcryptRegex := regexp.MustCompile(`^\$2[ayb]\$.{56}$`)

	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: testAccQuantumPasswordResource(12),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"quantum_password.test",
						"password",
						passwordRegex12,
					),
					resource.TestMatchResourceAttr(
						"quantum_password.test",
						"bcrypt",
						bcryptRegex,
					),
					testAccQuantumPasswordBcrypt("quantum_password.test"),
				),
			},
			// Force a rotation of password
			{
				Config: testAccQuantumPasswordResource(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"quantum_password.test",
						"password",
						passwordRegex10,
					),
					resource.TestMatchResourceAttr(
						"quantum_password.test",
						"bcrypt",
						bcryptRegex,
					),
					testAccQuantumPasswordBcrypt("quantum_password.test"),
				),
			},
		},
	})
}

func testAccQuantumPasswordResource(length int) string {
	return fmt.Sprintf(` 
	resource "quantum_password" "test" {
		special_chars = ""
		length        = %d
	}
	`, length)
}

func testAccQuantumPasswordBcrypt(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		err := bcrypt.CompareHashAndPassword([]byte(rs.Primary.Attributes["bcrypt"]), []byte(rs.Primary.Attributes["password"]))
		if err != nil {
			return fmt.Errorf("Bcrypt does not match: %s (%s)", resourceName, err)
		}

		return nil
	}
}
