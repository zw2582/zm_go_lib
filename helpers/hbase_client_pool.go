package helpers

import (
	"context"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/astaxie/beego"
	"github.com/jolestar/go-commons-pool"
	"github.com/pkg/errors"
	"net"
	"time"
	"zm_go_lib/libs/hbase"
)

var (
	hbase_host, hbase_port, hbase_list_table string
	hbase_p *pool.ObjectPool
)

//RegistHbaseClientPool 注册hbase连接池
func RegistHbaseClientPool() {
	hbase_host = beego.AppConfig.String("hbase_host")
	hbase_port = beego.AppConfig.String("hbase_port")
	hbase_list_table = beego.AppConfig.String("hbase_list_table")

	if hbase_host == "" || hbase_port == "" || hbase_list_table == "" {
		panic(errors.New("请在conf/app.conf中hbase参数：hbase_host，hbase_port，hbase_list_table"))
	}

	hbase_p = createHbaseClietPool()
}

//createHbaseClient：定义创建hbaseClient的函数
func createHbaseClient() *hbase.HbaseClient {
	//创建socket
	socket, err := thrift.NewTSocket(net.JoinHostPort(hbase_host, hbase_port))
	if err != nil {
		panic(err)
	}
	socket.SetTimeout(time.Second * 10)
	//创建tran
	transport := thrift.NewTBufferedTransport(socket, 1024)
	protocol := thrift.NewTBinaryProtocolFactory(false, true)
	client := hbase.NewHbaseClientFactory(transport, protocol)

	//if err := transport.Open(); err != nil {
	//	panic(err)
	//}
	//defer transport.Close()
	return client
}

//定义自己的hbase连接池工厂
type myHbaseClientFactory struct {
}

func (f *myHbaseClientFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	beego.Debug(`创建hbase连接client`)
	return pool.NewPooledObject(createHbaseClient()), nil
}

func (f *myHbaseClientFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	beego.Debug(`校验hbaseclient`)
	// do validate
	client := object.Object.(hbase.HbaseClient)
	enabled, err := client.IsTableEnabled(hbase.Bytes(hbase_list_table))
	if err != nil {
		return false
	}
	return enabled
}

// DestroyObject 关闭hbaseclient
func (f *myHbaseClientFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	client := object.Object.(*hbase.HbaseClient)

	if client.Transport.IsOpen() {
		beego.Debug("连接已打开，准备关闭")
		if err := client.Transport.Close(); err != nil {
			beego.Error(err)
			return err
		}
		if !client.Transport.IsOpen() {
			beego.Debug("连接已关闭")
		} else {
			beego.Debug("连接依然打开，没有关闭")
		}
	}
	return nil
}

//ActivateObject 激活hbaseclient
func (f *myHbaseClientFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	// do activate
	client := object.Object.(*hbase.HbaseClient)

	if !client.Transport.IsOpen() {
		beego.Debug("激活hbaseclient:连接未打开，打开连接")
		if err := client.Transport.Open(); err != nil {
			beego.Error(err)
		}
	}

	return nil
}

//PassivateObject 钝化对象
func (f *myHbaseClientFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	client := object.Object.(*hbase.HbaseClient)

	if client.Transport.IsOpen() {
		beego.Debug("连接已打开，准备关闭")
		if err := client.Transport.Close(); err != nil {
			beego.Error(err)
			return err
		}
	}
	return nil
}

//CreateHbaseClietPool:创建连接池,最好保证单例
func createHbaseClietPool() *pool.ObjectPool {
	beego.Debug(`创建hbaseclient连接池`)
	ctx := context.Background()
	p := pool.NewObjectPoolWithDefaultConfig(ctx, &myHbaseClientFactory{})
	p.Config.MaxTotal = 50
	return p
}

//GetClient：获取连接对象
func HbaseClient() *hbase.HbaseClient {
	if hbase_p == nil {
		panic(errors.New("请先注册hbase连接池：RegistHbaseClientPool"))
	}
	ctx := context.Background()
	obj, err := hbase_p.BorrowObject(ctx)
	if err != nil {
		panic(err)
	}
	client := obj.(*hbase.HbaseClient)
	return client
}

//Close:回收hbaseclient
func HbaseClose(client *hbase.HbaseClient) {
	ctx := context.Background()
	if err := hbase_p.ReturnObject(ctx, client); err != nil {
		panic(err)
	}
}