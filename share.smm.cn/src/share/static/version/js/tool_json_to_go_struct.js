(function(){

    var json_to_go_struct = function(json){
        try{
            if(!json||!json.trim()){
                return "";
            }
            var obj = JSON.parse(json);
            return obj == null ? null_to_go_struct : obj._to_go_struct(0,0);
        }catch(e){
            return e;
        }
    }

    Object_to_go_struct = function(spa,lev,outer_max_len){
        var outer ="struct{\r\n"
        var keys = [];
        var hkeys = [];
        var values = [];
        var max_length = 0;
        for(var p in this){
            if(!this.hasOwnProperty(p)){
                continue;
            }
            var hkey = underline_to_hump(p);
            keys.push(p)
            hkeys.push(hkey)
            values.push(this[p]);
            var key_len = hkey.length;
            if(key_len > max_length){
                max_length = key_len;
            }
        }
        max_length++;
        for (var i = 0; i < keys.length; i++) {
            outer += level(lev + 1);
            outer += space(hkeys[i],max_length);
            outer += (values[i] == null ? space(null__to_go_struct,type_len + spa) :
                                         (values[i]._to_go_struct(0,lev + 1,max_length)));
            outer += '`json:"' + keys[i] + '"`';
            outer += "\r\n";
        }
        outer += level(lev) + space("}",type_len + outer_max_len );
        return outer;
    };

    Object.defineProperty(Object.prototype, '_to_go_struct',{
        value : Object_to_go_struct,
        enumerable : false
    });

    Array.prototype._to_go_struct = function(spa,lev,max_len){
        if(!this.length || this[0] == null){
            return space("[]" + null__to_go_struct ,type_len + spa);
        }
        return "[]" + this[0]._to_go_struct(spa - 2,lev,max_len);
    };

    String.prototype._to_go_struct = function(spa){
        return space("string",type_len + spa);
    };

    Boolean.prototype._to_go_struct = function(spa){
        return space("bool",type_len + spa);
    };

    Number.prototype._to_go_struct = function(spa){
        if(window.parseInt(this) == this){
            if(this > 100 || this < -100){
                return space("int64",type_len + spa);
            }
            return space("int",type_len + spa);
        }
        return space("float64",type_len + spa);
    };

    var null_to_go_struct = "interface{}";

    var level=function(lev){
        return space("",lev * tab_len,true);
    }

    var space=function(str,len,rigid){
        var space_char = " ";
        var spa = len - str.length;
        if(spa <= 0){
            if(rigid){
                return str;
            }
            return str + space_char;
        }
        for (var i = 0; i < spa; i++) {
            str += space_char;
        }
        return str;
    }

    var underline_to_hump = function(str){
        var values = str.split("_");
        for (var i = 0; i < values.length; i++) {
            var vi = values[i];
            if(!vi.length){
                continue;
            }
            var first = vi.substring(0,1);
            var remain = vi.substring(1,vi.length);
            values[i] = first.toUpperCase() + remain.toLowerCase();
        }
        return values.join("");
    }

    var tab_len = 4
    var type_len = 15
    var source = window.document.getElementsByClassName('source')[0];
    var target = window.document.getElementsByClassName('target')[0];
    source.oninput = function(){
        target.value = json_to_go_struct(source.value);
        return false;
    };
})();
