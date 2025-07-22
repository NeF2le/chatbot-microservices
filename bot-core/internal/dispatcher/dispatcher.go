package dispatcher

import (
	"errors"
	"fmt"
	"github.com/NeF2le/chatbot-microservices/bot-core/internal/models"
)

type Dispatcher struct {
	Skills map[string]Skill
}

func New(skills map[string]Skill) *Dispatcher {
	return &Dispatcher{Skills: skills}
}

func (d *Dispatcher) Dispatch(msg models.Message) (models.Reply, error) {
	for name, skill := range d.Skills {
		if skill.Match(msg) {
			reply, err := skill.Execute(msg)
			if err != nil {
				return models.Reply{}, fmt.Errorf("error when executing skill %s", name)
			}
			return models.Reply{Reply: reply}, nil
		}
	}
	return models.Reply{}, errors.New("skill not found")
}
