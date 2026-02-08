package domain

import "fmt"

type ActorType string

const (
	ActorUser    ActorType = "USER"
	ActorAdmin   ActorType = "ADMIN"
	ActorSystem  ActorType = "SYSTEM"
	ActorService ActorType = "SERVICE"
)

var validActorTypes = map[ActorType]struct{}{
	ActorUser:    {},
	ActorAdmin:   {},
	ActorSystem:  {},
	ActorService: {},
}

func (a ActorType) IsValid() bool {
	_, ok := validActorTypes[a]
	return ok
}

func (a ActorType) Validate() error {
	if !a.IsValid() {
		return fmt.Errorf("invalid actor type: %q", a)
	}
	return nil
}
