create table `user_n_role` (
 `created_at` datetime not null default current_timestamp
,`created_by` int not null
,`organization_id` int not null
,`app_user_id` int not null
,`user_role_id` int not null
,primary key(`organization_id`, `app_user_id`, `user_role_id`)
,foreign key(`created_by`) references `app_user`(`id`) on delete cascade
,foreign key(`organization_id`) references `organization`(`id`) on delete cascade
,foreign key(`app_user_id`) references `app_user`(`id`) on delete cascade
,foreign key(`user_role_id`) references `user_role`(`id`) on delete cascade
);
