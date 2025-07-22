package dispatcher

import (
	"github.com/NeF2le/chatbot-microservices/bot-core/internal/models"
)

type Skill interface {
	Match(msg models.Message) bool
	Execute(msg models.Message) (string, error)
}
