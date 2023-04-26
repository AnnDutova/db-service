create or replace procedure insertIssue(issue_id int, project_title text, project_id int, author_name text, assigned_name text,
                                        key_ text, createTime_ timestamp, closedTime_ timestamp, updatedTime_ timestamp, summary_ text, description_ text, type_ text,
                                        priority_ text, status_  text, timeSpent_ bigint)
    LANGUAGE plpgsql
as $$
Declare author_id int;
    Declare assignee_id int;
begin
    if exists (select id from author where author.name=author_name) then
Select id into author_id from author where author.name=author_name;
else
select insertAuthor(author_name) into author_id;
end if;

    if exists (select id from author where author.name=assigned_name) then
Select id into assignee_id from author where author.name=assigned_name;
else
        if length(assigned_name) > 0 then
select insertAuthor(assigned_name) into assignee_id;
end if;
end if;

    if not exists (select id from project where project.id=project_id) then
        perform insertProject(project_title, project_id);
end if;

Insert Into issues(id, projectId, authorId, assigneeId, key, createdTime,closedTime, updatedTime, summary,
                   description, type, priority, status,timeSpent) values
    (issue_id, project_id, author_id,assignee_id, key_,createTime_,closedTime_,updatedTime_, summary_, description_,
     type_,priority_,status_, timeSpent_);
commit;
end;
$$;

create or replace procedure insertStatusChange(author_name text, issue_id int,
                                               changeTime_ timestamp, fromStatus_ text, toStatus_ text)
    LANGUAGE plpgsql
as $$
Declare author_id int;
begin
    if exists (select id from author where author.name=author_name) then
Select id into author_id from author where author.name=author_name;
else
select insertAuthor(author_name) into author_id;
end if;
Insert Into "statusChange"(issueId, authorId, changeTime, fromStatus, toStatus)
Values(issue_id, author_id, changeTime_,fromStatus_, toStatus_);
end;
$$;

create or replace procedure updateIssue(issue_id int, assignee_name text, new_updatedTime_ timestamp, closedTime_ timestamp,  summary_ text,
                                        description_ text, type_ text, priority_ text, status_  text, timeSpent_ bigint)
    LANGUAGE plpgsql
as $$
Declare assignee_id int;
    Declare old_date timestamp;
begin
    if exists (select id from author where author.name=assignee_name) then
Select id into assignee_id from author where author.name=assignee_name;
else
select insertAuthor(assignee_name) into assignee_id;
end if;

select updatedTime into old_date from issues where id = issue_id;

if (select updatedTime from issues where id = issue_id) < new_updatedTime_ then
Update issues set (authorId, updatedTime, closedTime, summary, description, type, priority, status, timeSpent) =
                      (assignee_id, new_updatedTime_, closedTime_, summary_, description_, type_, priority_, status_, timeSpent_)
where  id = issue_id;
end if;
commit;
end;
$$;

create or replace function insertProject(title_ text, id_ integer) returns void as $$
INSERT INTO project (id, title) VALUES (id_,title_);
$$ LANGUAGE SQL;


create or replace function returnProjectTitle(id_ int) returns text
    LANGUAGE plpgsql
as $$
Declare project_title text;
BEGIN
Select title into project_title from project where project.id=id_;
return project_title;
END;
$$;


create or replace function insertAuthor(name_ text) returns integer as $$
INSERT INTO author(name) VALUES (name_) returning id;
$$ LANGUAGE SQL;


create or replace function returnAuthorName(id_ int) returns text
    LANGUAGE plpgsql
as $$
Declare author_name text;
BEGIN
Select name into author_name from author where author.id=id_;
return author_name;
END;
$$;

create or replace procedure insertOrUpdateIssue(issue_id int, project_title text, project_id int, author_name text, assigned_name text,
                                                key_ text, createTime_ timestamp, closedTime_ timestamp, updatedTime_ timestamp, summary_ text, description_ text, type_ text,
                                                priority_ text, status_  text, timeSpent_ bigint)
    LANGUAGE plpgsql
as $$
BEGIN
    if exists (select id from issues where issue_id=issues.id) then
        Call updateIssue(issue_id, assigned_name ,
                         updatedTime_, closedTime_,summary_, description_, type_,
                         priority_, status_, timeSpent_);
else
        call insertIssue(issue_id, project_title, project_id, author_name, assigned_name ,
                         key_, createTime_, closedTime_, updatedTime_, summary_, description_, type_,
                         priority_, status_ , timeSpent_);
end if;
commit;
END;
$$;

create or replace function getLastChangeTime(id_ int) returns timestamp
    LANGUAGE plpgsql
as $$
Declare lastTime timestamp;
BEGIN
    if exists(select * from "statusChange" where "statusChange".issueid=id_) then
select "statusChange".changetime into lastTime from "statusChange" where "statusChange".issueid=id_
order by "statusChange".changetime desc;
else
select to_timestamp(0) into lastTime;
end if;
return lastTime;
END;
$$;


create or replace procedure addOpenTaskTime(id int, createTime timestamp, context json)
    LANGUAGE plpgsql
as $$
BEGIN
    if exists(Select op.data from "openTaskTime" as op where op.projectid = id) then
Delete from "openTaskTime" where projectid = id;
end if;
INSERT INTO "openTaskTime"(projectId, createdTime, data) VALUES (id, createTime, context);
END;
$$;

create or replace procedure addComplexityTaskTime(id int, createTime timestamp, context json)
    LANGUAGE plpgsql
as $$
BEGIN
    if exists(Select op.data from "complexityTaskTime" as op where op.projectid = id) then
Delete from "complexityTaskTime" where projectid = id;
end if;
INSERT INTO "complexityTaskTime"(projectId, createdTime, data) VALUES (id, createTime, context);
END;
$$;

create or replace procedure addTaskStateTime(id int, createTime timestamp, context json, state_ text)
    LANGUAGE plpgsql
as $$
BEGIN
    if exists(Select op.data from "taskStateTime" as op where op.projectid = id and op.state = state_) then
Delete from "taskStateTime" where projectid = id and state = state_;
end if;
INSERT INTO "taskStateTime"(projectId, createdTime, data, state) VALUES (id, createTime, context, state_);
END;
$$;

create or replace procedure addActivityByTask(id int, createTime timestamp, context json, state_ text)
    LANGUAGE plpgsql
as $$
BEGIN
    if exists(Select op.data from "activityByTask" as op where op.projectid = id and op.state = state_) then
Delete from "activityByTask" where projectid = id and state = state_;
end if;
INSERT INTO "activityByTask"(projectId, createdTime, data, state) VALUES (id, createTime, context, state_);
END;
$$;

create or replace procedure addTaskPriorityCount(id int, createTime timestamp, context json, state_ text)
    LANGUAGE plpgsql
as $$
BEGIN
    if exists(Select op.data from "taskPriorityCount" as op where op.projectid = id and op.state = state_) then
Delete from "taskPriorityCount" where projectid = id and state = state_;
end if;
INSERT INTO "taskPriorityCount"(projectId, createdTime, data, state) VALUES (id, createTime, context, state_);
END;
$$;
