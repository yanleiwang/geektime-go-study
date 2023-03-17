package builder

import "testing"

/*
builder模式 用于分步骤去构建一个复杂对象.
如果一个类的构造函数有很多可选参数, 那就适合用builder模式

在Builder模式中，通常会定义一个Builder接口，该接口定义了对象的构造过程中所需的各个步骤(可选参数)，
例如设置属性、添加部件等。
然后，在具体的Builder实现中，实现这些接口方法，最后实现一个Build方法, 将各个步骤组合起来, 最后完成对象的构造。

在下面的例子中，我们首先定义了一个User结构体，该结构体表示一个用户对象。
然后，我们定义了一个UserBuilder接口，该接口定义了对象构造的各个步骤，
并定义Build方法，该方法用于将各个步骤组合起来，生成最终的对象。
接着，我们定义了UserBuilder接口的实现类 userBuilderImpl，
最后，在客户端代码中，我们可以使用NewUserBuilder方法创建一个UserBuilder实例

为什么不直接在user上实现setname方法?
1. 这样就会生成一个处于中间状态的对象
2. 不同字段之间可能存在依赖关系. 直接在User结构体中实现SetXXX方法可能会使得代码逻辑变得复杂和难以维护。
*/

type User struct {
	name    string
	age     int
	address string
}

type UserBuilder interface {
	SetName(name string) UserBuilder
	SetAge(age int) UserBuilder
	SetAddress(address string) UserBuilder
	Build() *User
}

type userBuilderImpl struct {
	name    string
	age     int
	address string
}

func NewUserBuilder() UserBuilder {
	return &userBuilderImpl{}
}

func (b *userBuilderImpl) SetName(name string) UserBuilder {
	b.name = name
	return b
}

func (b *userBuilderImpl) SetAge(age int) UserBuilder {
	b.age = age
	return b
}

func (b *userBuilderImpl) SetAddress(address string) UserBuilder {
	b.address = address
	return b
}

func (b *userBuilderImpl) Build() *User {
	return &User{
		name:    b.name,
		age:     b.age,
		address: b.address,
	}
}

// 使用方法:
// 我们首先使用NewUserBuilder方法创建了一个UserBuilder实例。
// 然后，我们使用SetXXX方法来设置用户对象的各个属性，最后使用Build方法来生成最终的用户对象。
// 通过这种方式，我们可以方便地构造不同属性的用户对象，而不需要为每个属性都定义一个构造函数。
// 同时，我们也可以在Builder的实现中添加更多的构造步骤，从而使得对象构造的过程更加灵活和可扩展。
func Test_main(t *testing.T) {
	builder := NewUserBuilder()

	// 构造一个年龄为18岁的用户
	user1 := builder.SetName("Tom").SetAge(18).Build()
	println(user1)
	// 构造一个地址为Beijing的用户
	user2 := builder.SetName("Jerry").SetAddress("Beijing").Build()
	println(user2)
}
