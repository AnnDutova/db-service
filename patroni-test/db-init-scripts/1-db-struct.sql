Create user pguser with password 'pgpwd';

create database testdb;

GRANT CONNECT on database "testdb" to pguser;


create table if not exists "project" (
                                         id serial primary key,
                                         title text not null unique
);

create table if not exists "author" (
                                        id serial primary key,
                                        name text not null unique
);

create table if not exists "issues" (
                                        id serial primary key,
                                        projectId int not null,
                                        authorId int not null,
                                        assigneeId int,
                                        key text not null,
                                        createdTime timestamp not null,
                                        updatedTime timestamp not null,
                                        closedTime  timestamp,
                                        summary     text not null,
                                        description text,
                                        type        text not null,
                                        priority    text not null,
                                        status      text not null,
                                        timeSpent   bigint,
                                        info json,
                                        constraint "fk_issues_project" foreign key (projectId) references project(id) MATCH FULL,
    constraint "fk_issues_author" foreign key (authorId) references author(id) MATCH FULL,
    constraint "fk_issues_assignee" foreign key (assigneeId) references author(id) MATCH FULL
    );

create table if not exists "statusChange" (
                                              authorId int not null,
                                              issueId int not null,
                                              changeTime timestamp not null,
                                              fromStatus text,
                                              toStatus text not null,
                                              constraint "fk_statusChange_author" foreign key (authorId) references author(id) MATCH FULL
    );

create table if not exists "openTaskTime" (
                                              projectId int not null,
                                              createdTime timestamp not null,
                                              data json not null,
                                              constraint "fk_openTaskTime_project" foreign key (projectId) references project(id) MATCH FULL
    );

create table if not exists "taskStateTime" (
                                               projectId int not null,
                                               createdTime timestamp not null,
                                               data json not null,
                                               state text not null,
                                               constraint "fk_taskStateTime_project" foreign key (projectId) references project(id) MATCH FULL
    );

create table if not exists "complexityTaskTime" (
                                                    projectId int not null,
                                                    createdTime timestamp not null,
                                                    data json not null,
                                                    constraint "fk_complexityTaskTime_project" foreign key (projectId) references project(id) MATCH FULL
    );

create table if not exists "activityByTask" (
                                                projectId int not null,
                                                createdTime timestamp not null,
                                                data json not null,
                                                state text not null,
                                                constraint "fk_activityByTask_project" foreign key (projectId) references project(id) MATCH FULL
    );

create table if not exists "taskPriorityCount" (
                                                   projectId int not null,
                                                   createdTime timestamp not null,
                                                   data json not null,
                                                   state text not null,
                                                   constraint "fk_taskPriorityCount_project" foreign key (projectId) references project(id) MATCH FULL
    );

GRANT USAGE, SELECT ON SEQUENCE project_id_seq TO pguser;
GRANT USAGE, SELECT ON SEQUENCE author_id_seq TO pguser;
GRANT USAGE, SELECT ON SEQUENCE issues_id_seq TO pguser;

GRANT SELECT, INSERT, UPDATE, DELETE on table "project", "issues","author", "statusChange", "taskPriorityCount","activityByTask",
    "complexityTaskTime", "taskStateTime", "openTaskTime"  to pguser;

