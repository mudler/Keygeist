package keyboard

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

// InputEvent is the Linux input_event struct
type InputEvent struct {
	Time  unix.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

// KeyState represents the state of a key (pressed/released)
type KeyState bool

const (
	KeyReleased KeyState = false
	KeyPressed  KeyState = true
)

// KeyCombination represents a combination of keys to listen for
type KeyCombination struct {
	Name string
	Keys []uint16
}

// KeyboardListener listens for keyboard events and detects key combinations
type KeyboardListener struct {
	devicePath   string
	file         *os.File
	keyStates    map[uint16]KeyState
	combinations []KeyCombination
	callbacks    map[string][]func()
	running      bool
}

// NewKeyboardListener creates a new keyboard listener
func NewKeyboardListener(devicePath string) *KeyboardListener {
	return &KeyboardListener{
		devicePath:   devicePath,
		keyStates:    make(map[uint16]KeyState),
		combinations: make([]KeyCombination, 0),
		callbacks:    make(map[string][]func()),
	}
}

// FindKeyboardDevice attempts to find a keyboard device automatically
func (kl *KeyboardListener) FindKeyboardDevice() (string, error) {
	// Check /dev/input/by-path for keyboard devices
	byPathDir := "/dev/input/by-path"
	entries, err := os.ReadDir(byPathDir)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %v", byPathDir, err)
	}

	for _, entry := range entries {
		if strings.Contains(entry.Name(), "kbd") {
			devicePath := filepath.Join(byPathDir, entry.Name())
			// Resolve symlink to actual device
			actualPath, err := os.Readlink(devicePath)
			if err != nil {
				continue
			}
			fmt.Println("Found keyboard device:", actualPath)
			if !filepath.IsAbs(actualPath) {
				actualPath = filepath.Join("/dev/input/by-path", actualPath)
			}
			return actualPath, nil
		}
	}

	// Fallback: try common event devices
	commonDevices := []string{
		"/dev/input/event0",
		"/dev/input/event1",
		"/dev/input/event2",
		"/dev/input/event3",
		"/dev/input/event4",
		"/dev/input/event5",
	}

	for _, device := range commonDevices {
		if _, err := os.Stat(device); err == nil {
			return device, nil
		}
	}

	return "", fmt.Errorf("no keyboard device found")
}

// SetDevice sets the input device to listen on
func (kl *KeyboardListener) SetDevice(devicePath string) {
	kl.devicePath = devicePath
}

// AddCombination adds a key combination to listen for
func (kl *KeyboardListener) AddCombination(name string, keys ...uint16) {
	combination := KeyCombination{
		Name: name,
		Keys: keys,
	}
	kl.combinations = append(kl.combinations, combination)
}

// OnCombination registers a callback for when a combination is detected
func (kl *KeyboardListener) OnCombination(name string, callback func()) {
	kl.callbacks[name] = append(kl.callbacks[name], callback)
}

// Start begins listening for keyboard events
func (kl *KeyboardListener) Start() error {
	if kl.devicePath == "" {
		devicePath, err := kl.FindKeyboardDevice()
		if err != nil {
			return fmt.Errorf("failed to find keyboard device: %v", err)
		}
		kl.devicePath = devicePath
	}

	file, err := os.Open(kl.devicePath)
	if err != nil {
		return fmt.Errorf("failed to open device %s: %v", kl.devicePath, err)
	}
	kl.file = file
	kl.running = true

	fmt.Printf("Listening for keyboard events on %s\n", kl.devicePath)
	fmt.Printf("Registered combinations: %v\n", kl.getCombinationNames())

	go kl.listenLoop()
	return nil
}

// Stop stops listening for keyboard events
func (kl *KeyboardListener) Stop() {
	kl.running = false
	if kl.file != nil {
		kl.file.Close()
	}
}

// listenLoop is the main event listening loop
func (kl *KeyboardListener) listenLoop() {
	for kl.running {
		var event InputEvent
		err := kl.binaryRead(&event)
		if err != nil {
			if kl.running {
				fmt.Printf("Read error: %v\n", err)
			}
			break
		}

		if event.Type == EV_KEY {
			kl.handleKeyEvent(event)
		}
	}
}

// handleKeyEvent processes a key event and checks for combinations
func (kl *KeyboardListener) handleKeyEvent(event InputEvent) {
	keyState := event.Value != 0
	kl.keyStates[event.Code] = KeyState(keyState)

	// Check all combinations
	for _, combination := range kl.combinations {
		if kl.isCombinationActive(combination) {
			kl.triggerCallbacks(combination.Name)
		}
	}
}

// isCombinationActive checks if a key combination is currently active
func (kl *KeyboardListener) isCombinationActive(combination KeyCombination) bool {
	for _, key := range combination.Keys {
		if !kl.keyStates[key] {
			return false
		}
	}
	return true
}

// triggerCallbacks executes all callbacks for a combination
func (kl *KeyboardListener) triggerCallbacks(combinationName string) {
	if callbacks, exists := kl.callbacks[combinationName]; exists {
		for _, callback := range callbacks {
			callback()
		}
	}
}

// getCombinationNames returns a list of registered combination names
func (kl *KeyboardListener) getCombinationNames() []string {
	names := make([]string, len(kl.combinations))
	for i, combo := range kl.combinations {
		names[i] = combo.Name
	}
	return names
}

// binaryRead reads InputEvent from device file
func (kl *KeyboardListener) binaryRead(event *InputEvent) error {
	buf := make([]byte, 24)
	_, err := kl.file.Read(buf)
	if err != nil {
		return err
	}

	event.Time.Sec = int64(binary.LittleEndian.Uint64(buf[0:8]))
	event.Time.Usec = int64(binary.LittleEndian.Uint64(buf[8:16]))
	event.Type = binary.LittleEndian.Uint16(buf[16:18])
	event.Code = binary.LittleEndian.Uint16(buf[18:20])
	event.Value = int32(binary.LittleEndian.Uint32(buf[20:24]))

	return nil
}

// Common key codes
const (
	EV_KEY         = 0x01
	KEY_ESC        = 1
	KEY_1          = 2
	KEY_2          = 3
	KEY_3          = 4
	KEY_4          = 5
	KEY_5          = 6
	KEY_6          = 7
	KEY_7          = 8
	KEY_8          = 9
	KEY_9          = 10
	KEY_0          = 11
	KEY_MINUS      = 12
	KEY_EQUAL      = 13
	KEY_BACKSPACE  = 14
	KEY_TAB        = 15
	KEY_Q          = 16
	KEY_W          = 17
	KEY_E          = 18
	KEY_R          = 19
	KEY_T          = 20
	KEY_Y          = 21
	KEY_U          = 22
	KEY_I          = 23
	KEY_O          = 24
	KEY_P          = 25
	KEY_LEFTBRACE  = 26
	KEY_RIGHTBRACE = 27
	KEY_ENTER      = 28
	KEY_LEFTCTRL   = 29
	KEY_A          = 30
	KEY_S          = 31
	KEY_D          = 32
	KEY_F          = 33
	KEY_G          = 34
	KEY_H          = 35
	KEY_J          = 36
	KEY_K          = 37
	KEY_L          = 38
	KEY_SEMICOLON  = 39
	KEY_APOSTROPHE = 40
	KEY_GRAVE      = 41
	KEY_LEFTSHIFT  = 42
	KEY_BACKSLASH  = 43
	KEY_Z          = 44
	KEY_X          = 45
	KEY_C          = 46
	KEY_V          = 47
	KEY_B          = 48
	KEY_N          = 49
	KEY_M          = 50
	KEY_COMMA      = 51
	KEY_DOT        = 52
	KEY_SLASH      = 53
	KEY_RIGHTSHIFT = 54
	KEY_LEFTALT    = 56
	KEY_SPACE      = 57
	KEY_CAPSLOCK   = 58
	KEY_F1         = 59
	KEY_F2         = 60
	KEY_F3         = 61
	KEY_F4         = 62
	KEY_F5         = 63
	KEY_F6         = 64
	KEY_F7         = 65
	KEY_F8         = 66
	KEY_F9         = 67
	KEY_F10        = 68
	KEY_F11        = 87
	KEY_F12        = 88
	KEY_LEFTMETA   = 125
	KEY_RIGHTMETA  = 126
)
