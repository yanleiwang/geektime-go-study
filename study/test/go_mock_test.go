package test

import (
	"errors"
	"geektime-go-study/study/test/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
教程:   https://geektutu.com/post/quick-gomock.html
生成mock文件命令:   mockgen -source=db.go -destination=mocks/mock_db.gen.go -package=mocks

*/

func TestGetFromDB(t *testing.T) {

	// ctrl 是mock 的顶层控制器, 控制mock对象的生命周期,
	// 每个测试cases 都应该首先创建一个ctrl,
	//  并调用 defer ctrl.Finish()  结束
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock 对象
	m := mocks.NewMockDB(ctrl)
	// 打桩
	// 当 Get() 的参数为 Tom，则返回 error
	// 更多打桩api 参见参考
	m.EXPECT().Get(gomock.Eq("Tom")).Return(100, errors.New("not exist"))

	v := GetFromDB(m, "Tom")
	assert.Equal(t, -1, v)
}
