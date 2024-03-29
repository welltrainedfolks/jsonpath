# The jsonpath library

jsonpath is used to pull values out of a JSON document without unmarshalling the string into an object.  At the loss of post-parse random access and conversion to primitive types, you gain faster return speeds and lower memory utilization.  If the value you want is located near the start of the json, the evaluator will terminate after reaching and recording its destination.

The evaluator can be initialized with several paths, so you can retrieve multiple sections of the document with just one scan.  Naturally, when all paths have been reached, the evaluator will early terminate.
  
For each value returned by a path, you'll also get the keys & indexes needed to reach that value.  Use the `keys` flag to view this in the CLI.  The Go package will return an `[]interface{}` of length `n` with indexes `0 - (n-2)` being the keys and the value at index `n-1`.

## The history of a library

This fork is owned and currently maintained by [WTFTeam](https://wtfteam.pro). It's based on [that one](https://github.com/JumboInteractiveLimited/jsonpath) from Jumbo Interactive Limited, which sequentally based/copied from NodePrime/jsonpath Github repository, currently unavailable. The MIT license on this code is inherited from Jumbo Interactive Limited fork, so we can support and maintain this library freely. If there is by any chance someone's proprietary code, you can reach us by ``abuse <at> wtfteam <dot> pro`` with details. Unless proved against all code here is licensed under MIT.

### Mirrors and development

The main development goes in [fat0troll's Gitea](https://source.hodakov.me/libraries/jsonpath). There is an official mirror on [GitHub](https://github.com/welltrainedfolks/jsonpath), but we do not guarantee that PRs, opened at GitHub, will be merged.

The canonical import path hosted by [giredore](https://sources.dev.pztrn.name/pztrn/giredore), our own Go packages redirector.

### Go Package

Start with

```shell
go get source.wtfteam.pro/libraries/jsonpath
```

```go
paths, err := jsonpath.ParsePaths(pathStrings ...string) {
```

```go
eval, err := jsonpath.EvalPathsInBytes(json []byte, paths) 
// OR
eval, err := jsonpath.EvalPathsInReader(r io.Reader, paths)
```

then

```go
for {
    if result, ok := eval.Next(); ok {
        fmt.Println(result.Pretty(true)) // true -> show keys in pretty string
    } else {
        break
    }
}
if eval.Error != nil {
    return eval.Error
}
```  

`eval.Next()` will traverse JSON until another value is found.  This has the potential of traversing the entire JSON document in an attempt to find one.  If you prefer to have more control over traversing, use the `eval.Iterate()` method.  It will return after every scanned JSON token and return `([]*Result, bool)`.  This array will usually be empty, but occasionally contain results.

### Path Syntax

All paths start from the root node `$`.  Similar to getting properties in a JavaScript object, a period `.title` or brackets `["title"]` are used.  
  
Syntax|Meaning|Examples
------|-------|-------
`$`|root of doc|  
`.`|property selector |`$.Items`
`["abc"]`|quoted property selector|`$["Items"]`
`*`|wildcard property name|`$.*` 
`[n]`|Nth index of array|`[0]` `[1]`
`[n:m]`|Nth index to m-1 index (same as Go slicing)|`[0:1]` `[2:5]`
`[n:]`|Nth index to end of array|`[1:]` `[2:]`
`[*]`|wildcard index of array|`[*]`
`+`|get value at end of path|`$.title+`
`?(expression)`|where clause (expression can reference current json node with @)|`?(@.title == "ABC")`


Expressions

- paths (that start from current node `@`)
- numbers (integers, floats, scientific notation)
- mathematical operators (+ - / * ^)
- numerical comparisos (< <= > >=)
- logic operators (&& || == !=)
- parentheses `(2 < (3 * 5))`
- static values like (`true`, `false`)
- `@.value > 0.5`

Example:

```javascript
{
    "Items":
        [
            {
                "title": "A Midsummer Night's Dream",
                "tags":[
                    "comedy",
                    "shakespeare",
                    "play"
                ]
            },{
                "title": "A Tale of Two Cities",
                "tags":[
                    "french",
                    "revolution",
                    "london"
                ]
            }
        ]
}
```

Example Paths:

*Paths*  
`$.Items[*].title+`
... "A Midsummer Night's Dream"
... "A Tale of Two Cities"
  
`$.Items[*].tags+`
... ["comedy","shakespeare","play"]  
... ["french","revolution","london"]  
  
`$.Items[*].tags[*]+`  
... "comedy"  
... "shakespeare"  
... "play"  
... "french"  
... "revolution"  
...  "london"  
  
... = keys/indexes of path  
