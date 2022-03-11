CREATE DATABASE `db_rate_limit` DEFAULT CHARACTER SET = `utf8` DEFAULT COLLATE = `utf8_general_ci`;

CREATE TABLE `db_rate_limit`.`t_func_whitelist` (`id` int,`f_func` varchar(255), PRIMARY KEY (id));
INSERT INTO `db_rate_limit`.`t_func_whitelist` (`id`, `f_func`) VALUES ('1', 'set1.pod2.method3');
INSERT INTO `db_rate_limit`.`t_func_whitelist` (`id`, `f_func`) VALUES ('2', 'set2.pod7.method4');

CREATE TABLE `db_rate_limit`.`t_func_blacklist` (`id` int,`f_func` varchar(255), PRIMARY KEY (id));
INSERT INTO `db_rate_limit`.`t_func_blacklist` (`id`, `f_func`) VALUES ('1', 'a.b.c');
INSERT INTO `db_rate_limit`.`t_func_blacklist` (`id`, `f_func`) VALUES ('2', 'x.y.z');

CREATE TABLE `db_rate_limit`.`t_ip_whitelist` (`id` int,`f_ip` varchar(255), PRIMARY KEY (id));
INSERT INTO `db_rate_limit`.`t_ip_whitelist` (`id`, `f_ip`) VALUES ('1', '127.0.0.1');
INSERT INTO `db_rate_limit`.`t_ip_whitelist` (`id`, `f_ip`) VALUES ('2', '127.0.0.2');

CREATE TABLE `db_rate_limit`.`t_ip_blacklist` (`id` int,`f_ip` varchar(255), PRIMARY KEY (id));
INSERT INTO `db_rate_limit`.`t_ip_blacklist` (`id`, `f_ip`) VALUES ('1', '127.0.0.3');
INSERT INTO `db_rate_limit`.`t_ip_blacklist` (`id`, `f_ip`) VALUES ('2', '127.0.0.4');

CREATE TABLE `db_rate_limit`.`t_func_qps` (`id` int,`f_func` varchar(255),`f_qps` int, PRIMARY KEY (id));
INSERT INTO `db_rate_limit`.`t_func_qps` (`id`, `f_func`, `f_qps`) VALUES ('1', 'set1.pod2.method4', '3');
INSERT INTO `db_rate_limit`.`t_func_qps` (`id`, `f_func`, `f_qps`) VALUES ('2', 'z.y.x', '100');
