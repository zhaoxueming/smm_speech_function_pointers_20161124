(function(){
    var dom = {
        user : $("#navbar .name").text()
    }
    var fn_get_adapter = function(){
        var old = fn_get_ajaxvalue
        fn_get_ajaxvalue = function(url,jsondata){
            console.log(jsondata)
            console.log(url)
            return old(url,jsondata)
        }
    }

    var get_path = function(){
        var jss = document.scripts;
        var nowpath = ""
        for(var i=jss.length;i>0;i--){
            if(jss[i-1].src.indexOf("ops.js")>-1){
                nowpath=jss[i-1].src;
            }
        }
        return function(path){
            console.log(nowpath)
            console.log(path)
        }
    }();


    fn_get_adapter()

    get_path("/etc")

})();
