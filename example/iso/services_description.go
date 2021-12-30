package proto

import (
	"github.com/Speakerkfm/iso/pkg/models"

	"github.com/Speakerkfm/iso/pkg/proto/pb/service_b"
	"github.com/Speakerkfm/iso/pkg/proto/pb/service_c"
)

var Services = []*models.ProtoService{
	{
		Name: "UserService",
		Methods: []models.ProtoMethod{
			{
				Name:           "GetUser",
				RequestStruct:  &service_b.GetUserRequest{},
				ResponseStruct: &service_b.GetUserResponse{},
			},
		},
		ProtoPath: "pb/service_b.proto",
	},
	{
		Name: "PhoneService",
		Methods: []models.ProtoMethod{
			{
				Name:           "CheckPhone",
				RequestStruct:  &service_c.CheckPhoneRequest{},
				ResponseStruct: &service_c.CheckPhoneResponse{},
			},
		},
		ProtoPath: "pb/service_c.proto",
	},
}
