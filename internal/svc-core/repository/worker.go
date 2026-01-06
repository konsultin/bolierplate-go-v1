package repository

import (
	"github.com/konsultin/project-goes-here/internal/svc-core/constant"
)

func (r *Repository) PublishExampleEvent(message string) error {
	if r.nats == nil {
		return nil
	}
	// Publish to NATS
	return r.nats.Publish(constant.JobExample, []byte(message))
}
