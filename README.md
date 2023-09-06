# mserver
mini web server实现
## 基本功能：
1.路由注册
* 静态路由
* 参数路由
* 通配符路由

2.中间件
* 中间件注册
* 可路由中间件(尽可能匹配)，匹配多个，越具体越靠后
  * (/a/b, mws), 
    * /a/b/c 会执行mws
  * (/a/*, mws1), (/a/b/*, mws2)
    * /a/c 执行mws1
    * /a/b/c 执行mws1， mws2
  * (/a/*/c, mws1), (/a/b/c, mws2)
    * /a/d/c 执行mws1
    * /a/b/c 执行mws1，mws2
    * /a/b/d 不执行mws1，mws2!
  * (/a/:id, mws1), (/a/123/c, mws2)
    * /a/123,执行mws1
    * /a/123/c, 执行mws1和mws2
  * 不支持:
    * (/a/*/c, mws1),(/a/b/c, mws2), /a/b/c 执行ms2

3.支持group分组
* 基于group的uri前缀
* 不支持group的中间件(可使用可路由中间件替代)
## 快速使用：
