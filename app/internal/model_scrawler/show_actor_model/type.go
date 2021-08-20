package show_actor_model

type ShowActor struct {
	Name string `json:"name"`
	Play string `json:"play"`
	PlayType string `json:"play_type"`
}

type ShowActors []ShowActor