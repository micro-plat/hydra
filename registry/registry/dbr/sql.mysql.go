package dbr

type sqltexture struct {
	createStructure   string
	createNode        string
	exists            string
	getValue          string
	delete            string
	getChildren       string
	clear             string
	update            string
	getChildrenChange string
	getValueChange    string
	aclUpdate         string
	clearTmpNode      string
	getSeq            string
}

var mysqltexture sqltexture

func init() {

	mysqltexture.createStructure = `CREATE TABLE IF NOT EXISTS hydra_registry_info (
		id bigint not null auto_increment comment '编号' ,
		path varchar(256)  not null  comment '路径' ,
		value varchar(1024)  not null  comment '内容' ,
		is_temp tinyint default 0 not null  comment '临时节点' ,
		is_delete tinyint default 1 not null  comment '已删除' ,
		data_version bigint    comment '数据版本号' ,
		create_time datetime default current_timestamp not null  comment '创建时间' ,
		update_time datetime default current_timestamp not null  comment '更新时间' 
		,primary key (id)
		,unique index path(path)
	) ENGINE=InnoDB auto_increment = 100 DEFAULT CHARSET=utf8mb4 COMMENT='注册中心'`

	mysqltexture.createNode = `
	insert into hydra_registry_info(
		path,
		value,
		is_temp,
		is_delete,
		data_version,
		create_time,
		update_time 
	)
	select
		@path,
		@value,
		@is_temp,
		1,
		@data_version,
		now(),
		now()
	from dual
	where not exists(select 1 from hydra_registry_info t where t.path = @path)
	`
	mysqltexture.exists = `
	select count(0)
	from hydra_registry_info t
	where t.path like '#path%'
	and t.is_delete = 1
	`

	mysqltexture.getValue = `
	select
	t.id,
	t.path,
	t.value,
	t.is_temp,
	t.is_delete,
	t.data_version,
	t.create_time,
	t.update_time 
	from hydra_registry_info t
	where t.path = @path
	and t.is_delete = 1
	`

	mysqltexture.delete = `
	update hydra_registry_info t
	set t.update_time = now(),
	t.is_delete = 0
	where t.path = @path
	and t.is_delete = 1
	`
	mysqltexture.clear = `
	delete from hydra_registry_info
	where path = @path
	and (is_delete = 0
	or (is_temp = 0 and update_time <= date_add(now(),interval -30 second)))
	`
	mysqltexture.update = `
	update hydra_registry_info t set
	t.value = @value,
	t.data_version = t.data_version + 1,
	t.update_time = now() 
	where t.path = @path 
	and t.data_version = @data_version
	and t.is_delete = 1
	`

	mysqltexture.getChildren = `
	select
	t.path,
	t.value,
	t.is_temp,
	t.is_delete,
	t.data_version,
	t.create_time,
	t.update_time 
	from hydra_registry_info t
	where t.path like '#path%'
	and t.is_delete = 1
	`
	mysqltexture.getValueChange = `
	select
	t.path,
	t.value,
	t.is_temp,
	t.is_delete,
	t.data_version,
	t.create_time,
	t.update_time 
	from hydra_registry_info t
	where t.path in (#path) 
	and t.update_time > date_add(now(),interval -1*#sec second)
	`
	mysqltexture.getChildrenChange = `
	select
	t.path,
	t.value,
	t.is_temp,
	t.is_delete,
	t.data_version,
	t.create_time,
	t.update_time 
	from hydra_registry_info t
	where t.path like '#path%'
	and (t.create_time > date_add(now(),interval -1*#sec second) 
		or(t.is_delete = 0 and t.update_time > date_add(now(),interval -1*#sec second)))
	`

	mysqltexture.aclUpdate = `
	update hydra_registry_info t set
	t.update_time = now() 
	where t.path in (#path)
	and t.is_delete = 1
	`

	mysqltexture.clearTmpNode = `
	update hydra_registry_info t
	set t.update_time = now(),
	t.is_delete = 0
	where t.path in (#path)
	and t.is_delete = 1
	and is_temp = 0
	`

	mysqltexture.getSeq = "select replace(unix_timestamp(current_timestamp(3)),'.','')"
}
