Go-Promise

```go
func GetUserName(id time.Duration) *eventloop.Promise {
	return GlobalEventLoop.Async(func() (interface{}, error) {
		<-time.After(time.Second * id)
		if id == 0 {
			return nil, fmt.Errorf("some error id(%s)", id)
		}
		return fmt.Sprintf("id(%s): Test User", id), nil
	})
}

func GetUserNameWithPanic() *eventloop.Promise {
	return GlobalEventLoop.Async(func() (interface{}, error) {
		<-time.After(time.Second * 2)
		panic("panic attack")
	})
}

func main() {
	GlobalEventLoop.Main(func() {

		f := GlobalEventLoop.NewFuture()

		f.RegisterComplete(func(data interface{}) {
			fmt.Println("Received data from future --> ", data.(string))
		})

		result := GetUserName(2)

		result.Then(func(x interface{}) {
			fmt.Println("2 : user:", x)
		})

		fmt.Println("run before promise returns")

		GetUserName(0).Then(func(x interface{}) {
			fmt.Println("0 : user:", x)
		}).Catch(func(err error) {
			fmt.Println("0 : err:", err)
		})

		go func() {
			<-time.After(4 * time.Second)
			f.SignalComplete("completed after 4 seconds")
		}()

		f.SignalComplete("another hello world")

		GetUserName(5).Then(func(x interface{}) {
			fmt.Println("5 : user:", x)
			panic("a panic attack")
		}).Catch(func(err error) {
			fmt.Println("5 : err:", err)
		})

		GetUserName(15).Then(func(x interface{}) {
			fmt.Println("15 : user:", x)
		}).Catch(func(err error) {
			fmt.Println("15 : err:", err)
		})

		go func() {
			<-time.After(2 * time.Second)
			f.SignalComplete("completed after 2 seconds")
		}()

		//promise with panic
		GetUserNameWithPanic().Then(func(i interface{}) {
			fmt.Println("Then block never gets triggered")
		}).Catch(func(err error) {
			fmt.Println("GetUserNameWithPanic err: ", err)
		})

		//	await
		syncResult1, err := GlobalEventLoop.Await(GetUserName(4))
		fmt.Printf("syncResult1 - value: %v, err: %v\n", syncResult1, err)

		syncResult2, err := GlobalEventLoop.Await(GetUserName(1))
		fmt.Printf("syncResult2 - value: %v, err: %v\n", syncResult2, err)

		asyncResult := GetUserName(0)
		GetUserName(3)

		syncResult, err := GlobalEventLoop.Await(asyncResult)
		fmt.Printf("asyncResult - value: %v, err: %v\n", syncResult, err)

		fmt.Println("done")

		//nested promise
		GlobalEventLoop.Async(func() (interface{}, error) {
			fmt.Println("outer async")
			GlobalEventLoop.Async(func() (interface{}, error) {
				fmt.Println("inner async")
				return nil, nil
			}).Then(func(_ interface{}) {
				fmt.Println("resolved inner promise")
			})
			<-time.After(time.Second * 2)
			return nil, nil
		}).Then(func(_ interface{}) {
			fmt.Println("resolved outer promise")
		})
	})
}
```

#### result

```shell
run before promise returns
0 : err: some error id(0s)
Received data from future -->  another hello world
Received data from future -->  completed after 2 seconds
GetUserNameWithPanic err:  panic attack
2 : user: id(2ns): Test User
syncResult1 - value: id(4ns): Test User, err: <nil>
Received data from future -->  completed after 4 seconds
5 : user: id(5ns): Test User
5 : err: a panic attack
syncResult2 - value: id(1ns): Test User, err: <nil>
asyncResult - value: <nil>, err: some error id(0s)
done
outer async
inner async
resolved inner promise
resolved outer promise
15 : user: id(15ns): Test User
```

## TODO

- [x] nested promises
- [ ] chained .then

```go
GetUserName(7).Then(func(x interface{}) {
	fmt.Println("7 (1): user:", x)
}).Then(func(x interface{}) {
	fmt.Println("7 (2) : user:", x)
	panic("another panic attack")
}).Catch(func(err error) {
	fmt.Println("7 : err:", err)
})
```

- [ ] await all promises
- [x] basic Futures implementation
- [ ] handle error in Futures
- [ ] SignalFinally in Futures (called immediately after SignalComplete or SignalError)

> just a fun project, we might just learn something
