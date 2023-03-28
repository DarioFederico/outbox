package events

type CategoryEventName string

const (
	CreateCategoryEventV1_Name CategoryEventName = "CreateCategoryEventV1"
	UpdateCategoryEventV1_Name CategoryEventName = "UpdateCategoryEventV1"
)

type CategoryEvent interface {
	GetEventName() CategoryEventName
}
