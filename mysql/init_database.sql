drop table if exists runner cascade;
drop table if exists run cascade;
drop table if exists namespace cascade;
drop table if exists namespace_repos cascade;
drop table if exists batch cascade;
drop table if exists batch_runs cascade;

create table runner
(
	id					serial,
	status			varchar(16)		not null default 'waiting'
														check (status in ('waiting', 'running', 'dead')),
	last_conn		timestamp			default current_timestamp
														on update current_timestamp,
	ip          varchar(16)		not null,
	primary key (id)
);

create table run
(
	id					serial				not null,
	status			varchar(16)		not null default 'waiting'
														check (status in ('waiting', 'running', 'canceled'
																							'finished', 'failed', 'timeout')),
	runner			int						not null references runner(id),
	repo				varchar(64)		not null,
	logs				text					not null,
	primary key (id)
);

create table namespace
(
	name				varchar(64)		not null unique,
	primary key (name)
);

create table namespace_repos
(
	id					serial				not null,
	namespace   varchar(64)		not null references namespace(name),
	repo				varchar(64)		not null,
	primary key (id)
);

create table batch
(
	id					serial				not null,
	namespace		varchar(64)		not null references namespace(name),
	primary key (id)
);

create table batch_runs
(
	id					serial				not null,
	batch				int						not null references batch(id),
	run					int						not null references run(id),
	primary key (id)
);
