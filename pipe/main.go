package pipe

var Process = func(channel chan *CarBuilder) func(builder *CarBuilder) {
	return func(builder *CarBuilder) {
		channel <- builder
	}
}

type car struct {
	base,
	body,
	featureA,
	featureB,
	featureC bool
}

func (c *car) setBase(base bool) {
	c.base = base
}

func (c *car) setBody(body bool) {
	c.body = body
}
func (c *car) setFeatureA(featureA bool) {
	c.featureA = featureA
}
func (c *car) setFeatureB(featureB bool) {
	c.featureB = featureB
}
func (c *car) setFeatureC(featureC bool) {
	c.featureC = featureC
}

type NextProcess func(builder *CarBuilder)

type CarBuilder struct {
	car
	Next []NextProcess
}

func (c *CarBuilder) nextProcess() {
	next := c.Next[0]
	c.Next = c.Next[1:]
	next(c)
}

func (c *CarBuilder) Build() car {
	return c.car
}

func BaseBuilder(channel <-chan *CarBuilder, nextChannel chan<- *CarBuilder) {
	for {
		carBuilder := <-channel
		carBuilder.setBase(true)
		nextChannel <- carBuilder
	}
}

func BodyBuilder(channel <-chan *CarBuilder) {
	for {
		carBuilder := <-channel
		carBuilder.setBody(true)
		carBuilder.nextProcess()
	}
}

func FeatureABuilder(channel <-chan *CarBuilder) {
	for {
		carBuilder := <-channel
		carBuilder.setFeatureA(true)
		carBuilder.nextProcess()
	}
}

func FeatureBBuilder(channel <-chan *CarBuilder) {
	for {
		carBuilder := <-channel
		carBuilder.setFeatureB(true)
		carBuilder.nextProcess()
	}
}

func FeatureCBuilder(channel <-chan *CarBuilder) {
	for {
		carBuilder := <-channel
		carBuilder.setFeatureC(true)
		carBuilder.nextProcess()
	}
}
