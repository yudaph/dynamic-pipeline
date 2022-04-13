package main

import (
	"fmt"
	_ "net/http/pprof"
	"pipelines/pipe"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var wg sync.WaitGroup
	baseChan := make(chan *pipe.CarBuilder, 1000)
	bodyChan := make(chan *pipe.CarBuilder, 1000)
	featureAChan := make(chan *pipe.CarBuilder, 1000)
	featureBChan := make(chan *pipe.CarBuilder, 1000)
	featureCChan := make(chan *pipe.CarBuilder, 1000)
	buildChan := make(chan *pipe.CarBuilder, 1000)

	typeStandard := []pipe.NextProcess{pipe.Process(buildChan)}
	typeA := []pipe.NextProcess{pipe.Process(featureAChan), pipe.Process(featureCChan), pipe.Process(buildChan)}
	typeB := []pipe.NextProcess{pipe.Process(featureAChan), pipe.Process(featureBChan), pipe.Process(buildChan)}
	typeC := []pipe.NextProcess{pipe.Process(featureBChan), pipe.Process(featureCChan), pipe.Process(buildChan)}
	typeFullFeature := []pipe.NextProcess{pipe.Process(featureAChan), pipe.Process(featureBChan), pipe.Process(featureCChan), pipe.Process(buildChan)}
	features := [][]pipe.NextProcess{
		typeStandard, typeA, typeB, typeC, typeFullFeature,
	}

	var carCount int64
	for i := 1; i <= 4; i++ {
		go pipe.BaseBuilder(baseChan, bodyChan)
		go pipe.BodyBuilder(bodyChan)
		go pipe.FeatureABuilder(featureAChan)
		go pipe.FeatureBBuilder(featureBChan)
		go pipe.FeatureCBuilder(featureCChan)
		go func() {
			for {
				readyToBuild := <-buildChan
				readyToBuild.Build()
				wg.Done()
				atomic.AddInt64(&carCount, 1)
			}
		}()
	}

	var testBuilds []pipe.CarBuilder
	for i := 0; i < 1000000; i++ {
		testBuilds = append(testBuilds, pipe.CarBuilder{
			Next: features[i%5],
		})
	}

	defer func() func() {
		start := time.Now()
		return func() { fmt.Println(time.Since(start), carCount) }
	}()()
	for _, testBuild := range testBuilds {
		wg.Add(1)
		t := testBuild
		baseChan <- &t
	}
	wg.Wait()
}
