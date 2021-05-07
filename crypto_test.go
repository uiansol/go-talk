package main

import "testing"

func TestCreateAndCheckPassword_1(t *testing.T) {
	t.Parallel()

	testCases := []string{
		"frodo",
		"gand@lf",
		"leg0las",
		"g1ml1",
		"_#aragorn!4",
	}

	for _, test := range testCases {
		hashedPassword, err := CreatePassword([]byte(test))
		if err != nil {
			t.Fatalf("Error while hashing the password: %s", err)
		}

		if err := CheckPassword(hashedPassword, []byte(test)); err != nil {
			t.Fatal("Wrong password creation")
		}
	}
}

// This test do the same thing as before.
// It was put here to demonstrate the parallelism.
func TestCreateAndCheckPassword_2(t *testing.T) {
	t.Parallel()

	testCases := []string{
		"frodo",
		"gand@lf",
		"leg0las",
		"g1ml1",
		"_#aragorn!4",
	}

	for _, test := range testCases {
		hashedPassword, err := CreatePassword([]byte(test))
		if err != nil {
			t.Fatalf("Error while hashing the password: %s", err)
		}

		if err := CheckPassword(hashedPassword, []byte(test)); err != nil {
			t.Fatal("Wrong password creation")
		}
	}
}
