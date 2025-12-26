package terminal

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	e "github.com/JuanBiancuzzo/own_wiki/core/events"
	p "github.com/JuanBiancuzzo/own_wiki/core/platform"
	log "github.com/JuanBiancuzzo/own_wiki/core/system/logger"
	v "github.com/JuanBiancuzzo/own_wiki/view"

	ctxio "github.com/jbenet/go-context/io"
)

type TerminalPlatform struct {
	Reader       *bufio.Reader
	CancelReader context.CancelFunc
}

func NewTerminal() p.Platform {
	ctx, cancel := context.WithCancel(context.Background())
	stdin := ctxio.NewReader(ctx, os.Stdin)

	reader := bufio.NewReader(stdin)
	return &TerminalPlatform{
		Reader:       reader,
		CancelReader: cancel,
	}
}

func (hp *TerminalPlatform) HandleInput(eventQueue chan e.Event) {
	for {
		char, _, err := hp.Reader.ReadRune()
		if err == context.Canceled || err == io.EOF {
			break

		} else if err != nil {
			log.Error("Failed to read input from terminal, with error: '%v'", err)
			break
		}

		switch char {
		case 'h' | 'o' | 'l' | 'a' | 't' | 'n':
			eventQueue <- e.NewCharacterEvent(char)

		default:
			log.Debug("Character input: %d", char)
		}
	}

	log.Info("Closing HandleInput")
}

func (hp *TerminalPlatform) Render(viewRepresentation v.SceneRepresentation) {
	fmt.Printf("\033[H\033[2J")
	for _, valueRepresentation := range viewRepresentation {
		switch value := valueRepresentation.(type) {
		case *v.Heading:
			textHeading := fmt.Sprintf("%s %s", strings.Repeat("#", int(value.Level)), value.Data)
			fmt.Printf("%s\n%s\n", textHeading, strings.Repeat("-", len(textHeading)))

		case *v.Text:
			fmt.Printf("%s\n", value.Data)
		}
	}
}

func (hp *TerminalPlatform) Close() {
	hp.CancelReader()
}
