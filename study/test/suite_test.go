package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ExampleTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *ExampleTestSuite) SetupTest() {
	suite.VariableThatShouldStartAtFive = 5
}

func (suite *ExampleTestSuite) TestExample() {
	assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}

// 使用testify/suite 我们需要定义一个结构体, 然后内嵌一个匿名的suite.Suite结构。
// 然后每个以 Test开头的方法都会被视为一个测试用例
type MyTestSuit struct {
	suite.Suite
	testCount                     uint32
	VariableThatShouldStartAtFive int
}

// 所有测试用例执行前 会执行这个函数
func (s *MyTestSuit) SetupSuite() {
	fmt.Println("SetupSuite")
	s.VariableThatShouldStartAtFive = 5
}

// 所有测试用例执行后 会执行这个函数
func (s *MyTestSuit) TearDownSuite() {
	fmt.Println("TearDownSuite")
}

// 每个test用例执行前都会执行这个函数
func (s *MyTestSuit) SetupTest() {
	fmt.Printf("SetupTest test count:%d\n", s.testCount)
}

// 每个test用例执行后都会执行这个函数
func (s *MyTestSuit) TearDownTest() {
	s.testCount++
	fmt.Printf("TearDownTest test count:%d\n", s.testCount)
}

// 每个test用例执行前还会执行这个函数, 接受套件名和测试名作为参数。 先执行SetupTest 再执行BeforeTest
func (s *MyTestSuit) BeforeTest(suiteName, testName string) {
	fmt.Printf("BeforeTest suite:%s test:%s\n", suiteName, testName)
}

// 每个test用例执行后还会执行这个函数, 接受套件名和测试名作为参数。
func (s *MyTestSuit) AfterTest(suiteName, testName string) {
	fmt.Printf("AfterTest suite:%s test:%s\n", suiteName, testName)
}

// 其中一个测试用例
// 每个以 Test开头的方法都会被视为一个测试用例
// suite.Run(t, new(ExampleTestSuite)) 会执行所有执行用例
func (s *MyTestSuit) TestExampleV1() {
	fmt.Println("TestExampleV1")
	assert.Equal(s.T(), 5, s.VariableThatShouldStartAtFive)
	s.VariableThatShouldStartAtFive += 1
}

func (s *MyTestSuit) TestExampleV2() {
	fmt.Println("TestExampleV2")
}

// 运行测试套件的方法
// 我们需要借助go test运行，所以需要编写一个TestXxx函数，在该函数中调用suite.Run()运行测试套件：
func TestExample(t *testing.T) {
	suite.Run(t, new(MyTestSuit))
}
