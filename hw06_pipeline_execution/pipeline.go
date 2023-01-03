package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := WithDoneStage(done, in)

	for _, stage := range stages {
		// Skip stage if it's nil
		if stage != nil {
			out = WithDoneStage(done, stage(out))
		}
	}
	return out
}

func WithDoneStage(done In, in In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for range in {
			}
		}()

		for {
			select {
			case <-done:
				return
			case res, ok := <-in:
				if !ok {
					return
				}

				select {
				case <-done:
				case out <- res:
				}
			}
		}
	}()

	return out
}
