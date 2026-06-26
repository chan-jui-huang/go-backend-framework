package job

import "fmt"

type ExampleJob struct {
	cronExpression string
	name           string
}

func NewExampleJob() *ExampleJob {
	return &ExampleJob{
		cronExpression: "* * * * * *",
		name:           "example",
	}
}

func (job *ExampleJob) Name() string {
	return job.name
}

func (job *ExampleJob) CronExpression() string {
	return job.cronExpression
}

func (job *ExampleJob) Execute() {
	fmt.Printf("The %s is finished\n", job.name)
}
