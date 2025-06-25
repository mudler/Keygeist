package keyboard

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// KeyBindingConfig represents the configuration for keybindings
type KeyBindingConfig struct {
	ClipboardKey  string
	ScreenshotKey string
	AllContextKey string
	TextOnlyKey   string
}

// DefaultKeyBindingConfig returns the default keybinding configuration
func DefaultKeyBindingConfig() *KeyBindingConfig {
	return &KeyBindingConfig{
		ClipboardKey:  "win+c",
		ScreenshotKey: "win+s",
		AllContextKey: "win+e",
		TextOnlyKey:   "win+t",
	}
}

// LoadKeyBindingConfig loads keybinding configuration from environment variables
func LoadKeyBindingConfig() *KeyBindingConfig {
	config := DefaultKeyBindingConfig()

	// Override with environment variables
	if env := os.Getenv("CLIPBOARD_KEY"); env != "" {
		config.ClipboardKey = env
	}
	if env := os.Getenv("SCREENSHOT_KEY"); env != "" {
		config.ScreenshotKey = env
	}
	if env := os.Getenv("ALL_CONTEXT_KEY"); env != "" {
		config.AllContextKey = env
	}
	if env := os.Getenv("TEXT_ONLY_KEY"); env != "" {
		config.TextOnlyKey = env
	}

	return config
}

// ParseKeyCombination parses a key combination string into individual key codes
func ParseKeyCombination(combination string) ([]uint16, error) {
	parts := strings.Split(strings.ToLower(combination), "+")
	var keys []uint16

	for _, part := range parts {
		part = strings.TrimSpace(part)
		keyCode, err := stringToKeyCode(part)
		if err != nil {
			return nil, fmt.Errorf("invalid key '%s' in combination '%s': %v", part, combination, err)
		}
		keys = append(keys, keyCode)
	}

	return keys, nil
}

// stringToKeyCode converts a string representation of a key to its key code
func stringToKeyCode(key string) (uint16, error) {
	switch key {
	case "ctrl", "control":
		return KEY_LEFTCTRL, nil
	case "alt":
		return KEY_LEFTALT, nil
	case "shift":
		return KEY_LEFTSHIFT, nil
	case "win", "windows", "meta":
		return KEY_LEFTMETA, nil
	case "a":
		return KEY_A, nil
	case "b":
		return KEY_B, nil
	case "c":
		return KEY_C, nil
	case "d":
		return KEY_D, nil
	case "e":
		return KEY_E, nil
	case "f":
		return KEY_F, nil
	case "g":
		return KEY_G, nil
	case "h":
		return KEY_H, nil
	case "i":
		return KEY_I, nil
	case "j":
		return KEY_J, nil
	case "k":
		return KEY_K, nil
	case "l":
		return KEY_L, nil
	case "m":
		return KEY_M, nil
	case "n":
		return KEY_N, nil
	case "o":
		return KEY_O, nil
	case "p":
		return KEY_P, nil
	case "q":
		return KEY_Q, nil
	case "r":
		return KEY_R, nil
	case "s":
		return KEY_S, nil
	case "t":
		return KEY_T, nil
	case "u":
		return KEY_U, nil
	case "v":
		return KEY_V, nil
	case "w":
		return KEY_W, nil
	case "x":
		return KEY_X, nil
	case "y":
		return KEY_Y, nil
	case "z":
		return KEY_Z, nil
	case "0":
		return KEY_0, nil
	case "1":
		return KEY_1, nil
	case "2":
		return KEY_2, nil
	case "3":
		return KEY_3, nil
	case "4":
		return KEY_4, nil
	case "5":
		return KEY_5, nil
	case "6":
		return KEY_6, nil
	case "7":
		return KEY_7, nil
	case "8":
		return KEY_8, nil
	case "9":
		return KEY_9, nil
	case "f1":
		return KEY_F1, nil
	case "f2":
		return KEY_F2, nil
	case "f3":
		return KEY_F3, nil
	case "f4":
		return KEY_F4, nil
	case "f5":
		return KEY_F5, nil
	case "f6":
		return KEY_F6, nil
	case "f7":
		return KEY_F7, nil
	case "f8":
		return KEY_F8, nil
	case "f9":
		return KEY_F9, nil
	case "f10":
		return KEY_F10, nil
	case "f11":
		return KEY_F11, nil
	case "f12":
		return KEY_F12, nil
	case "space":
		return KEY_SPACE, nil
	case "enter", "return":
		return KEY_ENTER, nil
	case "tab":
		return KEY_TAB, nil
	case "escape", "esc":
		return KEY_ESC, nil
	case "backspace":
		return KEY_BACKSPACE, nil
	default:
		// Try to parse as a numeric key code
		if code, err := strconv.ParseUint(key, 10, 16); err == nil {
			return uint16(code), nil
		}
		return 0, fmt.Errorf("unknown key: %s", key)
	}
}
