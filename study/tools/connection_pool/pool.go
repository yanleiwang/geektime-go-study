package connection_pool

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

// 空闲连接
type idleConn struct {
	conn           net.Conn
	lastActiveTime time.Time
}

// 等待的请求
type connReq struct {
	connChan chan net.Conn
}

/*
Pool 连接池 采用懒惰删除策略
*/
type Pool struct {
	maxCnt      int                      // 最大连接数
	cnt         int                      // 当前连接数
	maxIdleCnt  int                      //最大空闲连接
	maxIdleTime time.Duration            // 最大空闲时间
	idleConns   chan *idleConn           // 空闲连接队列
	reqQueue    []connReq                // 请求队列
	factory     func() (net.Conn, error) // 创建连接函数
	close       func(net.Conn) error     // 关闭函数
	lock        sync.Mutex
}

func NewPool(initCnt int, maxIdleCnt int, maxCnt int, maxIdleTime time.Duration, factory func() (net.Conn, error)) (*Pool, error) {
	if initCnt > maxIdleCnt {
		return nil, errors.New("micro: 初始连接数量不能大于最大空闲连接数量")
	}

	idleConns := make(chan *idleConn, maxIdleCnt)
	for i := 0; i < initCnt; i++ {
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		idleConns <- &idleConn{conn: conn, lastActiveTime: time.Now()}
	}

	c := func(conn net.Conn) error {
		return conn.Close()
	}

	res := &Pool{
		idleConns:   idleConns,
		maxCnt:      maxCnt,
		cnt:         initCnt,
		maxIdleTime: maxIdleTime,
		factory:     factory,
		close:       c,
	}
	return res, nil
}

func (p *Pool) Get(ctx context.Context) (net.Conn, error) {
	for {
		select {
		case <-ctx.Done(): // 超时
			return nil, ctx.Err()
		case c := <-p.idleConns: // 拿到了空闲连接
			// 还没有过期
			if c.lastActiveTime.Add(p.maxIdleTime).Before(time.Now()) {
				_ = p.close(c.conn)
				continue
			}
			return c.conn, nil
		default: // 没有空闲连接
			// 看看还能不能 创建新的连接
			p.lock.Lock()
			if p.cnt < p.maxCnt {
				conn, err := p.factory()
				if err != nil {
					p.lock.Unlock()
					return nil, err
				}
				p.cnt++
				p.lock.Unlock()
				return conn, nil
			}
			// 创建不了, 就加入等待队列
			req := connReq{connChan: make(chan net.Conn, 1)}
			p.reqQueue = append(p.reqQueue, req)
			p.lock.Unlock()
			select {
			// 超时了
			case <-ctx.Done():
				// 选项1：从队列里面删掉 req 自己, 但还是有并发问题
				// 选项2：在这里转发 放回去
				go func() {
					c := <-req.connChan
					_ = p.Put(context.Background(), c)
				}()
				return nil, ctx.Err()
			//  从别人那拿到
			case c := <-req.connChan:
				return c, nil
			}

		}

	}

}

func (p *Pool) Put(ctx context.Context, c net.Conn) error {

	p.lock.Lock()
	defer p.lock.Unlock()
	// 先看看有没有在等着的
	if len(p.reqQueue) > 0 {
		p.reqQueue[0].connChan <- c
		p.reqQueue = p.reqQueue[1:]
		return nil
	}

	select {
	// 可以放回idleConns
	case p.idleConns <- &idleConn{
		conn:           c,
		lastActiveTime: time.Now(),
	}:
	// 放不进去了, close掉
	default:
		_ = p.close(c)
		p.cnt--
	}
	return nil
}
