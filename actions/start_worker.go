package actions

type StartWorker struct {
}

func writer(jobs <-chan Action, c Context) {
	for j := range jobs {
		err := j.Run(c)
		if err != nil {
			c.HandleErr(err)
		}
	}
}

func (a StartWorker) Run(c Context) error {
	jobs := c.BackgroundJobs()
	for w := 0; w < c.Concurrency(); w++ {
		go writer(jobs, c)
	}
	return nil
}

func (a StartWorker) MarshalJSON() ([]byte, error) {
	return []byte(`"Start workers"`), nil
}
