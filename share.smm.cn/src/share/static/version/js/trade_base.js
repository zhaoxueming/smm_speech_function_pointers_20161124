"use strict";
(function(){
    if(!window.smm){ window.smm = {}; }
    if(!window.smm.trade){ window.smm.trade = {}; }
    var trade = window.smm.trade;
    if(!window.console){window.console={}; }
    if(!window.console.log){window.console.log=function(){}; }

    trade.Cookie = {
        _cookie : {},
        set_cookie : function(key,value,expires,domain,path){
            if(!key){ return; }
            var cookie =  key + '=' + encodeURIComponent(value);
            if(typeof expires === 'number'){
                var date = new Date();
                date.setTime(date.getTime()+ expires * 24 * 60 * 60 * 1000);
                cookie += ';expires=' + date.toUTCString();
            }

            cookie += ';domain=' + (domain ? domain  : trade.Cookie.DOMAIN );
            cookie += ';path='   + (path   ? path    : trade.Cookie.PATH   );

            document.cookie = cookie;
            this._cookie[key]=value;
        },
        get_cookie : function(key){
            return this.__init()._cookie[key] ;
        },
        __init : function(){
            var cookies = $.trim(document.cookie?document.cookie:'');
            if(!cookies){
                return this;
            }
            cookies = cookies.split(';');

            for (var i = 0; i < cookies.length; i++) {
                var ci = cookies[i].split('=');
                this._cookie[$.trim((ci[0]?ci[0]:''))]=decodeURIComponent($.trim((ci[1]?ci[1]:'')));
            }
            return this;
        },
        DOMAIN : '.smm.cn',
        PATH   : '/'
    };

    trade.Format = {
        param : function(param){
            return param ? param : "_";
        },

        _crypto : null,

        base64 : function(obj){
            if(!this._crypto){
                this._crypto = window.CryptoJS;
            }
            return this._crypto.enc.Base64.stringify(this._crypto.enc.Utf8.parse(obj))

                .replace(/\//g, '*');     // replace  '/' to '*' because of the url

        },

        base64_decode : function(obj){
            if(!this._crypto){
                this._crypto = window.CryptoJS;
            }
            return this._crypto.enc.Base64.parse(

                    obj.replace(/\*/g,'/')     //  replace  '/' to '*' because of the url

                ).toString(this._crypto.enc.Utf8);

        },

        replace : function(template, obj){
            for(var p in obj){
                if(!obj.hasOwnProperty(p)){
                    continue;
                }
                template =  template.replace(
                    new window.RegExp( "__" + p.toUpperCase() + "__","mg" ),
                    obj[p]
                );
            }
            return template;
        }
    };
    trade.Api = {
        /*
            {
                path   : '',   must
                method : 'get',
                dataType : 'json',
                data   : {},
                success : defaulton,
                error : error
                quiet : false
            }
        */
        get : function(message){
            if(!message.quiet){
                trade.Pop.cover({
                    msg : "正在与服务器通信...",
                    icon : trade.Pop.loading
                });
            }
            $.ajax({
                url     : message.path + (message.path.indexOf('?') >= 0 ? '&' : '?') +'t=' + new Date().getTime()  ,
                type    : message.method    ? message.method    : 'get',
                dataType: message.dataType  ? message.dataType  : 'json',
                data    : message.data      ? message.data      :  {},
                success : function(data){
                    if(!message.quiet){
                        trade.Pop.rmcover();
                    }
                    message.success ? message.success(data)
                                    : (function(data){
                                       // I just want to print the data
                                    })(data);

                },
                error   : function(err){
                    if(!message.quiet){
                        trade.Pop.rmcover();
                    }
                    message.error   ? message.error(err)
                                    : (function(err){
                                        if(!message.quiet){
                                            trade.Pop.alert({
                                                title : "网络错误！",
                                                msg : "网络错误,与服务器通信失败！",
                                                icon : trade.Pop.error
                                            });
                                        }
                                    })(err);
                }
            });
        },

        _inputfile : null,

        _fileupload_callback : null,

        _fileupload_opt : null,
        /*
            {
                cover : false,
                success : function(){}
                input : input_dom
            }
        */
        fileupload : function(opt){
            this._fileupload_opt = opt|| {};

            this._inputfile = opt.input;
            this._inputfile .attr("name" , "file_name")
                            .fileupload({
                                url : trade.Config.host+'/fileupload/file_name/upload.json',
                                dataType: 'json',
                                done: function (e, data) {
                                    if(trade.Api._fileupload_opt.cover){
                                        trade.Pop.rmcover();
                                        trade.Api._fileupload_opt._cover = false;
                                    }
                                    if(!trade.Api._fileupload_opt.success){
                                        return;
                                    }
                                    trade.Api._fileupload_opt.success(data.result);
                                },
                                progressall: function (e, data) {
                                    if(trade.Api._fileupload_opt.cover && !trade.Api._fileupload_opt._cover ){
                                        trade.Api._fileupload_opt._cover = true;
                                        trade.Pop.cover({
                                            msg : "文件上传中...",
                                            icon : trade.Pop.loading
                                        });
                                    }
                                }
            });
        }
    };

    trade.SlideSearch = {
        /*
            opt = {
                input : $obj_input,
                list : $obj_list,                   //unnecessaries
                minlen : 4 ,                        //unnecessaries
                keyget : function(key){
                    return key;
                }                                   //unnecessaries
                onselect : function(name,id){

                },                                  //unnecessaries
                itemsget : function(key,callback){
                    callback([{
                        name : "XXXX",
                        id   : "id"
                    }])
                }
            }
        */
        bind : function(opt){
            if(!opt.input || !opt.itemsget){
                return;
            }
            if(!opt.minlen){
                opt.minlen = 0;
            }
            if(!opt.keyget){
                opt.keyget = function(key){
                    return key;
                };
            }
            if(!opt.list){
                opt.list = $( window.document.createElement("ul"))
                            .append($( window.document.createElement("li"))
                                    .append($(window.document.createElement("a"))
                                            .attr("href","javascript:void(0)")
                                            .text("未发现匹配的企业！")
                                    )
                            )
                            .addClass("dropdown-menu smm-slide-list")
                            .hide()
                            .insertAfter(opt.input);
            }
            opt.list._open = false;
            opt.list._keys = {};
            opt.list._key = "";
            opt.list._all_items = [];
            opt.list._empty = true;
            opt.list._time = new Date().getTime();

            opt.list._update = function (name){
                var key = opt.keyget(name)
                if(key == opt.list._key ){
                    return;
                }
                var time = new Date().getTime();
                opt.list._get_item(name,function(items){
                    if(time < opt.list._time){
                        return;
                    }
                    opt.list._key = key;
                    opt.list._time = time;
                    opt.list._set_items(items);
                });
            };
            opt.list._add_items = function(arr){
                for (var i = 0; i < arr.length; i++) {
                    opt.list._all_items.push(arr[i]);
                }
            }
            opt.list._get_by_name = function(value){
                if(!value){
                    value = opt.input.val();
                }
                for (var i = 0; i < opt.list._all_items.length; i++) {
                    var co = opt.list._all_items[i]
                    if(co.name == value){
                        return co
                    }
                }
                return null;
            }
            opt.list._get_item = function(name,onitemsget){
                var key = opt.keyget(name)
                var items = opt.list._keys[key];
                if(items){
                    onitemsget(items);
                    return;
                }
                opt.itemsget(name,function(list){
                    if(!list){
                        list = [];
                    }
                    opt.list._keys[key] = list;
                    opt.list._add_items(list);
                    onitemsget(list);
                });
            }

            opt.list._set_items = function(items){
                opt.list.empty();
                if( items.length == 0 ){
                    opt.list._empty  = true;
                    $(document.createElement("li"))
                        .appendTo(opt.list)
                        .append($(document.createElement("a"))
                                    .attr("href","javascript:void(0)")
                                    .text("未发现匹配的选项！")
                        );
                    return;
                }
                opt.list._empty  = false;

                for (var i = 0; i < items.length; i++) {
                    var co = items[i]
                    $(document.createElement("li"))
                        .attr("itemid",co.id)
                        .appendTo(opt.list)
                        .on('click',opt.list._li_click)
                        .append($(document.createElement("a"))
                                    .attr("href","javascript:void(0)")
                                    .text(co.name)
                        );
                }
            }

            opt.list._li_click = function(){
                var self = $(this);
                var name = self.text();
                opt.list._open = false;
                opt.list.hide();
                opt.input.val(name);
                opt.list._update();
                if(opt.onselect){
                    opt.onselect(opt.list._get_by_name(name));
                }

            }

            opt.list._start = function(){
                window.setInterval(function(){
                    var key = opt.input.val() || '';
                    if(key.length < opt.minlen || !opt.list._open){
                        opt.list.hide();
                        return;
                    }
                    opt.list.show();
                    opt.list._update(key);
                },1000);
            }

            opt.input._input = function(){
                opt.list._open = true;
            }

            opt.input._blur = function(){
                opt.list._open = false;
            };


            opt.list._start();
            opt.input.on("input",opt.input._input);
            opt.input.on("focus",opt.input._input);
            opt.input.on("blur" ,opt.input._blur );

            return opt;
        },
    }






})()
