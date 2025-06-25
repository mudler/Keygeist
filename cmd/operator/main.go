package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/mudler/keygeist/keyboard"
)

func main() {
	if _, err := exec.LookPath("zenity"); err != nil {
		log.Fatal("zenity is not installed. Please install it first.")
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		log.Fatal("OPENAI_MODEL environment variable is not set")
	}
	baseURL := os.Getenv("OPENAI_BASE_URL")
	systemPrompt := os.Getenv("OPENAI_SYSTEM_PROMPT")
	operator, err := keyboard.NewKeyboardOperator(apiKey, model, baseURL, os.Getenv("KEYBOARD_DEVICE"), systemPrompt)
	if err != nil {
		log.Fatalf("Failed to create Keygeist: %v", err)
	}
	defer operator.Close()
	fmt.Println("Keygeist initialized!")
	fmt.Println("Press the configured key combinations:")
	fmt.Printf("  - %s for clipboard context\n", operator.GetConfig().ClipboardKey)
	fmt.Printf("  - %s for screenshot context\n", operator.GetConfig().ScreenshotKey)
	fmt.Printf("  - %s for all context\n", operator.GetConfig().AllContextKey)
	fmt.Printf("  - %s for text-only context\n", operator.GetConfig().TextOnlyKey)
	fmt.Println("Press the same combination again to stop current interaction")
	fmt.Println("Press Ctrl+C to exit")
	if err := operator.Start(); err != nil {
		log.Fatalf("Failed to start Keygeist: %v", err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("\nExiting...")
}
