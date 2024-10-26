package wgg

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// InitWireGuardKeys initializes a new WireGuard key pair for both the node and the
// client, and writes the private and public keys to the respective files in the
// given directory. The files are created with 0600 permissions.
//
// If any error occurs during key generation or file I/O, the function returns an
// error. If the "wg" command is not available on the system, the function returns
// an error.
func InitAllKeys(outDir string) (string, string, string, string, error) {
	if !IsCommandAvailable("wg") {
		return "", "", "", "", errors.New("the wg command is not available on the system, please install it")
	}

	err := InitWireGuardKeys(
		outDir+"/node.private.key",
		outDir+"/node.public.key",
	)
	if err != nil {
		return "", "", "", "", errors.New("error while creating node keys: " + err.Error())
	}

	privateNodeKey, publicNodeKey, err := LoadWireGuardKey(
		outDir+"/node.private.key",
		outDir+"/node.public.key",
	)
	if err != nil {
		return "", "", "", "", errors.New("error while loading node keys: " + err.Error())
	}

	err = InitWireGuardKeys(
		outDir+"/client.private.key",
		outDir+"/client.public.key",
	)
	if err != nil {
		return "", "", "", "", errors.New("error while creating client keys: " + err.Error())
	}

	privateClientKey, publicClientKey, err := LoadWireGuardKey(
		outDir+"/client.private.key",
		outDir+"/client.public.key",
	)
	if err != nil {
		return "", "", "", "", errors.New("error while loading client keys: " + err.Error())
	}

	return privateNodeKey, publicNodeKey, privateClientKey, publicClientKey, nil
}

// GenerateWireGuardKeys generates a new WireGuard key pair.
// This function runs the "wg genkey" command to create a private key
// and then pipes it into the "wg pubkey" command to derive the corresponding
// public key. It returns the private key, public key, and an error if any
// command execution fails.
func GenerateWireGuardKeys() (string, string, error) {
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

// InitWireGuardKeys initializes a new WireGuard key pair and writes it to the
// given paths.
//
// If the private key and public key files already exist, the function does
// nothing and returns nil. Otherwise, it generates a new key pair and writes
// the private key and public key to the respective files, making sure that only
// the owner can read them.
//
// If any command execution or file I/O fails, the function returns an error.
func InitWireGuardKeys(privateKeyPath, publicKeyPath string) error {
	var keyIsMissing bool = false
	var err error
	_, err = os.Stat(privateKeyPath)
	if os.IsNotExist(err) {
		keyIsMissing = true
	}

	if !keyIsMissing {
		_, err = os.Stat(publicKeyPath)
		if os.IsNotExist(err) {
			keyIsMissing = true
		}
	}

	if !keyIsMissing {
		return nil
	}

	privateKey, publicKey, err := GenerateWireGuardKeys()
	if err != nil {
		return fmt.Errorf("error generating keys: %w", err)
	}

	err = os.WriteFile(privateKeyPath, []byte(privateKey), 0600)
	if err != nil {
		return fmt.Errorf("error writing private key to file: %w", err)
	}

	err = os.WriteFile(publicKeyPath, []byte(publicKey), 0600)
	if err != nil {
		return fmt.Errorf("error writing public key to file: %w", err)
	}

	return nil
}

func LoadWireGuardKey(privateKeyPath, publicKeyPath string) (string, string, error) {
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
