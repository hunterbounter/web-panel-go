CREATE OR REPLACE VIEW VULNERABILTY_REPORT_LIST as
SELECT id as db_id,
       name,
       url as url,
       risk as severity,
       'zap' as agent_type,
       create_date ,
       to_char(t.create_date, 'DD.MM.YYYY HH24:MI:SS') as formatted_create_date
FROM zap_scan_results t


union all

select
    id as db_id,
    t1.nvt_name as name,
    case
        when t1.hostname = '' then t1.ip
        when t1.ip = '' then t1.hostname
        else
            t1.ip || ' - ' || t1.hostname
        end as url,
    t1.severity as risk,
    'openvas' as agent_type,
    create_date ,
    to_char(t1.create_date,
            'DD.MM.YYYY HH24:MI:SS') as formatted_create_date
from
    openvas_scan_result t1

union all

select
    id as db_id,
    t1.name as name,
    t1.ip as url,
    t1.severity as risk,
    'nuclei' as agent_type,
    create_date,
    to_char(t1.create_date, 'DD.MM.YYYY HH24:MI:SS') as formatted_create_date
from
    nuclei_scan_results  t1
order by  create_date DESC;



CREATE OR REPLACE VIEW TARGET_VIEW as
select

    t.value ,
    case
        when t.type = 1 then 'Zap'
        when t.type = 2 then 'OpenVAS'
        when t.type = 3 then 'Nuclei'
        end as type ,
    case
        when t.status = 0 then 'Added'
        when t.status = 1 then 'Waiting'
        when t.status = 2 then 'Sent to Agent'
        when t.status = 3 then 'Finish'
        end as status,
    to_char(t.create_date, 'DD.MM.YYYY HH24:MI:SS') as create_date
from
    targets t
