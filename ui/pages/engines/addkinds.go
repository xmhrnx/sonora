package engines

import (
	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/engines/command"
	"github.com/mehran9675/sonora/internal/engines/openai"
	"github.com/mehran9675/sonora/ui/tr"
)

const (
	addAPI        = "api"
	addTranscribe = "transcribe"
	addCommand    = "command"
	addPlugin     = "plugin"
)

// addKinds maps the dialog's tabs to the engines they create; name is
// the default engine name, translated when it's given.
var addKinds = map[string]struct {
	kind     engines.Kind
	name     string
	template func(id, name string) *engines.Manifest
}{
	addAPI:        {engines.KindOpenAI, tr.Msg("API engine"), openai.Template},
	addTranscribe: {engines.KindOpenAITranscribe, tr.Msg("Transcription API"), openai.TranscribeTemplate},
	addCommand:    {engines.KindCommand, tr.Msg("Command engine"), command.Template},
}
