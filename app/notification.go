package app

import "github.com/RobleDev498/spaces/model"

func (a *App) SendNotifications(msg *model.Message, sender *model.User, stream *model.Stream, space *model.Space) error {
	message := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_POSTED, "", msg.StreamID, "", nil)

	message.Add("post", message.ToJson())
	message.Add("space_id", space.ID)
	message.Add("stream_type", stream.Type)

	a.Publish(message)

	return nil
}
