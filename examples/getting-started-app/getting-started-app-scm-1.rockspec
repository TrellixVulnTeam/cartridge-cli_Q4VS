package = 'getting-started-app'
version = 'scm-1'
source  = {
    url = '/dev/null',
}
-- Put any modules your app depends on here
dependencies = {
    'tarantool',
    'lua >= 5.1',
    'checks == 3.0.1-1',
    'cartridge == 2.7.3-1',
    'ldecnumber == 1.1.3-1',
    'metrics == 0.9.0-1',
}
build = {
    type = 'none';
}
