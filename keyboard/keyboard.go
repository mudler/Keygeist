package keyboard

import (
	"fmt"
	"time"

	"github.com/bendahl/uinput"
)

type KeyboardEmulator struct {
	keyboard uinput.Keyboard
}

func NewKeyboardEmulator() (*KeyboardEmulator, error) {
	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("Keyboard Emulator"))
	if err != nil {
		return nil, fmt.Errorf("failed to create keyboard: %v", err)
	}

	return &KeyboardEmulator{
		keyboard: keyboard,
	}, nil
}

func (ke *KeyboardEmulator) Close() {
	ke.keyboard.Close()
}

func (ke *KeyboardEmulator) PressKey(keyCode int) error {
	return ke.keyboard.KeyDown(keyCode)
}

func (ke *KeyboardEmulator) ReleaseKey(keyCode int) error {
	return ke.keyboard.KeyUp(keyCode)
}

func (ke *KeyboardEmulator) TapKey(keyCode int) error {
	if err := ke.PressKey(keyCode); err != nil {
		return err
	}
	// Avoid buffer issues
	time.Sleep(1 * time.Millisecond)
	return ke.ReleaseKey(keyCode)
}

func (ke *KeyboardEmulator) TypeText(text string) error {
	for _, char := range text {
		keyCode, shift := charToKeyCode(char)
		if keyCode != 0 {
			if shift {
				if err := ke.PressKey(int(uinput.KeyLeftshift)); err != nil {
					return err
				}
			}
			if err := ke.TapKey(keyCode); err != nil {
				if shift {
					_ = ke.ReleaseKey(int(uinput.KeyLeftshift))
				}
				return err
			}
			if shift {
				if err := ke.ReleaseKey(int(uinput.KeyLeftshift)); err != nil {
					return err
				}
			}
			//time.Sleep(10 * time.Millisecond)
		}
	}
	return nil
}

func (ke *KeyboardEmulator) PressHotkey(keys ...int) error {
	// Press all keys
	for _, key := range keys {
		if err := ke.PressKey(key); err != nil {
			return err
		}
	}

	time.Sleep(100 * time.Millisecond)

	// Release all keys in reverse order
	for i := len(keys) - 1; i >= 0; i-- {
		if err := ke.ReleaseKey(keys[i]); err != nil {
			return err
		}
	}

	return nil
}

// charToKeyCode returns the keycode and whether shift is required for a given rune
func charToKeyCode(char rune) (keyCode int, shift bool) {
	switch char {
	// Lowercase letters
	case 'a':
		return int(uinput.KeyA), false
	case 'b':
		return int(uinput.KeyB), false
	case 'c':
		return int(uinput.KeyC), false
	case 'd':
		return int(uinput.KeyD), false
	case 'e':
		return int(uinput.KeyE), false
	case 'f':
		return int(uinput.KeyF), false
	case 'g':
		return int(uinput.KeyG), false
	case 'h':
		return int(uinput.KeyH), false
	case 'i':
		return int(uinput.KeyI), false
	case 'j':
		return int(uinput.KeyJ), false
	case 'k':
		return int(uinput.KeyK), false
	case 'l':
		return int(uinput.KeyL), false
	case 'm':
		return int(uinput.KeyM), false
	case 'n':
		return int(uinput.KeyN), false
	case 'o':
		return int(uinput.KeyO), false
	case 'p':
		return int(uinput.KeyP), false
	case 'q':
		return int(uinput.KeyQ), false
	case 'r':
		return int(uinput.KeyR), false
	case 's':
		return int(uinput.KeyS), false
	case 't':
		return int(uinput.KeyT), false
	case 'u':
		return int(uinput.KeyU), false
	case 'v':
		return int(uinput.KeyV), false
	case 'w':
		return int(uinput.KeyW), false
	case 'x':
		return int(uinput.KeyX), false
	case 'y':
		return int(uinput.KeyY), false
	case 'z':
		return int(uinput.KeyZ), false
	// Uppercase letters
	case 'A':
		return int(uinput.KeyA), true
	case 'B':
		return int(uinput.KeyB), true
	case 'C':
		return int(uinput.KeyC), true
	case 'D':
		return int(uinput.KeyD), true
	case 'E':
		return int(uinput.KeyE), true
	case 'F':
		return int(uinput.KeyF), true
	case 'G':
		return int(uinput.KeyG), true
	case 'H':
		return int(uinput.KeyH), true
	case 'I':
		return int(uinput.KeyI), true
	case 'J':
		return int(uinput.KeyJ), true
	case 'K':
		return int(uinput.KeyK), true
	case 'L':
		return int(uinput.KeyL), true
	case 'M':
		return int(uinput.KeyM), true
	case 'N':
		return int(uinput.KeyN), true
	case 'O':
		return int(uinput.KeyO), true
	case 'P':
		return int(uinput.KeyP), true
	case 'Q':
		return int(uinput.KeyQ), true
	case 'R':
		return int(uinput.KeyR), true
	case 'S':
		return int(uinput.KeyS), true
	case 'T':
		return int(uinput.KeyT), true
	case 'U':
		return int(uinput.KeyU), true
	case 'V':
		return int(uinput.KeyV), true
	case 'W':
		return int(uinput.KeyW), true
	case 'X':
		return int(uinput.KeyX), true
	case 'Y':
		return int(uinput.KeyY), true
	case 'Z':
		return int(uinput.KeyZ), true
	// Numbers
	case '0':
		return int(uinput.Key0), false
	case '1':
		return int(uinput.Key1), false
	case '2':
		return int(uinput.Key2), false
	case '3':
		return int(uinput.Key3), false
	case '4':
		return int(uinput.Key4), false
	case '5':
		return int(uinput.Key5), false
	case '6':
		return int(uinput.Key6), false
	case '7':
		return int(uinput.Key7), false
	case '8':
		return int(uinput.Key8), false
	case '9':
		return int(uinput.Key9), false
	// Symbols above numbers
	case '!':
		return int(uinput.Key1), true
	case '@':
		return int(uinput.Key2), true
	case '#':
		return int(uinput.Key3), true
	case '$':
		return int(uinput.Key4), true
	case '%':
		return int(uinput.Key5), true
	case '^':
		return int(uinput.Key6), true
	case '&':
		return int(uinput.Key7), true
	case '*':
		return int(uinput.Key8), true
	case '(':
		return int(uinput.Key9), true
	case ')':
		return int(uinput.Key0), true
	// Other symbols
	case '-':
		return int(uinput.KeyMinus), false
	case '_':
		return int(uinput.KeyMinus), true
	case '=':
		return int(uinput.KeyEqual), false
	case '+':
		return int(uinput.KeyEqual), true
	case '[':
		return int(uinput.KeyLeftbrace), false
	case '{':
		return int(uinput.KeyLeftbrace), true
	case ']':
		return int(uinput.KeyRightbrace), false
	case '}':
		return int(uinput.KeyRightbrace), true
	case '\\':
		return int(uinput.KeyBackslash), false
	case '|':
		return int(uinput.KeyBackslash), true
	case ';':
		return int(uinput.KeySemicolon), false
	case ':':
		return int(uinput.KeySemicolon), true
	case '\'':
		return int(uinput.KeyApostrophe), false
	case '"':
		return int(uinput.KeyApostrophe), true
	case '`':
		return int(uinput.KeyGrave), false
	case '~':
		return int(uinput.KeyGrave), true
	case ',':
		return int(uinput.KeyComma), false
	case '<':
		return int(uinput.KeyComma), true
	case '.':
		return int(uinput.KeyDot), false
	case '>':
		return int(uinput.KeyDot), true
	case '/':
		return int(uinput.KeySlash), false
	case '?':
		return int(uinput.KeySlash), true
	case ' ':
		return int(uinput.KeySpace), false
	case '\n':
		return int(uinput.KeyEnter), false
	case '\t':
		return int(uinput.KeyTab), false
	case '\b':
		return int(uinput.KeyBackspace), false
	}
	return 0, false
}
