'use strict';
(function(){

    var trade = window.smm.trade;

    var param = {
        index_hidden        : "smm-index-hidden",
        index_hidden_cookie : "SMM_speech_index_hidden",
        index_fold          : "smm-index-fold",
        auth_token          : "SMM_auth_token",
        dropdown_li         : '.dropdown-menu li',
        notify_no           : '.smm_nav_notify span.badge'
    };

    var dom = {
        config          : $('.smm_config'),
        icon            : $('.smm_nav_icon'),
        exit            : $('.smm_nav_user_exit'),
        body            : $(document.body),
        doc             : $(document),
        venobox         : $(".venobox"),
        print           : $(".smm_print"),
        datepicker      : $('.smm_date_picker'),
        dropdown        : $('.dropdown-menu'),
        alert           : $('.smm_alert'),
        confirm         : $('.smm_confirm'),
        prompt          : $('.smm_prompt'),
        cover           : $('.smm_cover'),
        notify_no       : $('.smm_nav_notify span.badge'),
        notify          : $('.smm_nav_notify'),
        tooltip         : $('[data-toggle="tooltip"]'),
        popover         : $('[data-toggle="popover"]'),
        index_leader    : $('.smm_index_leader')
    };

    dom.icon._click = function(){
        if(!dom.body.hasClass(param.index_hidden)){
            dom.body.addClass(param.index_hidden);
            trade.Cookie.set_cookie(param.index_hidden_cookie,"true" );
        }else{
            dom.body.removeClass(param.index_hidden);
            trade.Cookie.set_cookie(param.index_hidden_cookie,"", -1 );
        }
    };

    dom.index_leader._click = function(){
        var self = $(this);
        var parent = self.parents('ul');
        if(!parent.hasClass(param.index_fold)){
            parent.addClass(param.index_fold);
            trade.Cookie.set_cookie(self.attr("hidden-cookie"),"true" );
        }else{
            parent.removeClass(param.index_fold);
            trade.Cookie.set_cookie(self.attr("hidden-cookie"),"", -1 );
        }
    };

    dom.print._click = function () {
        $($(this).attr('print')).print();
    }


    if(!trade.Config){
        trade.Config = {};
    }

    dom.config.each(function(){
        var self = $(this);
        trade.Config[self.attr("name")] = self.attr("value");
    });

    if(!trade.Pop){
        trade.Pop = {
            success     : "fa fa-3x fa-check-circle-o text-success",
            question    : "fa fa-3x fa-question-circle-o text-warning",
            error       : "fa fa-3x fa-times-circle-o text-danger",
            write       : "fa fa-3x fa-pencil-square-o text-default",
            loading     : "fa fa-3x fa-spinner fa-pulse fa-fw text-success"

        };
    }
    dom.alert.modal({
        show : false,
    }),

    dom.alert.on("hide.bs.modal",function(e){
        if(dom.alert._close){
            dom.alert._close(e);
        }
    });
    /*  {
            title : "title",
            msg  : "message",
            icon : "fa-check",
            onclose : func(){ }
        }*/
    trade.Pop.alert = function(opt){
        dom.alert.find('.smm_alert_title').html(opt.title ? opt.title : "");
        dom.alert.find('.smm_alert_body').html(opt.msg? opt.msg : "");
        dom.alert.find('.smm_alert_icon i').removeClass().addClass(opt.icon ? opt.icon : trade.Pop.success);
        dom.alert._close = opt.onclose ? opt.onclose : function(){};
        dom.alert.modal('show');
    };

    dom.confirm.modal({
        show : false,
    }),

    dom.confirm.on("hide.bs.modal",function(e){
        if(dom.confirm._close){
            dom.confirm._close(e);
        }
    });
    dom.confirm.find('.smm_confirm_ok').on( 'click',function(){
        dom.confirm._close = null;
        if(dom.confirm._ok){
            dom.confirm._ok();
        }
        dom.confirm.modal("hide");
    });

    dom.confirm.find('.smm_alert_checked').on('click',function(){
            dom.confirm.find('.smm_confirm_ok').prop("disabled" , !($(this).prop("checked")));
    })

    /*  {
            title : "title",
            msg  : "message",
            icon : "fa-question",
            check : "checkinfo",            //todo
            checked : false,
            onok : func(){ }
            oncancel : func(){ }
        }*/
    trade.Pop.confirm = function(opt){

        dom.confirm.find('.smm_alert_title').html(opt.title ? opt.title : "");
        dom.confirm.find('.smm_alert_body').html(opt.msg? opt.msg : "");
        if(opt.check){
            dom.confirm.find('.smm_alert_check').html(opt.check? opt.check : "");
            dom.confirm.find('.smm_alert_checked').prop("checked",!!opt.checked);
            dom.confirm.find('.smm_confirm_ok').prop("disabled",!opt.checked);
            dom.confirm.find('.smm_alert_check_line').show();
        }else{
            dom.confirm.find('.smm_alert_check_line').hide();
        }
        dom.confirm.find('.smm_alert_icon i').removeClass().addClass(opt.icon ? opt.icon : trade.Pop.question);
        dom.confirm._close = opt.oncancel ? opt.oncancel : function(){};
        dom.confirm._ok = opt.onok ? opt.onok : function(){};
        dom.confirm.modal('show');
    };

    dom.prompt.modal({
        show : false,
    });

    dom.prompt.on("hide.bs.modal",function(e){
        if(dom.prompt._close){
            dom.prompt._close(e);
        }
    });
    dom.prompt.find('.smm_prompt_ok').on( 'click',function(){
        dom.prompt._close = null;
        if(dom.prompt._ok){
            dom.prompt._ok(dom.prompt.find('.smm_prompt_input').val());
        }
        dom.prompt.modal("hide");
    });

    /*  {
            title : "title",
            msg  : "message",
            msg_foot : "foot msg",
            default_value : "default_value",
            placeholder : "placeholder",
            icon : "fa-question",
            onok : func(info){ }
            oncancel : func(){ }
        }*/
    trade.Pop.prompt = function(opt){

        dom.prompt.find('.smm_alert_title').html(opt.title ? opt.title : "");
        dom.prompt.find('.smm_alert_body').html(opt.msg? opt.msg : "");
        dom.prompt.find('.smm_alert_body_foot').html(opt.msg_foot? opt.msg_foot : "");
        dom.prompt.find('.smm_prompt_input').val(opt.default_value ? opt.default_value : "")
                                            .attr("placeholder",opt.placeholder ? opt.placeholder : "");
        dom.prompt.find('.smm_alert_icon i').removeClass().addClass(opt.icon ? opt.icon : trade.Pop.write);
        dom.prompt._close = opt.oncancel ? opt.oncancel : function(){};
        dom.prompt._ok = opt.onok ? opt.onok : function(){};
        dom.prompt.modal('show');
    };

    dom.cover.modal({
        show : false,
        keyboard : false,
        dropdown : 'static'
    }),
        /*  {
            msg  : "message",
            icon : "fa-question"
        }*/
    trade.Pop.cover = function(opt){

        dom.cover.find('.smm_alert_body').html(opt.msg? opt.msg : "");
        dom.cover.find('.smm_alert_icon i').removeClass().addClass(opt.icon ? opt.icon : trade.Pop.loading);
        dom.cover.modal('show');
    };

    trade.Pop.rmcover = function(){
        dom.cover.modal('hide');
    };

    window.onload = function(){
        dom.tooltip.tooltip({
            // template : '<div class="tooltip smm-tooltip-extract" role="tooltip"><div class="tooltip-arrow"></div><div class="tooltip-inner"></div></div>'
        });
        dom.popover.popover({
            template : '<div class="popover smm-popover" role="tooltip"><div class="arrow"></div><h3 class="popover-title"></h3><div class="popover-content"></div></div>'
        });

        dom.icon.on('click',dom.icon._click);
        dom.index_leader.on('click',dom.index_leader._click);
        if(dom.venobox.size()){
            dom.venobox.venobox();
        }
        if(dom.datepicker.size()){
            dom.datepicker.each(function(){
                var self = $(this);
                self.datetimepicker({
                    format : self.attr("format"),
                    minView : self.attr("minView"),
                    autoclose: 1,
                    language :"zh-CN"
                })
            });
        }
        if(dom.print.size()){
            dom.print.on("click",dom.print._click);
        }
    };
})()
