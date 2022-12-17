# go-redis-lock

#### 介绍
基于go-redis的分布式锁，具备自动续期，可重入等能力

#### 使用说明
```go
result, err := LockTransaction("say_hello", reqID, func() error {
	...
	return nil
})
// 加锁失败
if !result {
        ...
	return
}
// 执行错误
if err != nil {
	...
	return
}
```
