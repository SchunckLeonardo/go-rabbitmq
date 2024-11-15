package events

import (
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

type TestEvent struct {
	Name    string
	Payload any
}

func (te *TestEvent) GetName() string {
	return te.Name
}

func (te *TestEvent) GetPayload() any {
	return te.Payload
}

func (te *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	ID int
}

func (th *TestEventHandler) Handle(event EventInterface, wg *sync.WaitGroup) {}

type EventDispatcherTestSuite struct {
	suite.Suite
	event           TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.event = TestEvent{
		Name:    "test",
		Payload: "test",
	}
	suite.event2 = TestEvent{
		Name:    "test2",
		Payload: "test2",
	}
	suite.handler = TestEventHandler{ID: 1}
	suite.handler2 = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
	suite.eventDispatcher = NewEventDispatcher()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	suite.Equal(suite.eventDispatcher.handlers[suite.event.GetName()][0], &suite.handler)
	suite.Equal(suite.eventDispatcher.handlers[suite.event.GetName()][1], &suite.handler2)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register_WithSameHandler() {
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.IsType(ErrHandlerAlreadyRegistered, err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	// Event 1
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	// Event 2
	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

	err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

	suite.eventDispatcher.Clear()
	suite.Equal(0, len(suite.eventDispatcher.handlers))
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	// Event 1
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	isHandlerExists := suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler)
	suite.True(isHandlerExists)

	isHandler2Exists := suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler2)
	suite.True(isHandler2Exists)

	isHandler3Exists := suite.eventDispatcher.Has(suite.event.GetName(), &suite.handler3)
	suite.False(isHandler3Exists)
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	// Event 1
	err := suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Register(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event.GetName()]))

	err = suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler)
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()], 1)

	err = suite.eventDispatcher.Remove(suite.event.GetName(), &suite.handler2)
	suite.Nil(err)
	suite.Len(suite.eventDispatcher.handlers[suite.event.GetName()], 0)
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event EventInterface, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	eventHandler := &MockHandler{}
	eventHandler.On("Handle", &suite.event)

	eventHandler2 := &MockHandler{}
	eventHandler2.On("Handle", &suite.event)

	err := suite.eventDispatcher.Register(suite.event.GetName(), eventHandler)
	suite.Nil(err)
	err = suite.eventDispatcher.Register(suite.event.GetName(), eventHandler2)
	suite.Nil(err)

	suite.eventDispatcher.Dispatch(&suite.event)

	eventHandler.AssertExpectations(suite.T())
	eventHandler2.AssertExpectations(suite.T())
	eventHandler.AssertNumberOfCalls(suite.T(), "Handle", 1)
	eventHandler2.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(EventDispatcherTestSuite))
}
