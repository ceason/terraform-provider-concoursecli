package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParsePipelineList(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given output from `fly pipelines --all`", t, func() {
		flyPipelinesOutput := `
hello-world              main  no      no
some-activated-pipeline  main  yes     no
`
		Convey("When parsing into pipeline objects", func() {
			parsedPipelines := parseFlyPipelinesOutput(flyPipelinesOutput)

			Convey("Should..work?", func() {
				expectedPipelines := []*pipeline{
					{name: "hello-world", team: "main", paused: false, public: false},
					{name: "some-activated-pipeline", team: "main", paused: true, public: false},
				}
				So(parsedPipelines, ShouldResemble, expectedPipelines)
			})
		})
	})
}
