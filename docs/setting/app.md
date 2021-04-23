
### 1. 服务器参数

指平台、系统、
```go
	app := hydra.NewApp(
            hydra.WithPlatName("test"),
            hydra.WithSystemName("apiserver"),
            hydra.WithServerTypes(http.API),
            hydra.WithClusterName("prod"),			
    )
    app.Start()

```