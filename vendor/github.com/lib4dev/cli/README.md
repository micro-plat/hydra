# cli
命令行程序基础框架


#### 一、创建cli app
```go
func main(){
    var app=cli.New(cli.WithVersion("0.1.0"))
    app.Start()
}


```

#### 二、添加处理命令

```go

//Action .
func list(c *cli.Context) (err error) {
	//do something...
	return nil
}

func init() {
	cmds.Register(
		cli.Command{
			Name:   "ls",
			Usage:  "列出所有子节点",
			Action: list,
		})
}
```
