CREATE OR REPLACE FUNCTION update_openvas_status() RETURNS trigger AS $$
BEGIN
    IF NEW.type = 2 AND NEW.service_status = 'offline' AND EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - NEW.last_seen)) < 30 THEN
        NEW.service_status := 'Updating Database';
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;



CREATE TRIGGER trg_update_openvas_status
    BEFORE INSERT OR UPDATE ON machine_monitor
                         FOR EACH ROW
                         EXECUTE FUNCTION update_openvas_status();