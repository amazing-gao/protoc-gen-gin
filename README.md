# Install
```sh
go get github.com/gogo/protobuf/gogoproto && ls $GOPATH/src/github.com/gogo/protobuf

github.com/googleapis/googleapis/google/api

```


# gRPC -> REST

## Google标准

使用这套标准的系统：
  1. https://github.com/googleapis/googleapis
  2. https://cloud.google.com/endpoints/
  3. https://github.com/grpc-ecosystem/grpc-gateway
  4. https://github.com/envoyproxy/envoy

### HttpRule

定义了gRPC/REST映射的标准。

映射协议：
https://github.com/googleapis/googleapis/blob/master/google/api/annotations.proto
https://github.com/googleapis/googleapis/blob/master/google/api/http.proto
https://github.com/googleapis/googleapis/blob/master/google/api/httpbody.proto


### Rules for HTTP mapping
1. Leaf request fields (recursive expansion nested messages in the request
   message) are classified into three categories:
   - Fields referred by the path template. They are passed via the URL path.
   - Fields referred by the [HttpRule.body][google.api.HttpRule.body]. They are passed via the HTTP
     request body.
   - All other fields are passed via the URL query parameters, and the
     parameter name is the field path in the request message. A repeated
     field can be represented as multiple query parameters under the same
     name.
 2. If [HttpRule.body][google.api.HttpRule.body] is "*", there is no URL query parameter, all fields
    are passed via URL path and HTTP request body.
 3. If [HttpRule.body][google.api.HttpRule.body] is omitted, there is no HTTP request body, all
    fields are passed via URL path and URL query parameters.

#### Path template syntax

    Template = "/" Segments [ Verb ] ;
    Segments = Segment { "/" Segment } ;
    Segment  = "*" | "**" | LITERAL | Variable ;
    Variable = "{" FieldPath [ "=" Segments ] "}" ;
    FieldPath = IDENT { "." IDENT } ;
    Verb     = ":" LITERAL ;

The syntax `*` matches a single URL path segment. The syntax `**` matches
zero or more URL path segments, which must be the last part of the URL path
except the `Verb`.

The syntax `Variable` matches part of the URL path as specified by its
template. A variable template must not contain other variables. If a variable
matches a single path segment, its template may be omitted, e.g. `{var}`
is equivalent to `{var=*}`.

The syntax `LITERAL` matches literal text in the URL path. If the `LITERAL`
contains any reserved character, such characters should be percent-encoded
before the matching.

If a variable contains exactly one path segment, such as `"{var}"` or
`"{var=*}"`, when such a variable is expanded into a URL path on the client
side, all characters except `[-_.~0-9a-zA-Z]` are percent-encoded. The
server side does the reverse decoding. Such variables show up in the
[Discovery Document](https://developers.google.com/discovery/v1/reference/apis) as `{var}`.

If a variable contains multiple path segments, such as `"{var=foo/*}"`
or `"{var=**}"`, when such a variable is expanded into a URL path on the
client side, all characters except `[-_.~/0-9a-zA-Z]` are percent-encoded.
The server side does the reverse decoding, except "%2F" and "%2f" are left
unchanged. Such variables show up in the
[Discovery Document](https://developers.google.com/discovery/v1/reference/apis) as `{+var}`.


#### Request

##### Path

1. 使用request中的指定字段名称
`GetMessage(name: "123456")` -> `GET /v1/123456`

```protobuf
service Messaging {
  rpc GetMessage(GetMessageRequest) returns (Message) {
    option (google.api.http) = {
        get: "/v1/{name}"
    };
  }
}
message GetMessageRequest {
  string name = 1; // Mapped to URL path.
}
message Message {
  string text = 1; // The resource content.
}
```

2. 增加前缀
`GetMessage(name: "messages/123456")` -> `GET /v1/messages/123456`

```protobuf
service Messaging {
  rpc GetMessage(GetMessageRequest) returns (Message) {
    option (google.api.http) = {
        get: "/v1/{name=messages/*}"
    };
  }
}
message GetMessageRequest {
  string name = 1; // Mapped to URL path.
}
message Message {
  string text = 1; // The resource content.
}
```

##### Query

`GetMessage(message_id: "123456" revision: 2 sub: SubMessage(subfield: "foo"))` -> `GET /v1/messages/123456?revision=2&sub.subfield=foo`


```protobuf
service Messaging {
  rpc GetMessage(GetMessageRequest) returns (Message) {
    option (google.api.http) = {
        get:"/v1/messages/{message_id}"
    };
  }
}
message GetMessageRequest {
  message SubMessage {
    string subfield = 1;
  }
  string message_id = 1; // Mapped to URL path.
  int64 revision = 2;    // Mapped to URL query parameter `revision`.
  SubMessage sub = 3;    // Mapped to URL query parameter `sub.subfield`.
}
```

##### Body

1. 使用`*`映射所有字段到body
`UpdateMessage(message_id: "123456" message { text: "Hi!" })` -> `PATCH /v1/messages/123456 { "text": "Hi!" }`

```protobuf
service Messaging {
  rpc UpdateMessage(Message) returns (Message) {
    option (google.api.http) = {
      patch: "/v1/messages/{message_id}"
      body: "*"
    };
  }
}
message Message {
  string message_id = 1;
  string text = 2;
}
```

2. 使用字段名映射指定字段到body
`UpdateMessage(message_id: "123456" message { text: "Hi!" })` -> `PATCH /v1/messages/123456 { "text": "Hi!" }`

```protobuf
service Messaging {
  rpc UpdateMessage(UpdateMessageRequest) returns (Message) {
    option (google.api.http) = {
      patch: "/v1/messages/{message_id}"
      body: "message"
    };
  }
}
message Message {
  string text = 1;
}
message UpdateMessageRequest {
  string message_id = 1; // mapped to the URL
  Message message = 2;   // mapped to the body
}
```

#### Response