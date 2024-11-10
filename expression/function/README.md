# function
实现expression扩展函数，目前包含下列函数

// 返回 *

any()


// 对字符串进行crc

adler32(str string) uint32

// 获取网络的内容

// method: GET POST

// url: string

// encoding: form; json GET方法时忽略，POST方法时body的编码方式

// args: 类似struct的赋值模式，例如

{
    a: "年后",
    b: true,
    c: 15
}

fetch(method, url, encoding string, args struct) string
