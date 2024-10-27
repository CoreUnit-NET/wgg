package wgg

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GenerateWireGuardKeyPair generates a new WireGuard key pair.
//
// This function runs the "wg genkey" command to create a private key
// and then pipes it into the "wg pubkey" command to derive the corresponding
// public key. It returns the private key, public key, and an error if any
// command execution fails.
func GenerateWireGuardKeyPair() (string, string, error) {
	cmdGenKey := exec.Command("wg", "genkey")

	privateKeyBuf := &bytes.Buffer{}
	cmdGenKey.Stdout = privateKeyBuf

	if err := cmdGenKey.Run(); err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	privateKeyString := privateKeyBuf.String()

	cmdPubKey := exec.Command("wg", "pubkey")
	cmdPubKey.Stdin = bytes.NewReader([]byte(privateKeyString))

	publicKeyBuf := &bytes.Buffer{}
	cmdPubKey.Stdout = publicKeyBuf

	if err := cmdPubKey.Run(); err != nil {
		return "", "", fmt.Errorf("failed to generate public key: %w", err)
	}

	privateKeyString = strings.TrimSpace(privateKeyString)
	publicKeyString := strings.TrimSpace(publicKeyBuf.String())

	return privateKeyString, publicKeyString, nil
}

// LoadWireGuardKeyPair loads a WireGuard key pair from the given paths.
//
// It reads the private and public key files and returns the contents of both
// as strings. If any file I/O fails, it returns an error. If the file contents
// are empty, it returns an error.
func LoadWireGuardKeyPair(
	privateKeyPath string,
	publicKeyPath string,
) (string, string, error) {
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("error reading private key: %w", err)
	} else if len(privateKey) == 0 {
		return "", "", errors.New("private key is empty")
	}

	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("error reading public key: %w", err)
	} else if len(publicKey) == 0 {
		return "", "", errors.New("public key is empty")
	}

	return string(privateKey), string(publicKey), nil
}

// SaveWireGuardKeyPair saves the given WireGuard key pair to the specified file paths.
//
// It writes the private key to the provided privateKeyPath, and the public key
// to the provided publicKeyPath, ensuring that only the owner has read/write
// permissions on the files.
//
// If writing any of the keys to their respective files fails, the function
// returns an error.
func SaveWireGuardKeyPair(
	privateKeyPath string,
	publicKeyPath string,
	privateKey string,
	publicKey string,
) error {
	err := os.WriteFile(privateKeyPath, []byte(privateKey), 0600)
	if err != nil {
		return fmt.Errorf("error writing private key: %w", err)
	}

	err = os.WriteFile(publicKeyPath, []byte(publicKey), 0600)
	if err != nil {
		return fmt.Errorf("error writing public key: %w", err)
	}

	return nil
}

// InitWireGuardKeyPair initializes a new WireGuard key pair if the given
// files do not exist yet.
//
// If the private key and public key files already exist, the function does
// nothing and returns the contents of both files as strings. Otherwise, it
// generates a new key pair and writes the private key and public key to the
// respective files, making sure that only the owner can read them.
//
// If any command execution or file I/O fails, the function returns an error.
func InitWireGuardKeyPair(
	privateKeyPath string,
	publicKeyPath string,
) (string, string, error) {
	privateKey, publicKey, err := LoadWireGuardKeyPair(privateKeyPath, publicKeyPath)
	if err != nil {
		// gen keys and save them
		privateKey, publicKey, err = GenerateWireGuardKeyPair()
		if err != nil {
			return "", "", err
		}

		err = SaveWireGuardKeyPair(privateKeyPath, publicKeyPath, privateKey, publicKey)
		if err != nil {
			return "", "", err
		}
	}

	return privateKey, publicKey, nil

}
