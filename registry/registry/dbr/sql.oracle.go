package dbr

var oracletexture sqltexture

func init() {

	//以下为oracle
	oracletexture.createStructure = `
	create table HYDRA_REGISTRY_INFO
	(
	  id           NUMBER(20) not null,
	  path         VARCHAR2(256) not null,
	  value        VARCHAR2(1024) not null,
	  is_temp      NUMBER(2) default 0 not null,
	  is_delete    NUMBER(2) default 1 not null,
	  data_version NUMBER(20),
	  create_time  DATE default sysdate not null,
	  update_time  DATE default sysdate not null
	)
	tablespace USERS
	  pctfree 10
	  initrans 1
	  maxtrans 255;
	comment on table HYDRA_REGISTRY_INFO
	  is '注册中心';
	comment on column HYDRA_REGISTRY_INFO.id
	  is '编号';
	comment on column HYDRA_REGISTRY_INFO.path
	  is '路径';
	comment on column HYDRA_REGISTRY_INFO.value
	  is '内容';
	comment on column HYDRA_REGISTRY_INFO.is_temp
	  is '临时节点(0:是，1否）';
	comment on column HYDRA_REGISTRY_INFO.is_delete
	  is '已删除';
	comment on column HYDRA_REGISTRY_INFO.data_version
	  is '数据版本号';
	comment on column HYDRA_REGISTRY_INFO.create_time
	  is '创建时间';
	comment on column HYDRA_REGISTRY_INFO.update_time
	  is '更新时间';
	create index IDX_HYDRA_REGISTRY_INFO_PATH on HYDRA_REGISTRY_INFO (PATH)
	  tablespace USERS
	  pctfree 10
	  initrans 2
	  maxtrans 255;
	alter table HYDRA_REGISTRY_INFO
	  add constraint PK_HYDRA_REGISTRY_INFO primary key (ID)
	  using index 
	  tablespace USERS
	  pctfree 10
	  initrans 2
	  maxtrans 255;

	  create sequence SEQ_HYDRA_REGISTRY_INFO_ID
	  minvalue 10000
	  maxvalue 99999999999999999999
	  start with 10000
	  increment by 1
	  cache 20
	  cycle;
	`

	oracletexture.createNode = `
	insert into hydra_registry_info(
		id,
		path,
		value,
		is_temp,
		is_delete,
		data_version,
		create_time,
		update_time 
	)
	select
		SEQ_HYDRA_REGISTRY_INFO_ID.Nextval,
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
	where t.path like '#path%'
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
	or (is_temp = 0 and update_time <= sysdate -30/24/60/60))
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
	`

	oracletexture.aclUpdate = `
	update hydra_registry_info t set
	t.update_time = sysdate
	where t.path in (#path)
	and t.is_delete = 1
	`

	oracletexture.clearTmpNode = `
	update hydra_registry_info t
	set t.update_time = sysdate,
	t.is_delete = 0
	where t.path in (#path)
	and t.is_delete = 1
	and is_temp = 0
	`

	oracletexture.getSeq = "select SEQ_HYDRA_REGISTRY_INFO_ID.Nextval from dual"
}
