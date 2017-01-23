set foreign_key_checks = 0;
drop table if exists runner cascade;
drop table if exists run cascade;
drop table if exists namespace cascade;
drop table if exists namespace_repos cascade;
drop table if exists batch cascade;
set foreign_key_checks = 1;

create table runner
(
	id		serial,
	name            varchar(32)	not null default 'no_name',
	status		varchar(16)	not null default 'waiting'
					check (status in ('waiting', 'running', 'paused', 'dead')),
	last_conn	timestamp	default current_timestamp
					on update current_timestamp,
	ip		varchar(64)	not null,
	primary key (id)
);

create table namespace
(
	name	varchar(64)	not null unique,
        is_ci   boolean,
	primary key (name)
);

create table namespace_repos
(
	id		serial		not null,
	namespace	varchar(64)	not null,
	repo		varchar(64)	not null,
	ssh             varchar(128)	not null,
        active          boolean         default true,
	primary key (id),
	foreign key(namespace) references namespace(name)
);

create table batch
(
	id		serial		not null,
	namespace	varchar(64)	not null,
	init_script	mediumtext      not null,
	update_time	int		not null,
	timeout		int		not null,
	creation        timestamp	default current_timestamp,
	primary key (id),
	foreign key(namespace) references namespace(name)
);

create table run
(
	id	    serial	    not null,
	status	    varchar(16)	    not null default 'waiting'
                                    check (status in ('waiting', 'running'
                                    'canceled', 'finished', 'failed', 'timeout')),
	runner	    bigint unsigned not null,
	repo        bigint unsigned not null,
	batch	    bigint unsigned not null,
	logs	    mediumtext	    not null,
	creation    timestamp	    default current_timestamp,
	last_update timestamp	    default current_timestamp
                                    on update current_timestamp,
        start_time  timestamp       default current_timestamp,
        priority    int unsigned    default 0,

	primary key (id),
	foreign key(repo) references namespace_repos(id),
	foreign key(batch) references batch(id)
);
