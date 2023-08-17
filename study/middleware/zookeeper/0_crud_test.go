package zookeeper

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	hosts = []string{"127.0.0.1:2181", "127.0.0.1:2182", "127.0.0.1:2183"}
)

func createConn() (*zk.Conn, error) {
	// 连接zk
	conn, _, err := zk.Connect(hosts, time.Second*5)
	return conn, err
}

func TestCRUD(t *testing.T) {
	conn, err := createConn()
	assert.NoError(t, err)
	defer conn.Close()

	assert.NoError(t, add(conn))
	assert.NoError(t, get(conn))
	assert.NoError(t, modify(conn))
	assert.NoError(t, del(conn))

}

const path = "/test"

// 增
func add(conn *zk.Conn) error {
	var data = []byte("test value")
	// flags有4种取值：
	// 0:永久，除非手动删除
	// zk.FlagEphemeral = 1:短暂，session断开则该节点也被删除
	// zk.FlagSequence  = 2:会自动在节点后面添加序号
	// 3:Ephemeral和Sequence，即，短暂且自动添加序号
	var flags int32 = 0
	// 节点访问控制权限, 谁都访问
	acls := zk.WorldACL(zk.PermAll)
	s, err := conn.Create(path, data, flags, acls)
	if err != nil {
		return err
	}
	fmt.Printf("创建: %s 成功\n", s)
	return nil
}

// 查
func get(conn *zk.Conn) error {
	data, _, err := conn.Get(path)
	if err != nil {
		return err
	}
	fmt.Printf("%s 的值为 %s\n", path, string(data))
	return nil
}

// 删改 需要带上version, 作为乐观锁
// 1） version > 0, 会进行版本号检查 只有版本号一致的情况下，
//
//	才会进行写操作， 避免数据并发修改的一致性 不然会返回 zk.ErrBadVersion
//
// 2）  version == -1 忽略版本号检查
// 改
func modify(conn *zk.Conn) error {
	new_data := []byte("hello zookeeper")
	_, sate, _ := conn.Get(path)
	_, err := conn.Set(path, new_data, sate.Version)
	if err != nil {
		return err
	}
	return nil
}

// 删
func del(conn *zk.Conn) error {
	_, sate, _ := conn.Get(path)
	err := conn.Delete(path, sate.Version)
	if err != nil {
		return err
	}
	return nil
}
