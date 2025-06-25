package keyboard

import (
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"github.com/atotto/clipboard"
	"github.com/kbinani/screenshot"
	"github.com/sashabaranov/go-openai"
)

type KeyboardOperator struct {
	listener     *KeyboardListener
	emulator     *KeyboardEmulator
	client       *openai.Client
	apiKey       string
	model        string
	baseURL      string
	systemPrompt string
	config       *KeyBindingConfig

	interactionMutex sync.Mutex
	isInteracting    bool
	cancelContext    context.CancelFunc
}

func NewKeyboardOperator(apiKey, model, baseURL, keyboardDevice, systemPrompt string) (*KeyboardOperator, error) {
	emulator, err := NewKeyboardEmulator()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize keyboard emulator: %v", err)
	}
	listener := NewKeyboardListener(keyboardDevice)
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	client := openai.NewClientWithConfig(config)

	// Use default system prompt if none provided
	if systemPrompt == "" {
		systemPrompt = "You are a helpful AI assistant. You may have access to the user's clipboard and/or a screenshot of their current screen, depending on the context. Use this context to provide more relevant and helpful responses."
	}

	// Load keybinding configuration from environment variables
	keyConfig := LoadKeyBindingConfig()

	return &KeyboardOperator{
		listener:     listener,
		emulator:     emulator,
		client:       client,
		apiKey:       apiKey,
		model:        model,
		baseURL:      baseURL,
		systemPrompt: systemPrompt,
		config:       keyConfig,
	}, nil
}

func (ko *KeyboardOperator) Close() {
	ko.StopCurrentInteraction()
	if ko.listener != nil {
		ko.listener.Stop()
	}
	if ko.emulator != nil {
		ko.emulator.Close()
	}
}

func (ko *KeyboardOperator) StopCurrentInteraction() {
	ko.interactionMutex.Lock()
	defer ko.interactionMutex.Unlock()
	if ko.isInteracting && ko.cancelContext != nil {
		ko.cancelContext()
		ko.isInteracting = false
		ko.cancelContext = nil
	}
}

func (ko *KeyboardOperator) showZenityInput() (string, error) {
	cmd := exec.Command("zenity", "--entry", "--title=Zeygeist", "--text=What would you like me to help you with?", "--width=400", "--height=100")
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
		}
		return "", fmt.Errorf("zenity error: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (ko *KeyboardOperator) getClipboardContent() string {
	content, err := clipboard.ReadAll()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(content)
}

func (ko *KeyboardOperator) takeScreenshotBase64() ([]string, error) {
	displays := screenshot.NumActiveDisplays()
	var screenshots []string

	// Try using the screenshot library first
	for i := 0; i < displays; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			fmt.Printf("Screenshot library failed for display %d: %v, trying fallback tools\n", i, err)
			// Fallback to external tools for this display
			fallbackScreenshot, err := ko.takeScreenshotWithFallbackTools()
			if err != nil {
				return nil, fmt.Errorf("both screenshot library and fallback tools failed for display %d: %v", i, err)
			}
			screenshots = append(screenshots, fallbackScreenshot)
			break // Use one fallback screenshot for all displays
		}

		var buf strings.Builder
		encoder := base64.NewEncoder(base64.StdEncoding, &buf)
		if err := png.Encode(encoder, img); err != nil {
			fmt.Printf("Failed to encode screenshot for display %d: %v, trying fallback tools\n", i, err)
			// Fallback to external tools for this display
			fallbackScreenshot, err := ko.takeScreenshotWithFallbackTools()
			if err != nil {
				return nil, fmt.Errorf("both screenshot library and fallback tools failed for display %d: %v", i, err)
			}
			screenshots = append(screenshots, fallbackScreenshot)
			break // Use one fallback screenshot for all displays
		}
		encoder.Close()
		screenshots = append(screenshots, buf.String())
	}

	if len(screenshots) == 0 {
		// If no displays were found or all failed, try fallback tools
		fmt.Printf("No displays found with screenshot library, trying fallback tools\n")
		fallbackScreenshot, err := ko.takeScreenshotWithFallbackTools()
		if err != nil {
			return nil, fmt.Errorf("no displays found and fallback tools failed: %v", err)
		}
		screenshots = append(screenshots, fallbackScreenshot)
	}

	return screenshots, nil
}

func (ko *KeyboardOperator) takeScreenshotWithFallbackTools() (string, error) {
	// Try multiple screenshot tools in order of preference
	tools := []struct {
		name string
		args []string
	}{
		{"grim", []string{}},                 // Wayland - saves to file
		{"gnome-screenshot", []string{"-f"}}, // GNOME - saves to file
		{"scrot", []string{}},                // X11 - saves to file
	}

	for _, tool := range tools {
		fmt.Printf("Trying screenshot tool: %s\n", tool.name)
		screenshot, err := ko.tryScreenshotTool(tool.name, tool.args)
		if err == nil {
			fmt.Printf("Successfully captured screenshot with %s\n", tool.name)
			return screenshot, nil
		}
		fmt.Printf("Failed to capture screenshot with %s: %v\n", tool.name, err)
	}

	return "", fmt.Errorf("all screenshot tools failed")
}

func (ko *KeyboardOperator) tryScreenshotTool(toolName string, args []string) (string, error) {
	// Check if tool exists
	if _, err := exec.LookPath(toolName); err != nil {
		return "", fmt.Errorf("tool not found: %v", err)
	}

	// Create a temporary file for the screenshot
	tmpFile, err := ioutil.TempFile("", "screenshot-*.png")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Build command with filename
	cmdArgs := append(args, tmpFile.Name())
	cmd := exec.Command(toolName, cmdArgs...)
	cmd.Env = os.Environ()

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("tool failed: %v, output: %s", err, string(output))
	}

	// Read the file and encode to base64
	imageData, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read screenshot file: %v", err)
	}

	// Convert to base64
	return base64.StdEncoding.EncodeToString(imageData), nil
}

func (ko *KeyboardOperator) handleCombinationContext(ctxType string) func() {
	return func() {
		ko.interactionMutex.Lock()
		if ko.isInteracting {
			if ko.cancelContext != nil {
				ko.cancelContext()
			}
			ko.isInteracting = false
			ko.cancelContext = nil
			ko.interactionMutex.Unlock()
			return
		}
		ko.isInteracting = true
		ctx, cancel := context.WithCancel(context.Background())
		ko.cancelContext = cancel
		ko.interactionMutex.Unlock()
		go func() {
			defer func() {
				ko.interactionMutex.Lock()
				ko.isInteracting = false
				ko.cancelContext = nil
				ko.interactionMutex.Unlock()
			}()

			// Type backspace to clear the combination that was pressed
			ko.emulator.TypeText("\b")

			// Take screenshot before showing Zenity dialog
			var screenshots []string
			if ctxType == "screenshot" || ctxType == "all" {
				var err error
				screenshots, err = ko.takeScreenshotBase64()
				if err != nil {
					fmt.Printf("Failed to take screenshot: %v\n", err)
					screenshots = nil
				}
			}

			input, err := ko.showZenityInput()
			if err != nil || input == "" {
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
			response, err := ko.queryOpenAIWithContext(ctx, input, ctxType, screenshots)
			if err != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
			// Clean the response before typing
			cleanedResponse := ko.cleanResponse(response)
			_ = ko.emulator.TypeText(cleanedResponse)
		}()
	}
}

func (ko *KeyboardOperator) queryOpenAIWithContext(ctx context.Context, prompt, ctxType string, screenshots []string) (string, error) {
	clipboardContent := ""
	if ctxType == "clipboard" || ctxType == "all" {
		clipboardContent = ko.getClipboardContent()
	}

	systemMessage := ko.systemPrompt
	userMessage := fmt.Sprintf("User question: %s\n\n", prompt)
	if clipboardContent != "" {
		userMessage += fmt.Sprintf("Clipboard content:\n%s\n\n", clipboardContent)
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemMessage,
		},
	}

	if len(screenshots) > 0 {
		fmt.Printf("Sending %d screenshots to OpenAI\n", len(screenshots))

		// Create multi-content message with text and all screenshots
		multiContent := []openai.ChatMessagePart{
			{
				Type: openai.ChatMessagePartTypeText,
				Text: userMessage,
			},
		}

		// Add all screenshots as image parts
		for _, screenshot := range screenshots {
			dataURL := "data:image/png;base64," + screenshot
			multiContent = append(multiContent, openai.ChatMessagePart{
				Type: openai.ChatMessagePartTypeImageURL,
				ImageURL: &openai.ChatMessageImageURL{
					URL: dataURL,
				},
			})
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:         openai.ChatMessageRoleUser,
			MultiContent: multiContent,
		})
	} else {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: userMessage,
		})
	}

	resp, err := ko.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    ko.model,
			Messages: messages,
		},
	)
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %v", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}
	return resp.Choices[0].Message.Content, nil
}

func (ko *KeyboardOperator) Start() error {
	// Parse key combinations from configuration
	clipboardKeys, err := ParseKeyCombination(ko.config.ClipboardKey)
	if err != nil {
		return fmt.Errorf("invalid clipboard key combination '%s': %v", ko.config.ClipboardKey, err)
	}

	screenshotKeys, err := ParseKeyCombination(ko.config.ScreenshotKey)
	if err != nil {
		return fmt.Errorf("invalid screenshot key combination '%s': %v", ko.config.ScreenshotKey, err)
	}

	allContextKeys, err := ParseKeyCombination(ko.config.AllContextKey)
	if err != nil {
		return fmt.Errorf("invalid all context key combination '%s': %v", ko.config.AllContextKey, err)
	}

	textOnlyKeys, err := ParseKeyCombination(ko.config.TextOnlyKey)
	if err != nil {
		return fmt.Errorf("invalid text only key combination '%s': %v", ko.config.TextOnlyKey, err)
	}

	// Register combinations with parsed keys
	ko.listener.AddCombination("clipboard", clipboardKeys...)
	ko.listener.OnCombination("clipboard", ko.handleCombinationContext("clipboard"))

	ko.listener.AddCombination("screenshot", screenshotKeys...)
	ko.listener.OnCombination("screenshot", ko.handleCombinationContext("screenshot"))

	ko.listener.AddCombination("all", allContextKeys...)
	ko.listener.OnCombination("all", ko.handleCombinationContext("all"))

	ko.listener.AddCombination("textonly", textOnlyKeys...)
	ko.listener.OnCombination("textonly", ko.handleCombinationContext("textonly"))

	return ko.listener.Start()
}

func (ko *KeyboardOperator) cleanResponse(response string) string {
	// Unwrap code blocks with language specifier (```lang\n...```)
	codeBlockWithLang := regexp.MustCompile("(?m)```[a-zA-Z0-9_+-]*\\n([\\w\\W]*?)```[ \t\r\n]*")
	response = codeBlockWithLang.ReplaceAllString(response, "$1")

	// Unwrap code blocks without language specifier (```...```)
	codeBlock := regexp.MustCompile("(?m)```\\n?([\\w\\W]*?)```[ \t\r\n]*")
	response = codeBlock.ReplaceAllString(response, "$1")

	// Remove thinking tags (<thinking>...</thinking>)
	thinkingRegex := regexp.MustCompile("(?m)<thinking>[\\w\\W]*?</thinking>")
	response = thinkingRegex.ReplaceAllString(response, "")

	// Remove think tags (<think>...</think>)
	thinkRegex := regexp.MustCompile("(?m)<think>[\\w\\W]*?</think>")
	response = thinkRegex.ReplaceAllString(response, "")

	// Clean up extra whitespace and newlines
	response = strings.TrimSpace(response)

	// Remove multiple consecutive newlines
	newlineRegex := regexp.MustCompile("\\n{3,}")
	response = newlineRegex.ReplaceAllString(response, "\n\n")

	return response
}

func (ko *KeyboardOperator) GetConfig() *KeyBindingConfig {
	return ko.config
}
