package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mudler/keygeist/keyboard"
)

func main() {
	kl := keyboard.NewKeyboardListener("")
	// Listen for Windows + I
	kl.AddCombination("win+i", keyboard.KEY_LEFTMETA, keyboard.KEY_I)
	kl.OnCombination("win+i", func() {
		fmt.Println("Detected: Windows + I")
	})

	if err := kl.Start(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer kl.Stop()

	// Wait for Ctrl+C to exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("Exiting...")
}
