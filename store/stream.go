package store

import (
	"log"

	"github.com/RobleDev498/spaces/model"
)

func (s *Store) CreateStream(stream *model.Stream) (*model.Stream, error) {
	stream.PreSave()
	result := s.DB.Create(&stream)
	if err := result.Error; err != nil {
		log.Println(err)
	}
	return stream, result.Error
}

func (s *Store) UpdateStream(stream *model.Stream) (*model.Stream, error) {
	stream.PreUpdate()
	result := s.DB.Save(&stream)
	if err := result.Error; err != nil {
		log.Println(err)
	}
	return stream, result.Error
}

func (s *Store) GetStream(streamId string) (*model.Stream, error) {
	var stream *model.Stream
	result := s.DB.First(&stream, "id = ?", streamId)
	return stream, result.Error
}

func (s *Store) AddFriendToRoom(roomMember *model.RoomMember) (*model.RoomMember, error) {
	result := s.DB.Table("room_members").Create(&roomMember)
	return roomMember, result.Error
}

func (s *Store) AddFriendsToRoom(roomMembers []*model.RoomMember) ([]*model.RoomMember, error) {
	result := s.DB.Table("room_members").Create(&roomMembers)
	return roomMembers, result.Error
}

func (s *Store) GetRooms(userId string) ([]*model.Stream, error) {
	var streams []*model.Stream
	//SELECT * FROM streams WHERE owner_id = ?
	result := s.DB.Raw(`SELECT s.* FROM streams s
	JOIN room_members r ON s.id = r.stream_id
	JOIN users u ON r.user_id = u.id
	WHERE u.id = ?`, userId).Scan(&streams)

	return streams, result.Error
}

func (s *Store) GetStreamMembers(streamId string) ([]*model.User, error) {
	var members []*model.User
	result := s.DB.Raw(`SELECT u.* FROM users u
	JOIN room_members r ON u.id = r.user_id 
	JOIN streams s ON r.stream_id = s.id 
	WHERE s.id = ?`, streamId).Scan(&members)

	return members, result.Error
}
