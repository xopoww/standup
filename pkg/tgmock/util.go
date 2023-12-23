package tgmock

import (
	"strings"

	"github.com/xopoww/standup/pkg/tgmock/control"
)

//nolint:unparam // it will return error after more features are implemented
func parseMessageEntities(text string) ([]*control.MessageEntity, error) {
	if len(text) == 0 {
		return nil, nil
	}
	entities := make([]*control.MessageEntity, 0)

	// bot_command
	if text[0] == '/' {
		length := strings.IndexByte(text, ' ')
		if length == -1 {
			length = len(text)
		}
		entities = append(entities, &control.MessageEntity{
			Type:   "bot_command",
			Offset: 0,
			Length: int32(length),
		})
	}

	return entities, nil
}
