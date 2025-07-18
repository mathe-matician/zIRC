package chat

import "github.com/phuslu/log"

const (
	saslC                = "sasl"
	accountRegistrationC = "account-registration"
)

type Capability interface {
	Execute()
	String() string
}

var CapMap = map[string]Capability{
	saslC:                &CAP_SASL{},
	accountRegistrationC: &CAP_AccountRegistration{},
}

type CAP_SASL struct {
	Enabled bool
}

func (sasl *CAP_SASL) Execute() {

}

func (sasl *CAP_SASL) String() string {
	return saslC
}

type CAP_AccountRegistration struct {
	Enabled bool
}

func (sasl *CAP_AccountRegistration) Execute() {

}

func (sasl *CAP_AccountRegistration) String() string {
	return accountRegistrationC
}

type ServerCapabilities struct {
	SASL                *CAP_SASL
	AccountRegistration *CAP_AccountRegistration
}

func NewServerCapabilities(capabilities ...string) *ServerCapabilities {
	serverCaps := ServerCapabilities{}
	for _, cap := range capabilities {
		if _, ok := CapMap[cap]; !ok {
			log.Info().Msgf("Unsupported capability: %s", cap)
			continue
		}
		// log.Debug().Msgf("cap %v", c)
	}

	return &serverCaps
}

func (sc *ServerCapabilities) String() {

}
