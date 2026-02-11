package crypto

import (
	"encoding/base64"
	"os"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	// Setup: Generate a test key
	testKey := make([]byte, 32)
	for i := range testKey {
		testKey[i] = byte(i)
	}
	os.Setenv("ENCRYPTION_KEY", base64.StdEncoding.EncodeToString(testKey))
	defer os.Unsetenv("ENCRYPTION_KEY")

	// Initialize
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	tests := []struct {
		name      string
		plaintext string
	}{
		{
			name:      "short text",
			plaintext: "hello",
		},
		{
			name:      "oauth token",
			plaintext: "gho_16C7e42F292c6912E7710c838347Ae178B4a",
		},
		{
			name:      "long text",
			plaintext: "This is a much longer text that simulates a GitHub OAuth token with lots of characters and special symbols!@#$%^&*()",
		},
		{
			name:      "empty string",
			plaintext: "",
		},
		{
			name:      "unicode",
			plaintext: "Hello 世界 🌍",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			// Verify it's different from plaintext
			if encrypted == tt.plaintext && tt.plaintext != "" {
				t.Error("Encrypted text should be different from plaintext")
			}

			// Verify it's base64 encoded
			_, err = base64.StdEncoding.DecodeString(encrypted)
			if err != nil {
				t.Errorf("Encrypted text should be valid base64: %v", err)
			}

			// Decrypt
			decrypted, err := Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			// Verify it matches original
			if decrypted != tt.plaintext {
				t.Errorf("Decrypted text doesn't match original. Got %q, want %q", decrypted, tt.plaintext)
			}
		})
	}
}

func TestEncryptDifferentNonces(t *testing.T) {
	// Setup
	testKey := make([]byte, 32)
	for i := range testKey {
		testKey[i] = byte(i)
	}
	os.Setenv("ENCRYPTION_KEY", base64.StdEncoding.EncodeToString(testKey))
	defer os.Unsetenv("ENCRYPTION_KEY")

	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	plaintext := "test"

	// Encrypt same text twice
	encrypted1, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("First encrypt failed: %v", err)
	}

	encrypted2, err := Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Second encrypt failed: %v", err)
	}

	// Should produce different ciphertexts due to different nonces
	if encrypted1 == encrypted2 {
		t.Error("Encrypting same text twice should produce different ciphertexts")
	}

	// But both should decrypt to same plaintext
	decrypted1, err := Decrypt(encrypted1)
	if err != nil {
		t.Fatalf("First decrypt failed: %v", err)
	}

	decrypted2, err := Decrypt(encrypted2)
	if err != nil {
		t.Fatalf("Second decrypt failed: %v", err)
	}

	if decrypted1 != plaintext || decrypted2 != plaintext {
		t.Error("Both decrypted texts should match original plaintext")
	}
}

func TestInitErrors(t *testing.T) {
	tests := []struct {
		name    string
		keyEnv  string
		wantErr bool
	}{
		{
			name:    "missing key",
			keyEnv:  "",
			wantErr: true,
		},
		{
			name:    "invalid base64",
			keyEnv:  "not-valid-base64!!!",
			wantErr: true,
		},
		{
			name:    "wrong key length",
			keyEnv:  base64.StdEncoding.EncodeToString([]byte("short")),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset encryptionKey
			encryptionKey = nil

			if tt.keyEnv != "" {
				os.Setenv("ENCRYPTION_KEY", tt.keyEnv)
				defer os.Unsetenv("ENCRYPTION_KEY")
			} else {
				os.Unsetenv("ENCRYPTION_KEY")
			}

			err := Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecryptErrors(t *testing.T) {
	// Setup valid key
	testKey := make([]byte, 32)
	os.Setenv("ENCRYPTION_KEY", base64.StdEncoding.EncodeToString(testKey))
	defer os.Unsetenv("ENCRYPTION_KEY")

	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	tests := []struct {
		name       string
		ciphertext string
		wantErr    bool
	}{
		{
			name:       "invalid base64",
			ciphertext: "not-valid-base64!!!",
			wantErr:    true,
		},
		{
			name:       "too short",
			ciphertext: base64.StdEncoding.EncodeToString([]byte("short")),
			wantErr:    true,
		},
		{
			name:       "corrupted data",
			ciphertext: base64.StdEncoding.EncodeToString([]byte("this is definitely not a valid AES-GCM ciphertext with nonce")),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Decrypt(tt.ciphertext)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptWithoutInit(t *testing.T) {
	// Reset encryptionKey
	encryptionKey = nil

	_, err := Encrypt("test")
	if err == nil {
		t.Error("Encrypt should fail when not initialized")
	}
}

func TestDecryptWithoutInit(t *testing.T) {
	// Reset encryptionKey
	encryptionKey = nil

	_, err := Decrypt("test")
	if err == nil {
		t.Error("Decrypt should fail when not initialized")
	}
}
