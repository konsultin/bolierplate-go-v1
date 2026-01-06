package service

func (s *Service) RunSimulation() error {
	s.log.Info("Running simulation: Server -> Service -> Repo -> NATS -> Worker")

	// Call Repo to publish
	err := s.repo.PublishExampleEvent("Hello from Repository!")
	if err != nil {
		return err
	}

	return nil
}
