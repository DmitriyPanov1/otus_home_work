package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	return pipeline(in, done, stages...)
}

func pipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return execute(in, done)
	}

	for _, stage := range stages {
		in = stage(execute(in, done))
	}

	return in
}

func execute(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}

				out <- v
			}
		}
	}()

	return out
}
