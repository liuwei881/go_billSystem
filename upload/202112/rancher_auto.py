#coding=utf-8


import cattle
import redis
import dns.tsigkeyring
import dns.update
import dns.query
import json
import salt.client
import nginx


def GetName():
    """从rancher中获取service名称字典"""
    client = cattle.Client(url='http://10.100.17.124:8080/v1',
                           access_key='58108661D7583006AF76',
                           secret_key='SAyABdU9ocbQPncMTpxC4BenBrH4kAqVGQ348bs4')
    name = {}
    p = client.list_service()
    for i in p.items():
        if i[0] == 'data':
            for j in i[1]:
                try:
                    name[j.name] = int(j.data.fields.launchConfig.ports[0].split(':')[0])
                except Exception as e:
                    pass
    return name


def AddDns(server, zone, name, ttl, _type, ip):
    """添加dns A记录"""
    keyring = dns.tsigkeyring.from_text({'other-key': 'WWFjaI4lkvXNkRAIExbFYA=='})
    up = dns.update.Update(zone, keyring=keyring)
    up.add(name, ttl, _type, ip)
    return dns.query.tcp(up, server)


def DelDns(server, zone, name, _type):
    """添加dns A记录"""
    keyring = dns.tsigkeyring.from_text({'other-key': 'WWFjaI4lkvXNkRAIExbFYA=='})
    up = dns.update.Update(zone, keyring=keyring)
    up.delete(name, _type)
    return dns.query.tcp(up, server)


def AddConsul(name, port):
    """创建健康检查的json并发动到rancher host,重启consul"""
    res = {}
    res["id"] = name
    res["name"] = name
    res["tags"] = []
    res["tags"].append(name)
    res["port"] = int(port)
    res["checks"] = []
    res["checks"].append({})
    res["checks"][0]["name"] = "_".join(["nginx", name])
    res["checks"][0]["tcp"] = "localhost:{0}".format(int(port))
    res["checks"][0]["interval"] = "5s"
    res["checks"][0]["timeout"] = "2s"

    with open("/data/salt/tcp_check.json", "r") as f:
        a = json.loads(f.read())
        a = dict(a)
        a["services"].append(res)
    with open("/data/salt/tcp_check.json", "w") as f:
        f.write(json.dumps(a, indent=4))
    print('create tcp_check.json finish')
    local = salt.client.LocalClient()
    local.cmd('p_rancher_group', [
        'cp.get_file',
        'cmd.run',
    ],
              [
                  ['salt://tcp_check.json', '/etc/consul.d/agent/tcp_check.json'],
                  ['/usr/local/bin/consul reload'],
              ], tgt_type='nodegroup')

    print('cp tcp_check.json to consul client and reload consul finish')


def DelConsul(name):
    """删除健康检查的json并发动到rancher host,重启consul"""
    with open("/data/salt/tcp_check.json", "r") as f:
        a = json.loads(f.read())
        a = dict(a)
        [a['services'].remove(i) for i in a['services'] if name == i['name']]

    with open("/data/salt/tcp_check.json", "w") as f:
        f.write(json.dumps(a, indent=4))
    print('del tcp_check.json finish')
    # salt cp file to consul client and reload consul
    local = salt.client.LocalClient()
    local.cmd('p_rancher_group', [
        'cp.get_file',
        'cmd.run',
    ],
              [
                  ['salt://tcp_check.json', '/etc/consul.d/agent/tcp_check.json'],
                  ['/usr/local/bin/consul reload'],
              ], tgt_type='nodegroup')


def AddNginx(name):
    """添加nginx,vhost,upstream及consul-template"""
    with open('/data/salt/pre_nginx_temp/upstream/{0}.conf'.format(name), 'w') as f:
        pass
    ####create server####
    c = nginx.Conf()
    s = nginx.Server()
    s.add(
        nginx.Key('listen', '80'),
        nginx.Key('server_name', '{0}.testyf.ak'.format(name)),
        nginx.Key('access_log', '/var/log/nginx/{0}.testyf.ak.log  logstash'.format(name)),
        nginx.Location('/',
                       nginx.Key('proxy_pass', 'http://{0}'.format(name)),
                       nginx.Key('proxy_set_header', 'Host $host'),
                       nginx.Key('proxy_set_header', 'X-Real-IP $remote_addr'),
                       nginx.Key('proxy_set_header', 'X-Forwarded-For $proxy_add_x_forwarded_for')
                       )
    )
    c.add(s)
    nginx.dumpf(c, '/data/salt/pre_nginx_temp/vhosts/{0}.testyf.ak.conf'.format(name))
    print("create {0}.testyf.cn.conf finish".format(name))
    ####create consul template
    c = nginx.Conf()
    u = nginx.Upstream('{0}'.format(name),
                       nginx.Key('ip_hash', ''),
                       nginx.Key('{{range service "%s"}}' % name, ''),
                       nginx.Key('server', '{{.Address}}:{{.Port}} fail_timeout=0'),
                       nginx.Key('{{else}}', ''),
                       nginx.Key('server', '10.100.20.31:80'),
                       nginx.Key('{{end}}', ''),
                       nginx.Key('keepalive', '64')
                       )
    c.add(u)
    nginx.dumpf(c, '/data/salt/pre_nginx_temp/consul_nginx_temp/{0}.ctmpl'.format(name))
    with open('/data/salt/pre_nginx_temp/consul_nginx_temp/{0}.ctmpl'.format(name), 'r') as f:
        fs = f.readlines()
    fst = []
    for i in fs:
        if i.endswith('} ;\n'):
            i = i.replace('} ;\n', '} \n')
        fst.append(i)
    with open('/data/salt/pre_nginx_temp/consul_nginx_temp/{0}.ctmpl'.format(name), 'w') as f:
        f.write(''.join(fst))
    print('create {0}.ctmpl finish'.format(name))

    # update consul_temp.conf
    with open("/data/salt/pre_nginx_temp/consul_nginx_temp/consul_temp.conf", 'r') as f:
        fs = f.readlines()
    a = 'template {\n', ' source = "/data/nginx_consul_template/{0}.ctmpl"\n'.format(
        name), ' destination = "/usr/local/nginx/upstream/{0}.conf"\n'.format(
        name), ' command = "systemctl reload nginx"\n', '}\n'
    for i in a:
        fs.append(i)
    with open("/data/salt/pre_nginx_temp/consul_nginx_temp/consul_temp.conf", 'w') as f:
        f.write(''.join(fs))
    print('update consul_temp.conf finish')

    # cp upstream vhosts ctmpl consul_temp to nginx server and reload consul-template
    local = salt.client.LocalClient()
    local.cmd('vmlin7542.open.com.cn', [
        'cp.get_file',
        'cp.get_file',
        'cp.get_file',
        'cp.get_file',
        'cmd.run',
    ],
              [
                  ['salt://pre_nginx_temp/upstream/{0}.conf'.format(name),
                   '/usr/local/nginx/upstream/{0}.conf'.format(name)],
                  ['salt://pre_nginx_temp/vhosts/{0}.testyf.ak.conf'.format(name),
                   '/usr/local/nginx/vhosts/{0}.testyf.ak.conf'.format(name)],
                  ['salt://pre_nginx_temp/consul_nginx_temp/{0}.ctmpl'.format(name),
                   '/data/nginx_consul_template/{0}.ctmpl'.format(name)],
                  ['salt://pre_nginx_temp/consul_nginx_temp/consul_temp.conf',
                   '/data/consul_template/consul_temp.conf'],
                  ['kill -HUP `cat /var/run/consul-template.pid`'],
              ])
    print('reload nginx and consul-template finish')


def DelNginx(name):
    """删除nginx,vhost,upstream,consul-template"""
    with open("/data/salt/pre_nginx_temp/consul_nginx_temp/consul_temp.conf", 'r') as fp:
        fp = fp.readlines()
    for i in fp:
        if i.startswith(' source') and i.endswith('{0}.ctmpl"\n'.format(name)):
            fp = fp[0:fp.index(i) - 1] + fp[fp.index(i) + 4:]
    with open("/data/salt/pre_nginx_temp/consul_nginx_temp/consul_temp.conf", 'w') as f:
        f.write(''.join(fp))
    print('update consul_temp.conf finish')
    local = salt.client.LocalClient()
    local.cmd('vmlin7542.open.com.cn', [
        'cmd.run',
        'cmd.run',
        'cmd.run',
        'cp.get_file',
        'cmd.run',
        ],
    [
        ['rm /usr/local/nginx/upstream/{0}.conf -rf'.format(name)],
        ['rm /usr/local/nginx/vhosts/{0}.testyf.ak.conf -rf'.format(name)],
        ['rm /data/nginx_consul_template/{0}.ctmpl -rf'.format(name)],
        ['salt://pre_nginx_temp/consul_nginx_temp/consul_temp.conf',
         '/data/consul_template/consul_temp.conf'],
        ['kill -HUP `cat /var/run/consul-template.pid`'],
    ])
    print('reload nginx and consul-template finish')
    print("del nginx and consul-template finish domainname {0}".format(name))


def DelConsulTemp(name):
    """删除consul模板"""
    with open("/data/salt/pre_nginx_temp/consul_nginx_temp/consul_temp.conf", 'r') as fp:
        fp = fp.readlines()
    for i in fp:
        if i.startswith(' source') and i.endswith('{0}.ctmpl"\n'.format(name)):
            fp = fp[0:fp.index(i) - 1] + fp[fp.index(i) + 4:]
    with open("/data/salt/pre_nginx_temp/consul_nginx_temp/consul_temp.conf", 'w') as f:
        f.write(''.join(fp))
    print('update consul_temp.conf finish')
    local = salt.client.LocalClient()
    local.cmd('vmlin7542.open.com.cn', [
        'cp.get_file',
        'cmd.run',
        ],
    [
        ['salt://pre_nginx_temp/consul_nginx_temp/consul_temp.conf',
         '/data/consul_template/consul_temp.conf'],
        ['kill -HUP `cat /var/run/consul-template.pid`'],
    ])
    print("update consul-template finish domainname {0}".format(name))


if __name__ == '__main__':
    NameCache = redis.StrictRedis(host='10.96.141.112', port=6379, db=2)
    NowName = GetName() # 从rancher接口获取servername
    # NameCache.set('name', list(name_dict.keys()))
    # 以上是在夜里执行
    OriginalName = eval(NameCache.get('newname'))   # 凌晨5点存的servername
    DiffValue = list(set(NowName.iteritems()) - set(OriginalName.iteritems()))  # 现在的与凌晨5点的取差集(或与redis中比差集), 统计出现在有的
    DiffValue_del = list(set(OriginalName.iteritems()) - set(NowName.iteritems()))  # 凌晨5点的与现在的取差集, 统计出原来有的
    try:
        Increment = eval(NameCache.get('newdiffvalue')) # 获取 现在有的servername
    except Exception as e:
        Increment = {}
    DiffValue_in = list(set(Increment.iteritems()) - set(NowName.iteritems()))  # 现在有的与 现在的server取差集
    if DiffValue:
        for i in DiffValue:
            name, port = i
            if name not in list(Increment):
                DelConsul(name)
                DelConsulTemp(name)
            AddConsul(name, port)
            AddNginx(name)
            AddDns('10.100.14.219', 'testyf.ak', name, 60, 'A', '10.100.20.126')
            AddDns('10.100.132.16', 'testyf.ak', name, 60, 'A', '10.100.20.126')

    if DiffValue_del:
        for j in DiffValue_del:
            name, port = j
            if name not in NowName:
                DelConsul(name)
                DelConsulTemp(name)
                DelNginx(name)
                DelDns('10.100.14.219', 'testyf.ak', name, 'A')
                DelDns('10.100.132.16', 'testyf.ak', name, 'A')
        NameCache.set('newname', GetName())
    if DiffValue_in:
        for i in DiffValue_in:
            name, port = i
            if name not in NowName:
                DelConsul(name)
                DelConsulTemp(name)
                DelNginx(name)
                DelDns('10.100.14.219', 'testyf.ak', name, 'A')
                DelDns('10.100.132.16', 'testyf.ak', name, 'A')
    NameCache.set('newdiffvalue', dict(DiffValue))  # 将现在有的入库
    NameCache.set('newname', NowName)