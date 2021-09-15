package dbr

const createNode = `
insert into 
hydra_registry_info
(
	path,
	value
	temp,
	data_version,
	acl_version,
)
values
(
	@path,
	@value
	@data_version,
	@acl_version
)
`

const exists = `
select count(1)
from  hydra_registry_info t
where t.path = @path'
`

const getValue = `
select 
t.path,
t.data_version
from hydra_registry_info t
where t.path = @path
`

const delete = `
update hydra_registry_info t
set t.delete_time = now(),
t.is_delete = 0
where t.path=@path
`

const update = `
update hydra_registry_info t
set
t.value = @value,
t.data_version = t.data_version + 1,
t.acl_version = t.acl_version + 1
t.update_time = now()
where
t.path=@path and t.data_version=@data_version
`

const getChildren = `
select t.path,t.value,t.data_version
from hydra_registry_info t
where
t.path like '#path%'
`
const getValueChange = `
select t.path,t.value,t.data_version
from hydra_registry_info t
where
t.path in(#path) and t.update_time > date_add(now(),interval -1*#sec second)
`
const getChildrenChange = `
select t.path,t.value,t.data_version
from hydra_registry_info t
where
t.path like '#path%'
and (t.create_time > date_add(now(),interval -1*#sec second) or t.delete_time > date_add(now(),interval -1*#sec second))
and limit 1
`
const aclUpdate = `
update hydra_registry_info t
set
t.update_time = now(),
t.acl_version = t.acl_version + 1
where t.path in(#path)
`
