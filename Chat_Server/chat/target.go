package chat

type Target interface {
	IsTarget()
}

type RemoteTask struct {
	DNS        string
	ClientNick string
	Msg        string
}

func (rt *RemoteTask) IsTarget() {}
