USE `settingsdb`;


INSERT INTO `settings`
    (`id`, `name`, `notify`)
VALUES
    (1, 'profit', 1),
    (2, 'max-deals', 1),
    (3, 'max-deal-amount', 0),
    (4, 'extra-status-icon', 0),
    (5, 'extra-access', 0);


INSERT INTO `settings_values`
    (`id`, `setting_id`, `name`, `value`)
VALUES
    ( 1, 1, 'profit-low', '85'),
    ( 2, 1, 'profit-mid', '90'),
    ( 3, 1, 'profit-high', '92'),
    ( 4, 1, 'profit-god', '146'),
    ( 5, 2, 'deals-10', '10'),
    ( 6, 2, 'deals-15' , '15'),
    ( 7, 2, 'deals-20', '20'),
    ( 8, 3, 'amount-100', '100'),
    ( 9, 3, 'amount-500', '500'),
    (10, 3, 'amount-1000', '1000'),
    (11, 4, 'icon-jun', 'http://foo.bar/junior.png'),
    (12, 4, 'icon-mid', 'http://foo.bar/middle.png'),
    (13, 4, 'icon-sen', 'http://foo.bar/senior.png'),
    (14, 5, 'access-jun', ''),
    (15, 5, 'access-mid', 'courses'),
    (16, 5, 'access-sen', 'courses;consultant');


INSERT INTO `bundles`
    (`id`, `tag`, `name`, `parent_id`)
VALUES
    ( 1, 'jun', 'profit-jun',  0),
    ( 2, 'mid', 'profit-mid',  1),
    ( 3, 'sen', 'profit-sen',  2),
    ( 4,  NULL, 'profit-god',  3),
    ( 5, 'jun',  'deals-jun',  0),
    ( 6, 'mid',  'deals-mid',  5),
    ( 7, 'sen',  'deals-sen',  6),
    ( 8, 'jun', 'amount-jun',  0),
    ( 9, 'mid', 'amount-mid',  8),
    (10, 'sen', 'amount-sen',  9),
    (11, 'jun',  'extra-jun',  0),
    (12, 'mid',  'extra-mid', 11),
    (13, 'sen',  'extra-sen', 12);


INSERT INTO `bundles_values`
    (`bundle_id`, `value_id`)
VALUES
    ( 1,  1),
    ( 2,  2),
    ( 3,  3),
    ( 4,  4),
    ( 5,  5),
    ( 6,  6),
    ( 7,  7),
    ( 8,  8),
    ( 9,  9),
    (10, 10),
    (11, 11),
    (11, 14),
    (12, 12),
    (12, 15),
    (13, 13),
    (13, 16);
