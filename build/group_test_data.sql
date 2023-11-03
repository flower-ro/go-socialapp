select tg_id,group_id,group_url,is_exit,enable_collect ,collect_status,node_id, account_status, is_banned,collect_fail_reason from tg_group left join
                                                                                                                                   (select tg_id idid,node_id,status account_status, is_banned from account) a on tg_group.tg_id=a.idid ORDER BY enable_collect desc, node_id,collect_status



update tg_group set collect_status='collecting',is_exit=false,enable_collect=true where tg_id='5877603988' and group_id='33';
update tg_group set collect_status='collecting',is_exit=false,enable_collect=true where tg_id='6034068074' and group_id='1';
update tg_group set collect_status='collecting',is_exit=false,enable_collect=true where tg_id='5992485352' and group_id='1';
update tg_group set collect_status='claimed for collect',is_exit=false,enable_collect=true where tg_id='6034068074' and group_id='123';

update tg_group set collect_status='claimed for collect',is_exit=false,enable_collect=true where tg_id='6039551633' and group_id='123';
update tg_group set collect_status='claimed for collect',is_exit=false,enable_collect=true where tg_id='6039551633' and group_id='33';
update tg_group set collect_status='claimed for collect' ,is_exit=false,enable_collect=true where tg_id='5898118872' and group_id='123';

update account set node_id='1001',status='working' where tg_id='5877603988'

update tg_group set enable_collect=FALSE