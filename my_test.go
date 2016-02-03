package trie

import (
    "fmt"
)

func init() {
    var t = NewTrie()
    t.Add("中国人民")
    
    
    t.Add("中国无敌")
    t.Add("中国abc")  
    t.Add("中国人民共和国万岁")
    t.Add("中美友好")
    //var d = t.Dump()
    //fmt.Println(d)
    
    var l = t.PrefixMembersList("中")
    fmt.Println(l)
    
    
    var p = t.PrefixMembersList("中国")
    fmt.Println(p)
    
    //fmt.Println(t)
}
