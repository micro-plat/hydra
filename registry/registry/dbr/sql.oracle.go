package dbr

var oracletexture sqltexture

func init() {

	//以下为oracle
	oracletexture.createStructure = `CREATE TABLE IF NOT EXISTS hydra_registry_info (
		id bigint not null auto_increment comment '编号' ,
		path varchar(64)  not null  comment '路径' ,
		value varchar(1024)  not null  comment '内容' ,
		is_temp tinyint default 0 not null  comment '临时节点' ,
		is_delete tinyint default 1 not null  comment '已删除' ,
		data_version bigint    comment '数据版本号' ,
		create_time datetime default current_timestamp not null  comment '创建时间' ,
		update_time datetime default current_timestamp not null  comment '更新时间' 
		,primary key (id)
		,unique index path(path)
	) ENGINE=InnoDB auto_increment = 100 DEFAULT CHARSET=utf8mb4 COMMENT='注册中心'`

	oracletexture.createNode = `
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
		sysdate,
		sysdate
	from dual
	where not exists(select 1 from hydra_registry_info t where t.path = @path)
	`
	oracletexture.exists = `
	select count(0)
	from hydra_registry_info t
	where t.path = @path
	and t.is_delete = 1
	`

	oracletexture.getValue = `
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

	oracletexture.delete = `
	update hydra_registry_info t
	set t.update_time = sysdate,
	t.is_delete = 0
	where t.path = @path
	and t.is_delete = 1
	`
	oracletexture.clear = `
	delete from hydra_registry_info
	where path = @path
	and (is_delete = 0
	or is_temp = 0 and update_time > sysdate -30/24/60/60)
	`
	oracletexture.update = `
	update hydra_registry_info t set
	t.value = @value,
	t.data_version = t.data_version + 1,
	t.update_time = sysdate
	where t.path = @path 
	and t.data_version = @data_version
	and t.is_delete = 1
	`

	oracletexture.getChildren = `
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
	oracletexture.getValueChange = `
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
	and t.update_time > sysdate -1*#sec/24/60/60 
	`
	oracletexture.getChildrenChange = `
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
	and (t.create_time > sysdate -1*#sec/24/60/60 
		or(t.is_delete = 0 and t.update_time > sysdate -1*#sec/24/60/60 ))
	and limit 1
	`

	oracletexture.aclUpdate = `
	update hydra_registry_info t set
	t.update_time = sysdate
	where t.path in (#path)
	and t.is_delete = 1
	`

	oracletexture.clearTmpNode = `
	delete from hydra_registry_info
	where path in (#path)
	and is_temp = 0
	`
}
