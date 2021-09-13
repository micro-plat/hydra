package sql

const CreatePersistentNode = `
insert into 
sys_registry_info
(
	path,
	value
)
values
(
	@path,
	@value
)
`

const Exists = `
select count(1)
from 
sys_registry_info t
where
t.path like '#path%'
`

const GetValue = `
select 
t.path,
t.value
from sys_registry_info t
where
t.path=@path
`

const Delete = `
delete from sys_registry_info
where
path=@path
`

const Update = `
update
sys_registry_info t
set
t.value=@value,
t.last_update_time=now()
where
t.path=@path
`

const GetChildren = `
select 
t.path,
t.value
from sys_registry_info t
where
t.path like '#path%'
`

const GetData = `
select *
from 
sys_registry_info t
`
