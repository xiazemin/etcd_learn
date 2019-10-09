package main

import (
	"github.com/coreos/etcd/clientv3"
	"time"
	"context"
	"fmt"
)

func main(){
	config := clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},//{"192.168.50.250:2379","172.16.196.129:2379"},
		DialTimeout:10*time.Second,
	}
	client,err := clientv3.New(config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	kv := clientv3.NewKV(client)
	//kv是一个用于操作kv的连接，其实它本质上是用了client的conn，为了更加专注于键值对的操作，关闭client后也会使kv无法用。（kv的操作client也能实现）

	//设置一个超时的context:

	ctx,cancleFunc:= context.WithTimeout(context.TODO(),5*time.Second)
	/**
	context.WithTimeout()会返回一个timerCtx{}，并在这个结构体里注入了超时时间。cancleFunc是一个取消操作的函数。put,get等操作是阻塞型操作，context里有一个用于管理超时的select,当时间一到就会隐式执行cancelFunc，使操作停止并返回错误。如果显式的调用cancelFunc()则会立即停止操作，返回错误。
	 */
	//put操作：
	putResp,err := kv.Put(ctx,"/job/v3","push the box",clientv3.WithPrevKV())  //withPrevKV()是为了获取操作前已经有的key-value
	if err != nil{
		panic(err)
	}
	fmt.Printf("%v",putResp.PrevKv)
	getResp,err := kv.Get(ctx,"/job/",clientv3.WithPrefix()) //withPrefix()是未了获取该key为前缀的所有key-value
	if err != nil{
		panic(err)
	}
	fmt.Printf("%v",getResp.Kvs)

	//由于etcd是有序存储键值对的，还可以附加clientv3.WithFromKey(),clientv3.WithLimit()来实现分页获取的效果。

	//监听etcd集群键的改变
	wc := client.Watch(context.Background(), "/job/", clientv3.WithPrefix(),clientv3.WithPrevKV())
	for v := range wc {
		if v.Err() != nil {
			panic(err)
		}
		for _, e := range v.Events {
			fmt.Printf("type:%v\n kv:%v  prevKey:%v  ", e.Type, e.Kv, e.PrevKv)
		}
	}

	//删除操作
	if delResp,err := kv.Delete(context.TODO(),"/cron/jobs/job2",clientv3.WithPrevKV()/*得到删除之前的值*/);err != nil {
		fmt.Println(err)
		return
	} else {
		if len(delResp.PrevKvs) != 0 {
			fmt.Println(delResp.PrevKvs)
		}
	}


	go func() {
		time.Sleep(20 * time.Second)
		cancleFunc() // 在调用处主动取消
	}()
}
