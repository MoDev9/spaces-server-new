package store

import "github.com/RobleDev498/spaces/model"

func (s *Store) CreateMessage(msg *model.Message) (*model.Message, error) {
	msg.PreSave()
	result := s.DB.Create(&msg)
	return msg, result.Error
}

func (s *Store) GetMessages(streamId string) ([]*model.Message, error) {
	var messages []*model.Message
	result := s.DB.Where("stream_id = ?", streamId).Find(&messages)
	return messages, result.Error
}

func (s *Store) GetDefaultStream(userId string) (*model.Stream, error) {
	var stream *model.Stream
	//result := s .DB.Raw("SELECT * FROM messages WHERE author_id = ? ORDER BY created_at limit 1", userId).Scan(&stream)
	result := s.DB.Table("streams").First(&stream)
	return stream, result.Error
}
