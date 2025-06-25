package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mudler/keygeist/keyboard"
)

func printUsage() {
	fmt.Println("Keyboard Emulator Usage:")
	fmt.Println("  press <keycode>     - Press a key")
	fmt.Println("  release <keycode>   - Release a key")
	fmt.Println("  tap <keycode>       - Tap a key (press and release)")
	fmt.Println("  type <text>         - Type text")
	fmt.Println("  hotkey <key1> <key2> ... - Press multiple keys simultaneously")
	fmt.Println("  help                - Show this help")
	fmt.Println("  quit                - Exit the program")
	fmt.Println()
	fmt.Println("Common Key Codes:")
	fmt.Println("  A-Z: 30-57")
	fmt.Println("  0-9: 11-20")
	fmt.Println("  Space: 57")
	fmt.Println("  Enter: 28")
	fmt.Println("  Tab: 15")
	fmt.Println("  Ctrl: 29")
	fmt.Println("  Alt: 56")
	fmt.Println("  Shift: 42")
	fmt.Println("  Escape: 1")
}

func main() {
	fmt.Println("Initializing Keyboard Emulator...")

	ke, err := keyboard.NewKeyboardEmulator()
	if err != nil {
		log.Fatalf("Failed to initialize keyboard emulator: %v", err)
	}
	defer ke.Close()

	fmt.Println("Keyboard Emulator initialized successfully!")
	fmt.Println("Type 'help' for usage instructions.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]

		switch command {
		case "help":
			printUsage()

		case "quit", "exit":
			fmt.Println("Goodbye!")
			return

		case "press":
			if len(parts) < 2 {
				fmt.Println("Usage: press <keycode>")
				continue
			}
			keyCode, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Printf("Invalid key code: %s\n", parts[1])
				continue
			}
			if err := ke.PressKey(keyCode); err != nil {
				fmt.Printf("Error pressing key: %v\n", err)
			} else {
				fmt.Printf("Pressed key %d\n", keyCode)
			}

		case "release":
			if len(parts) < 2 {
				fmt.Println("Usage: release <keycode>")
				continue
			}
			keyCode, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Printf("Invalid key code: %s\n", parts[1])
				continue
			}
			if err := ke.ReleaseKey(keyCode); err != nil {
				fmt.Printf("Error releasing key: %v\n", err)
			} else {
				fmt.Printf("Released key %d\n", keyCode)
			}

		case "tap":
			if len(parts) < 2 {
				fmt.Println("Usage: tap <keycode>")
				continue
			}
			keyCode, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Printf("Invalid key code: %s\n", parts[1])
				continue
			}
			if err := ke.TapKey(keyCode); err != nil {
				fmt.Printf("Error tapping key: %v\n", err)
			} else {
				fmt.Printf("Tapped key %d\n", keyCode)
			}

		case "type":
			if len(parts) < 2 {
				fmt.Println("Usage: type <text>")
				continue
			}
			text := strings.Join(parts[1:], " ")
			if err := ke.TypeText(text); err != nil {
				fmt.Printf("Error typing text: %v\n", err)
			} else {
				fmt.Printf("Typed: %s\n", text)
			}

		case "hotkey":
			if len(parts) < 2 {
				fmt.Println("Usage: hotkey <key1> <key2> ...")
				continue
			}
			var keyCodes []int
			for _, part := range parts[1:] {
				keyCode, err := strconv.Atoi(part)
				if err != nil {
					fmt.Printf("Invalid key code: %s\n", part)
					continue
				}
				keyCodes = append(keyCodes, keyCode)
			}
			if len(keyCodes) > 0 {
				if err := ke.PressHotkey(keyCodes...); err != nil {
					fmt.Printf("Error pressing hotkey: %v\n", err)
				} else {
					fmt.Printf("Pressed hotkey: %v\n", keyCodes)
				}
			}

		default:
			fmt.Printf("Unknown command: %s. Type 'help' for usage.\n", command)
		}
	}
}
