/* account sql, admin/admin */
INSERT INTO t_user (name, password, salt, email, verified, created_at, updated_at) VALUES ('admin', '-buSGGz5B6BaqICVeBCXkp_FyEnR3u8wnuwuyaGxncq8A8qOmnw6ez6c_8fX0StnX1QqmcYP5u8kTiswMpKHpA==', 'xhxKQFDaFpLS', '2637309949@qq.com', 1, strftime('%Y-%m-%d %H-%M-%S','now'), strftime('%Y-%m-%d %H-%M-%S','now'))
ON CONFLICT (email) DO
UPDATE SET name='admin';
