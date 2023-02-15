package api

import (
	"github.com/RobleDev498/spaces/app"
	"github.com/gorilla/mux"
)

const (
	API_URL_SUFFIX = "/api"
)

type Routes struct {
	Root    *mux.Router
	ApiRoot *mux.Router

	Users *mux.Router // 'api/users'
	User  *mux.Router // 'api/users/{user_id:[0-9]+}'

	Spaces        *mux.Router // 'api/spaces'
	Space         *mux.Router // 'api/spaces/{space_id:[A-Za-z0-9@]+}'
	SpacesForUser *mux.Router // 'api/users/{user_id:[0-9]+}/spaces'
	SpaceMembers  *mux.Router // 'api/spaces/{space_id:[0-9]+}/members'

	Streams         *mux.Router // 'api/streams'
	Stream          *mux.Router // 'api/streams/{stream_id:[0-9]+}'
	StreamsForSpace *mux.Router // 'api/spaces/{space_id:[0-9]+}/streams'
	StreamMembers   *mux.Router // 'api/streams/{stream_id:[0-9]+}/members'
	StreamsDefault  *mux.Router // 'api/streams/default'

	RoomsForUser *mux.Router // 'api/users/{user_id:[0-9]+}/rooms'

	Messages          *mux.Router // 'api/messages'
	Message           *mux.Router // 'api/messages/{message_id:[0-9]+}'
	MessagesForStream *mux.Router // 'api/streams/{stream_id:[0-9]+}/messages'
	MessagesForUser   *mux.Router // 'api/users/{user_id:[0-9]+}/messages'
}

type Api struct {
	app        *app.App
	BaseRoutes *Routes
}

func Init(a *app.App, root *mux.Router) *Api {
	api := &Api{
		app:        a,
		BaseRoutes: &Routes{},
	}

	api.BaseRoutes.Root = root
	api.BaseRoutes.ApiRoot = root.PathPrefix(API_URL_SUFFIX).Subrouter()

	api.BaseRoutes.Users = api.BaseRoutes.ApiRoot.PathPrefix("/users").Subrouter()
	api.BaseRoutes.User = api.BaseRoutes.ApiRoot.PathPrefix("/users/{user_id:[A-Za-z0-9]+}").Subrouter()

	api.BaseRoutes.Spaces = api.BaseRoutes.ApiRoot.PathPrefix("/spaces").Subrouter()
	api.BaseRoutes.Space = api.BaseRoutes.ApiRoot.PathPrefix("/spaces/{{space_id:[A-Za-z0-9]+}}").Subrouter()
	api.BaseRoutes.SpacesForUser = api.BaseRoutes.ApiRoot.PathPrefix("/users/{user_id:[A-Za-z0-9]+}/spaces").Subrouter()
	api.BaseRoutes.SpaceMembers = api.BaseRoutes.ApiRoot.PathPrefix("/spaces/{{space_id:[A-Za-z0-9]+}}/members").Subrouter()

	api.BaseRoutes.Streams = api.BaseRoutes.ApiRoot.PathPrefix("/streams").Subrouter()
	api.BaseRoutes.Stream = api.BaseRoutes.ApiRoot.PathPrefix("/streams/{stream_id:[A-Za-z0-9]+}").Subrouter()
	api.BaseRoutes.StreamsForSpace = api.BaseRoutes.ApiRoot.PathPrefix("/spaces/{{space_id:[A-Za-z0-9]+}}/streams").Subrouter()
	api.BaseRoutes.StreamMembers = api.BaseRoutes.ApiRoot.PathPrefix("/streams/{stream_id:[A-Za-z0-9]+}/members").Subrouter()

	api.BaseRoutes.StreamsDefault = api.BaseRoutes.ApiRoot.PathPrefix("/streams/default").Subrouter()

	api.BaseRoutes.RoomsForUser = api.BaseRoutes.ApiRoot.PathPrefix("/users/{user_id:[A-Za-z0-9]+}/rooms").Subrouter()

	api.BaseRoutes.Messages = api.BaseRoutes.ApiRoot.PathPrefix("/messages").Subrouter()
	api.BaseRoutes.Message = api.BaseRoutes.ApiRoot.PathPrefix("/message/{{message_id:[A-Za-z0-9]+}}").Subrouter()
	api.BaseRoutes.MessagesForStream = api.BaseRoutes.ApiRoot.PathPrefix("/streams/{stream_id:[A-Za-z0-9]+}/messages").Subrouter()
	api.BaseRoutes.MessagesForUser = api.BaseRoutes.ApiRoot.PathPrefix("/users/{user_id:[A-Za-z0-9]+}/messages").Subrouter()

	api.InitUser()
	api.InitStream()
	api.InitSpace()
	api.InitMessage()
	api.InitWebSocket()
	return api
}
