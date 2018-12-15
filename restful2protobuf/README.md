# Json2Protobuf

``` json
{
    "fields": [
        "name": {
            "order": 1,
            "type": "string",
            "is_required": true,
            "desc": "学生名称"
        },
        "age": {
            "order": 2,
            "type": "int32",
            "is_required": false,
            "desc": "学生年龄"
        },
        "status": {
            "order": 3,
            "type": "int32",
            "default": "0",
            "is_required": true
        },
        "join_time": {
            "order": 4,
            "type": "int64",
            "is_required": false,
            "desc": "加入日期"
        },
        "department": {
            "order": 5,
            "type": "object",
            "is_requierd": true,
            "desc":  "系",
            "fields": [
                "id": {
                    "order": 0,
                    "type": "int64",
                    "is_required": true,
                    "desc": "系id"
                },
                "name": {
                    "order": 1,
                    "type": "string",
                    "is_required": true,
                    "desc": "系名称"
                }
            ]
        },
    ]
}
```

## Json格式

+ order

    order对应protobuf以及thrift的序列号

+ type

    |类型   |描述   |
    |-------|--------|
    |string|字符串|
    |int32|32位整型|
    |int64|64位整型|
    |double|浮点类型|
    |object|对象类型|
    |bool|boolean类型|

+ is_required:

    true为required; false为optional

+ desc:

    字段注释

+ fields

    object的字段

## Protobuf格式

json转化为Protobuf的格式

``` protobuf
message RequestBody
{
    required string name = 1; // 学生名称
    optional int32 age = 2;   // 学生年龄
    required int32 status = 3 [default = 0];
    optional int64 join_date = 4; // 加入日期
    required Department department = 5; // 部门
}

message Department
{
    required int64 id = 1; // 部分id
    required string name = 2; // 系名称
}
```