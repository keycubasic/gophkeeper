SELECT * FROM users;

MERGE INTO creds c
USING (
    SELECT CAST(user_id AS INT) user_id, user_site, user_login, user_password, metadata
      FROM (VALUES (1, 'sdf', '1sdf', '1sf', '1sf')) 
        AS v (user_id, user_site, user_login, user_password, metadata)
    ) t ON t.user_id = c.user_id AND t.user_site = c.user_site
    WHEN MATCHED THEN
        UPDATE SET user_login = t.user_login, user_password = t.user_password, metadata = t.metadata
    WHEN NOT MATCHED THEN
        INSERT (user_id, user_site, user_login, user_password, metadata)
        VALUES (t.user_id, t.user_site, t.user_login, t.user_password, t.metadata);


SELECT * FROM creds;
