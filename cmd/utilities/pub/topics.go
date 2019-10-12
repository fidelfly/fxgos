package pub

const (
	TopicRoleUpdate = "topic.role.update"
	TopicResource   = "topic.resource"
)

const (
	ResourceCreate = iota
	ResourceUpdate
	ResourceDelete
)

type ResourceEvent struct {
	Type   string
	Action int
	Id     int64
	Code   string
}
